# basicserver

## What is this

basicserver is a preconfigured [Iris based](https://iris-go.com/) web server with user authentication (using [JWT](https://jwt.io/) token) and key/values+files storage (based on [MongoDB](https://www.mongodb.com/)).

## Examples usage

.env

```ini
MONGO=mongodb://127.0.0.1:27017/project
PORT=8081
SECRET=mySecretKey
LOG_LEVEL=info
SMTP_URL=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=...@gmail.com
SMTP_PASS=...
```

main.go

```go
package main

import (
	"os"

	"github.com/bonnevoyager/basicserver"
	"github.com/joho/godotenv"
)

type MyApp struct{}

func getSettings() *basicserver.Settings {
	godotenv.Load()

	secret := os.Getenv("SECRET")
	mongoString := os.Getenv("MONGO")
	serverPort := os.Getenv("PORT")
	logLevel := os.Getenv("LOG_LEVEL")

	smtpURL := os.Getenv("SMTP_URL")
	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")

	SMTPSettings := &basicserver.SMTPSettings{
		URL:  smtpURL,
		Port: smtpPort,
		User: smtpUser,
		Pass: smtpPass,
	}

	return &basicserver.Settings{
		Secret:      []byte(secret),
		MongoString: mongoString,
		ServerPort:  serverPort,
		LogLevel:    logLevel,
		SMTP:        *SMTPSettings,
	}
}

func main() {
	settings := getSettings()

	myApp := &MyApp{}

	app := basicserver.CreateApp(settings)

	api := app.Iris.Party("/api")
	api.Use(app.RequireAuth())
	{
		api.Get("/profile", myApp.ServeProfileGet(app))
		api.Post("/profile", myApp.ServeProfilePost(app))
	}

	app.Init()
	app.Start(settings.ServerPort)
}

```

In order to use basicserver, you need to at least provide `MongoString` and `ServerPort` [configuration options](https://github.com/bonnevoyager/basicserver/blob/master/main.go#L21-L49) to `basicserver.CreateApp(settings)`.

Preconfigured [routes](https://github.com/bonnevoyager/basicserver/blob/master/routes.go#L7-L20) are:

- [POST /register](https://github.com/bonnevoyager/basicserver/blob/master/register_post.go)
- [POST /signin](https://github.com/bonnevoyager/basicserver/blob/master/signin_post.go)
- [GET /recover](https://github.com/bonnevoyager/basicserver/blob/master/recover_get.go)
- [POST /recover](https://github.com/bonnevoyager/basicserver/blob/master/recover_post.go)
- [POST /change](https://github.com/bonnevoyager/basicserver/blob/master/change_post.go)
- [DELETE /account](https://github.com/bonnevoyager/basicserver/blob/master/account_delete.go)
- [GET /keepalive](https://github.com/bonnevoyager/basicserver/blob/master/keepalive_get.go)
- [POST /api/data](https://github.com/bonnevoyager/basicserver/blob/master/data_post.go)
- [POST /api/file](https://github.com/bonnevoyager/basicserver/blob/master/file_post.go)
- [GET /api/data](https://github.com/bonnevoyager/basicserver/blob/master/data_get.go)
- [GET /api/file/{id:string}](https://github.com/bonnevoyager/basicserver/blob/master/file_get.go)
- [DELETE /api/data](https://github.com/bonnevoyager/basicserver/blob/master/data_delete.go)
- [DELETE /api/file](https://github.com/bonnevoyager/basicserver/blob/master/file_delete.go)

You can add additional routes as in the example above, by adding more handlers.

In case you need user authorization, you can use [app.RequireAuth()](https://github.com/bonnevoyager/basicserver/blob/master/require_auth.go).

[Godoc link](https://godoc.org/github.com/BonneVoyager/basicserver).

## Testing

Since basicserver needs MongoDB connection, a running instance of MongoDB Server should be running.

Test configuration tries to connect to default `mongodb://127.0.0.1:27017/test`.

## License

MIT
