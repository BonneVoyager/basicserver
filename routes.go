package basicserver

import (
	"github.com/kataras/iris"
)

// Init configures default server routes:
//
//    `POST /register` serves for user registration
//    `POST /signin` serves for user login
//    `POST /api/data` serves to update user state
//    `POST /api/file` serves to upload user file
//    `GET /api/data` serves to get user data
//    `GET /api/file/{id:string}` serves to get user file
//    `DELETE /api/data` serves to delete user data
//    `DELETE /api/file` serves to delete user file
//
// Check BasicApp.Serve* functions for more details about specific handlers.
//
func (app *BasicApp) Init() {
	app.Iris.Post("/register", app.ServeRegisterPost())
	app.Iris.Post("/signin", app.ServeSigninPost())

	api := app.Iris.Party("/api")
	api.Use(app.RequireAuth())
	{
		api.Post("/data", app.ServeDataPost())
		api.Post("/file", app.ServeFilePost())
		api.Get("/data", app.ServeDataGet())
		api.Get("/file/{id:string}", app.ServeFileGet())
		api.Delete("/data", app.ServeDataDelete())
		api.Delete("/file", app.ServeFileDelete())
	}
}

// Start starts listening on given port.
func (app *BasicApp) Start(port string) {
	app.Iris.Run(iris.Addr(":" + port))
}
