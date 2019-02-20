package basicserver

import (
	"io"

	"github.com/kataras/iris"
)

// ServeFilePost serves
// Method:   POST
// Resource: http://localhost/api/file
//
// This resource requires `Authorization` header, e.g.:
//
//    Content-Type: multipart/form-data
//    Authorization: Bearer {token}
//
// In order to store a file, a POST request to /api/file resource need to be send as
// `multipart/form-data`. Only one file is accepted per request. Field name of uploaded
// file should be "file" and it’s filename will be it’s id. Resubmitted files
// (files with the same filenames) are overwritten.
//
// If everything goes well, then this will return status code `200` and no response body.
//
// In case of error, this will return status code `400` or `500` and `text/plain` error
// message as response.
//
// In case of invalid/expired token, this will return status code `401` and `text/plain`
// error message as a response.
//
func (app *BasicApp) ServeFilePost() iris.Handler {
	return func(ctx iris.Context) {
		file, info, err := ctx.FormFile("file")
		if err != nil {
			app.HandleError(err, ctx, iris.StatusInternalServerError)
			return
		}
		defer file.Close()

		uid := ctx.Values().Get("uid").(string)
		fileName := uid + ":" + info.Filename
		_ = app.Coll.Files.Remove(fileName)
		newFile, err := app.Coll.Files.Create(fileName)
		if err != nil {
			app.HandleError(err, ctx, iris.StatusInternalServerError)
			return
		}
		defer newFile.Close()

		fileReader := io.Reader(file)
		_, err = io.Copy(newFile, fileReader)
		if err != nil {
			app.HandleError(err, ctx, iris.StatusInternalServerError)
			return
		}
	}
}
