package basicserver

import (
	"github.com/globalsign/mgo/bson"
	"github.com/kataras/iris"
)

// ServeDataGet serves
// Method:   GET
// Resource: http://localhost/api/data
//
// This resource requires `Authorization` header, e.g.:
//
//		Authorization: Bearer {token}
//
// If everything goes well, then this will return status code `200` and `application/json`
// response with the stored data:
//
//    {
//      "foo": "bar",
//      "bar": "foo"
//    }
//
// In case of error, this will return status code `400` or `500` and `text/plain` error
// message as a response.
//
// In case of invalid/expired token, this will return status code `401` and `text/plain`
// error message as a response.
//
func (app *BasicApp) ServeDataGet() iris.Handler {
	return func(ctx iris.Context) {
		uid := ctx.Values().Get("uid").(string)
		objectUID := bson.ObjectIdHex(uid)

		var state State
		err := app.Coll.States.FindId(objectUID).One(&state)
		if err != nil {
			if err.Error() == "not found" {
				ctx.JSON(State{})
			} else {
				app.HandleError(err, ctx, iris.StatusInternalServerError)
			}
			return
		}

		ctx.JSON(state.Data)
	}
}
