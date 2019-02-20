package basicserver

import (
	"errors"

	"github.com/globalsign/mgo/bson"
	"github.com/kataras/iris"
	"golang.org/x/crypto/bcrypt"
)

type changeInput struct {
	Password string `bson:"password" form:"password"`
	Code     string `bson:"code" form:"code"`
}

// ServeChangePasswordPut serves
// Method:   PUT
// Resource: http://localhost/change
//
// This resource accepts `application/json` and `application/x-www-form-urlencoded`
// `Content-Type` headers.
//
// Sample request to be sent as `PUT` to the /change resource as `application/json`:
//
//    {
//      "password": "newSecretPassword",
//      "code": "abcde123"
//    }
//
// If everything goes well, then this will return status code `200` and no response body.
//
// In case of error, this will return status code `400` or `500` and `text/plain` error
// message as a response.
//
// In case of invalid/expired token, this will return status code `401` and `text/plain`
// error message as a response.
//
func (app *BasicApp) ServeChangePasswordPut() iris.Handler {
	return func(ctx iris.Context) {
		var formInput, jsonInput changeInput
		err := ctx.ReadForm(&formInput)
		err = ctx.ReadJSON(&jsonInput)

		var inputPassword string
		if formInput.Password != "" {
			inputPassword = formInput.Password
		} else if jsonInput.Password != "" {
			inputPassword = jsonInput.Password
		}
		var inputCode string
		if formInput.Code != "" {
			inputCode = formInput.Code
		} else if jsonInput.Code != "" {
			inputCode = jsonInput.Code
		}
		if inputPassword == "" {
			err = errors.New("Password cannot be Empty")
			app.HandleError(err, ctx, iris.StatusBadRequest)
			return
		}
		var user User
		err = app.Coll.Users.Find(bson.M{"recovery_code": inputCode}).One(&user)
		if err != nil {
			if err.Error() == "not found" {
				app.HandleError(err, ctx, iris.StatusUnauthorized)
				ctx.WriteString("No Such User")
			} else {
				app.HandleError(err, ctx, iris.StatusInternalServerError)
			}
			return
		}

		passByte := []byte(inputPassword)
		passEnc, err := bcrypt.GenerateFromPassword(passByte, bcrypt.DefaultCost)
		if err != nil {
			app.HandleError(err, ctx, iris.StatusInternalServerError)
			return
		}
		app.Coll.Users.UpdateId(user.ID, bson.M{
			"$set":   bson.M{"password": string(passEnc)},
			"$unset": bson.M{"recovery_code": ""},
		})

		ctx.Redirect("/recover/done")
	}
}
