package basicserver

import (
	"time"

	"github.com/kataras/iris"
)

// ServeFileGet serves
// Method:   GET
// Resource: http://localhost/api/file
//
// This resource requires `Authorization` header, e.g.:
//
//		Authorization: Bearer {token}
//
// If everything goes well then we will receive status code 200 and response with the file
// associated with provided id in query parameter.
//
// In case of error, this will return status code `400` or `500` and `text/plain` error
// message as a response.
//
// In case of invalid/expired token, this will return status code `401` and `text/plain`
// error message as a response.
//
func (app *BasicApp) ServeFileGet() iris.Handler {
	return func(ctx iris.Context) {
		fileID := ctx.Params().Get("id")
		uid := ctx.Values().Get("uid").(string)
		fileName := uid + ":" + fileID

		file, err := app.Coll.Files.Open(fileName)
		if err != nil {
			if err.Error() == "not found" {
				app.HandleError(err, ctx, iris.StatusUnauthorized)
				ctx.WriteString("No Such File")
			} else {
				app.HandleError(err, ctx, iris.StatusInternalServerError)
			}
			return
		}
		defer file.Close()

		ctx.ServeContent(file, fileName, time.Now(), true)
	}
}
