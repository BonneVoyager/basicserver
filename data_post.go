package basicserver

import (
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/kataras/iris"
)

// ServeDataPost serves
// Method:   POST
// Resource: http://localhost/api/data
//
// This resource requires `Authorization` header, e.g.:
//
//		Content-Type: application/json
//		Authorization: Bearer {token}
//
// Sample request to be `POST`ed to the /api/data resource as `application/json`:
//
//    {
//      "foo": "bar",
//      "bar": "foo"
//    }
//
// If everything goes well, then this will return status code `200` and no response body.
//
// In case of error, this will return status code `400` or `500` and `text/plain` error
// message (e.g. "Incorrect Credentials") as a response.
//
// In case of invalid/expired token, this will return status code `401` and `text/plain`
// error message as a response.
//
func (app *BasicApp) ServeDataPost() iris.Handler {
	return func(ctx iris.Context) {
		var input bson.M
		err := ctx.ReadJSON(&input)
		if err != nil {
			app.HandleError(err, ctx, iris.StatusBadRequest)
			return
		}

		uid := ctx.Values().Get("uid").(string)
		objectUID := bson.ObjectIdHex(uid)

		parsedInput := make(bson.M)
		for key, value := range input {
			parsedInput["data."+key] = value
		}
		parsedInput["updated_at"] = time.Now()

		_, err = app.Coll.States.UpsertId(objectUID, bson.M{"$set": parsedInput})
		if err != nil {
			app.HandleError(err, ctx, iris.StatusInternalServerError)
			return
		}
	}
}
