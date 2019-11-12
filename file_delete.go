package basicserver

import (
	"errors"

	"github.com/globalsign/mgo/bson"
	"github.com/kataras/iris/v12"
)

// ServeFileDelete serves
// Method:   DELETE
// Resource: http://localhost/api/file
//
// This resource requires `Authorization` header, e.g.:
//
//    Content-Type: application/json
//    Authorization: Bearer {token}
//
// In order to delete a file, a DELETE request to /api/file resource need to be send.
//
//    {
//      "name": "uploaded_image.jpg"
//    }
//
// If everything goes well, then this will return status code `200` and no response body.
//
// In case of error, this will return status code `400` or `500` and `text/plain` error
// message as response.
//
// In case of invalid/expired token, this will return status code `401` and `text/plain`
// error message as a response.
//
func (app *BasicApp) ServeFileDelete() iris.Handler {
	return func(ctx iris.Context) {
		var input bson.M
		err := ctx.ReadJSON(&input)
		if err != nil {
			app.HandleError(err, ctx, iris.StatusBadRequest)
			return
		}

		uid := ctx.Values().Get("uid").(string)

		filename := input["name"].(string)
		if filename == "" {
			err = errors.New("Name Field Not Provided")
			app.HandleError(err, ctx, iris.StatusBadRequest)
			return
		}

		err = app.Coll.Files.Remove(uid + ":" + filename)
		if err != nil {
			app.HandleError(err, ctx, iris.StatusBadRequest)
			return
		}
	}
}
