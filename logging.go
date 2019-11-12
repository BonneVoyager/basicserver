package basicserver

import (
	"github.com/kataras/iris/v12"
)

// HandleError logs the error and sets status code response header.
func (app *BasicApp) HandleError(err error, ctx iris.Context, status int) {
	app.Iris.Logger().Error(err)
	ctx.StatusCode(status)
}

// LogMessage logs string message on info level
func (app *BasicApp) LogMessage(message string) {
	app.Iris.Logger().Infof(message)
}
