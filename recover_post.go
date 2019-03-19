package basicserver

import (
	"log"

	"github.com/globalsign/mgo/bson"
	"github.com/kataras/iris"
	gomail "gopkg.in/gomail.v2"
)

type recoverInput struct {
	Email string `bson:"email" form:"email"`
}

// ServeRecoverPasswordPost serves
// Method:   POST
// Resource: http://localhost/recover
//
// This resource accepts `application/json` and `application/x-www-form-urlencoded`
// `Content-Type` headers.
//
// Sample request to be `POST`ed to the /recover resource as `application/json`:
//
//    {
//      "email": "user@example.com"
//    }
//
// If everything goes well, then this will return status code `200` and no response body.
//
// In case of error, this will return status code `400` or `500` and `text/plain` error
// message (e.g. "No Such User") as a response.
//
// In case of invalid/expired token, this will return status code `401` and `text/plain`
// error message as a response.
//
func (app *BasicApp) ServeRecoverPasswordPost() iris.Handler {
	return func(ctx iris.Context) {
		var formInput, jsonInput recoverInput
		err := ctx.ReadForm(&formInput)
		err = ctx.ReadJSON(&jsonInput)

		var inputEmail string
		if formInput.Email != "" {
			inputEmail = formInput.Email
		} else if jsonInput.Email != "" {
			inputEmail = jsonInput.Email
		}
		var user User
		err = app.Coll.Users.Find(bson.M{"email": inputEmail}).One(&user)
		if err != nil {
			if err.Error() == "not found" {
				app.HandleError(err, ctx, iris.StatusUnauthorized)
				ctx.WriteString("No Such User")
			} else {
				app.HandleError(err, ctx, iris.StatusInternalServerError)
			}
			return
		}

		var recCode string
		var foundNewCode bool
		updateUserCode := func() {
			recCode = user.GenerateRecoveryCode()
			err = app.Coll.Users.Find(bson.M{"recovery_code": recCode}).One(&user)
			if err != nil {
				if err.Error() == "not found" {
					app.Coll.Users.UpdateId(user.ID, bson.M{
						"$set": bson.M{"recovery_code": recCode},
					})
					foundNewCode = true
				}
			}
		}
		for !foundNewCode {
			updateUserCode()
		}

		if app.Settings.SMTP.URL == "" || app.Settings.SMTP.Port == 0 {
			ctx.WriteString("SMTP account not configured.")
			return
		}

		m := gomail.NewMessage()
		m.SetHeader("From", app.Settings.SMTP.User)
		m.SetHeader("To", inputEmail)
		m.SetHeader("Subject", "Password Recovery Link")
		msg := ""
		if app.Settings.RecoverTemplate != "" {
			msg += app.Settings.RecoverTemplate
		} else {
			msg += "Use link below to reset your password:<br />http://localhost/"
		}
		msg += "recover/" + recCode + "</body></html>"
		m.SetBody("text/html", msg)
		d := gomail.NewPlainDialer(
			app.Settings.SMTP.URL,
			app.Settings.SMTP.Port,
			app.Settings.SMTP.User,
			app.Settings.SMTP.Pass,
		)
		if err := d.DialAndSend(m); err != nil {
			log.Println(err)
			app.HandleError(err, ctx, iris.StatusInternalServerError)
			return
		}

		ctx.WriteString("Recovery Email sent. Check your inbox.")
	}
}
