package basicserver

import (
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
)

// ServeKeepAliveGet serves
// Method:   GET
// Resource: http://localhost/keepalive
//
// This resource requires `Authorization` header, e.g.:
//
//    Authorization: Bearer {token}
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
// message (e.g. "Incorrect Authorization Header") as a response.
//
func (app *BasicApp) ServeKeepAliveGet() iris.Handler {
	return func(ctx iris.Context) {
		uid := ctx.Values().Get("uid").(string)

		expiresAt := time.Now().Add(time.Hour * time.Duration(72)).Unix() // 72 hours
		claimsMap := jwt.MapClaims{
			"uid": uid,
			"exp": expiresAt,
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsMap)
		if app.Settings.SingleLogin { // substain single login value
			claimsMap["sl"] = ctx.Values().Get("sl").(string)
		}
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
