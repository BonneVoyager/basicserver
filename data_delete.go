package basicserver

import (
	"errors"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/kataras/iris/v12"
)

// ServeDataDelete serves
// Method:   DELETE
// Resource: http://localhost/api/data
//
// This resource requires `Authorization` header, e.g.:
//
//		Content-Type: application/json
//		Authorization: Bearer {token}
//
// Sample requests needs to be send as DELETE to the /api/data resource as application/json:
//
// This will remove all they keys contained in the array:
//
//    [ "-12198394893", "23749713845" ]
//
// Bellow will removed all the user data from user storage
//
//    true
//
// If everything goes well, then this will return status code `200` and no response body.
//
// In case of error, this will return status code `400` or `500` and `text/plain` error
// message (e.g. "Unsupported Input") as a response.
//
// In case of invalid/expired token, this will return status code `401` and `text/plain`
// error message as a response.
//
func (app *BasicApp) ServeDataDelete() iris.Handler {
	return func(ctx iris.Context) {
		var input interface{}
		err := ctx.ReadJSON(&input)
		if err != nil {
			app.HandleError(err, ctx, iris.StatusBadRequest)
			return
		}

		uid := ctx.Values().Get("uid").(string)
		objectUID := bson.ObjectIdHex(uid)

		unsetInput := make(bson.M)
		switch i := input.(type) {
		case []interface{}:
			for _, v := range i {
				unsetInput["data."+v.(string)] = ""
			}
		case bool:
			unsetInput["data"] = ""
		default:
			err = errors.New("Unsupported Input")
			app.HandleError(err, ctx, iris.StatusBadRequest)
			return
		}
		updateInput := bson.M{
			"$set":   bson.M{"updated_at": time.Now()},
			"$unset": unsetInput,
		}

		_, err = app.Coll.States.UpsertId(objectUID, updateInput)
		if err != nil {
			app.HandleError(err, ctx, iris.StatusInternalServerError)
			return
		}
	}
}
