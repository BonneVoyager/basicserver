package basicserver

import (
	"errors"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/globalsign/mgo/bson"
	"github.com/kataras/iris/v12"
)

type registerInput struct {
	Email     string    `bson:"email"`
	Password  string    `bson:"password"`
	CreatedAt time.Time `bson:"created_at"`
}

var (
	emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

// ServeRegisterPost serves:
// Method:   POST
// Resource: http://localhost/register
//
// This resource requires `Content-Type` header, e.g.:
//
//    Content-Type: application/json
//
// Sample request to be `POST`ed to the /register resource as `application/json`:
//
//    {
//      "email": "user@example.com",
//      "password": "myPassword"
//    }
//
// If everything goes well, then this will return status code `200` and no response body.
//
// In case of error, this will return status code `400` or `500` and `text/plain` error
// message (e.g. "Incorrect Email") as a response.
//
func (app *BasicApp) ServeRegisterPost() iris.Handler {
	return func(ctx iris.Context) {
		var input registerInput
		err := ctx.ReadJSON(&input)
		if err != nil {
			app.HandleError(err, ctx, iris.StatusBadRequest)
			return
		}

		inputEmail := input.Email
		if !emailRegexp.MatchString(inputEmail) {
			err := errors.New("Incorrect " + inputEmail + " Email")
			app.HandleError(err, ctx, iris.StatusBadRequest)
			ctx.WriteString("Incorrect Email")
			return
		}

		var user User
		err = app.Coll.Users.Find(bson.M{"email": inputEmail}).One(&user)
		if err == nil && user.Email != "" {
			err := errors.New("Email " + inputEmail + " Taken")
			app.HandleError(err, ctx, iris.StatusBadRequest)
			ctx.WriteString("Email Taken")
			return
		}

		passByte := []byte(input.Password)
		passEnc, err := bcrypt.GenerateFromPassword(passByte, bcrypt.DefaultCost)
		if err != nil {
			app.HandleError(err, ctx, iris.StatusInternalServerError)
			return
		}

		input.Password = string(passEnc)
		input.CreatedAt = time.Now()
		err = app.Coll.Users.Insert(input)
		if err != nil {
			app.HandleError(err, ctx, iris.StatusInternalServerError)
			return
		}

		app.LogMessage("User " + inputEmail + " registered.")
	}
}
