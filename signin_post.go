package basicserver

import (
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/globalsign/mgo/bson"
	"github.com/kataras/iris"
	"golang.org/x/crypto/bcrypt"
)

type signinInput struct {
	Email    string `bson:"email"`
	Password string `bson:"password"`
}

// ServeSigninPost serves
// Method:   POST
// Resource: http://localhost/signin
//
// This resource requires `Content-Type` header, e.g.:
//
//    Content-Type: application/json
//
// Sample request to be `POST`ed to the /signin resource as `application/json`:
//
//    {
//      "email": "user@example.com",
//      "password": "myPassword"
//    }
//
// If everything goes well, then this will return status code `200` and `application/json`
// response with JWT token and it’s expiration milliseconds elapsed since UNIX epoch:
//
//    {
//      "expires": 1543567182285,
//      "token": "..."
//    }
//
// In case of error, this will return status code `400` or `500` and `text/plain` error
// message (e.g. "Incorrect Credentials") as a response.
//
func (app *BasicApp) ServeSigninPost() iris.Handler {
	return func(ctx iris.Context) {
		var input signinInput
		err := ctx.ReadJSON(&input)
		if err != nil {
			app.HandleError(err, ctx, iris.StatusBadRequest)
			return
		}

		inputEmail := input.Email
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

		userPassByte := []byte(user.Password)
		inputPassByte := []byte(input.Password)
		err = bcrypt.CompareHashAndPassword(userPassByte, inputPassByte)
		if err != nil {
			app.HandleError(err, ctx, iris.StatusUnauthorized)
			ctx.WriteString("Incorrect Credentials")
			return
		}

		expiresAt := time.Now().Add(time.Hour * time.Duration(72)).Unix() // 72 hours
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"uid": user.ID.Hex(),
			"exp": expiresAt,
		})
		tokenString, err := token.SignedString(app.Settings.Secret)
		if err != nil {
			log.Print(err)
			app.HandleError(err, ctx, iris.StatusInternalServerError)
			return
		}

		ctx.JSON(iris.Map{
			"expires": expiresAt,
			"token":   tokenString,
		})
	}
}
