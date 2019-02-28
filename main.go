package basicserver

import (
	"log"

	"github.com/kataras/iris"

	mgo "github.com/globalsign/mgo"
)

const usersCollection = "users"
const statesCollection = "states"
const filesCollection = "files"

type collections struct {
	Users  *mgo.Collection
	States *mgo.Collection
	Files  *mgo.GridFS
}

// SMTPSettings values are used by BasicApp to send emails.
type SMTPSettings struct {
	URL  string
	Port int
	User string
	Pass string
}

// Settings values are used by BasicApp. At least `MongoString` and `ServerPort` are required.
//
//  Following values are possible:
//
//   `LogLevel` - available values are: "disable", "fatal", "error", "warn", "info", "debug"
//   `MongoString` - URI format described at http://docs.mongodb.org/manual/reference/connection-string/
//   `Secret` - secret value used by JWT parser
//   `SingleLogin` - allows to access restricted resources only with fresh token received from signin
//   `ServerPort` - port on which the server should listen to
//   `RecoverTemplate` - html content to be sent along with password recovery email
//   `SMTP` - SMTP configuration to send emails
//
type Settings struct {
	LogLevel        string
	MongoString     string
	Secret          []byte
	SingleLogin     bool
	ServerPort      string
	RecoverTemplate string
	SMTP            SMTPSettings
}

// BasicApp contains following fields:
//
//   `Coll.Users` - MongoDB "users" collection
//   `Coll.State` - MongoDB "states" collection
//   `Coll.File` - MongoDB "files" collection
//   `Db` - MongoDB named database
//   `Iris` - iris.Default() instance
//   `Settings` - Settings passed as an argument
//
type BasicApp struct {
	Coll     *collections
	Db       *mgo.Database
	Iris     *iris.Application
	Settings *Settings
}

// CreateApp returns BasicApp.
//
// `settings` argument should contain at least `MongoString` and `ServerPort` fields.
//
// BasicApp contains following fields:
//
//   `Coll.Users` - MongoDB "users" collection
//   `Coll.State` - MongoDB "states" collection
//   `Coll.File` - MongoDB "files" collection
//   `Db` - MongoDB named database
//   `Iris` - iris.Default() instance
//   `Settings` - Settings passed as an argument
//
func CreateApp(settings *Settings) *BasicApp {
	if settings.MongoString == "" {
		log.Fatal("MongoString cannot be empty!")
	}
	if settings.ServerPort == "" {
		log.Fatal("ServerPort cannot be empty!")
	}

	session, err := mgo.Dial(settings.MongoString)
	if err != nil {
		log.Fatal(err)
	}
	db := session.DB("")

	usersC := db.C(usersCollection)
	statesC := db.C(statesCollection)
	filesC := db.GridFS(filesCollection)

	usersC.EnsureIndex(mgo.Index{
		Key:        []string{"recovery_code"},
		Background: true,
	})

	filesC.Files.EnsureIndex(mgo.Index{
		Key:        []string{"filename"},
		Unique:     true,
		DropDups:   true,
		Background: true,
	})

	app := &BasicApp{
		Coll: &collections{
			Users:  usersC,
			States: statesC,
			Files:  filesC,
		},
		Db:       db,
		Iris:     iris.Default(),
		Settings: settings,
	}

	app.Iris.Logger().SetLevel(settings.LogLevel)

	return app
}
