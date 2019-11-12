package basicserver

import (
	"errors"

	"github.com/globalsign/mgo/bson"
	"github.com/kataras/iris/v12"
)

type fileItem struct {
	ID       bson.ObjectId `bson:"_id"`
	Filename string        `bson:"filename"`
}

// ServeRemoveAccountDelete serves
// Method:   DELETE
// Resource: http://localhost/account
//
// This resource requires `Authorization` header, e.g.:
//
//		Authorization: Bearer {token}
//
// If everything goes well, then this will return status code `200` and no response body.
//
// In case of error, this will return status code `400` or `500` and `text/plain` error
// message as a response.
//
// In case of invalid/expired token, this will return status code `401` and `text/plain`
// error message as a response.
//
func (app *BasicApp) ServeRemoveAccountDelete() iris.Handler {
	return func(ctx iris.Context) {
		uid := ctx.Values().Get("uid").(string)
		objectUID := bson.ObjectIdHex(uid)

		// remove all user files
		filenamePrefix := "^" + uid + ":"
		var item fileItem
		iter := app.Coll.Files.Find(bson.M{"filename": bson.RegEx{Pattern: filenamePrefix}}).Iter()
		for iter.Next(&item) {
			err := app.Coll.Files.RemoveId(item.ID)
			if err != nil {
				app.HandleError(err, ctx, iris.StatusInternalServerError)
				return
			}
		}

		// then remove user state
		err := app.Coll.States.RemoveId(objectUID)
		if err != nil {
			if err.Error() != "not found" { // state might not be existing yet
				app.HandleError(err, ctx, iris.StatusInternalServerError)
			}
		}

		// to finally remove the user
		err = app.Coll.Users.RemoveId(objectUID)
		if err != nil {
			if err.Error() == "not found" {
				err := errors.New("Account Not Found")
				app.HandleError(err, ctx, iris.StatusUnauthorized)
				return
			}
			app.HandleError(err, ctx, iris.StatusInternalServerError)
			return
		}
	}
}
