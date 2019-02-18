package basicserver

import (
	"github.com/kataras/iris"
)

// ServeRecoverPasswordGet serves
// Method:   GET
// Resource: http://localhost/recover/{code:string}
//
// This should return status code `200` and response body with html form.
//
func (app *BasicApp) ServeRecoverPasswordGet(v string) iris.Handler {
	return func(ctx iris.Context) {
		switch v {
		case "":
			ctx.HTML(`
			<form method="POST">
				<input name="email" type="text" placeholder="Email" />
				<button type="submit">Send</button>
			</form>`)
		case "code":
			recCode := ctx.Params().Get("code")
			ctx.HTML(`
			<form action="/change" method="POST">
				<input name="password" type="text" placeholder="Password" />
				<input name="code" type="hidden" value="` + recCode + `" />
				<button type="submit">Send</button>
			</form>`)
		case "done":
			ctx.WriteString(`Password Changed!`)
		}
	}
}
