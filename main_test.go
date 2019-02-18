package basicserver

import (
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/globalsign/mgo/bson"
	"github.com/kataras/iris"
	"github.com/kataras/iris/httptest"
	"golang.org/x/crypto/bcrypt"
)

var testUID = bson.NewObjectId()

const testEmail = "testing@example.com"
const testPassword = "mySecret123Password"
const testSecret = "testSecret"

func getSettings() *Settings {
	secret := testSecret
	mongoString := "mongodb://127.0.0.1:27017/test"
	serverPort := "8080"
	logLevel := ""

	return &Settings{
		Secret:      []byte(secret),
		MongoString: mongoString,
		ServerPort:  serverPort,
		LogLevel:    logLevel,
	}
}

var settings = getSettings()
var app = CreateApp(settings)

func createTestUser() {
	passByte := []byte(testPassword)
	encryptedPassword, _ := bcrypt.GenerateFromPassword(passByte, bcrypt.DefaultCost)
	app.Coll.Users.Insert(bson.M{
		"_id":      testUID,
		"email":    testEmail,
		"password": encryptedPassword,
	})
}

func removeTestUser() {
	app.Coll.Users.Remove(bson.M{"email": testEmail})
}

func removeTestState() {
	app.Coll.States.Remove(bson.M{"_id": testUID})
}

func createTestToken() string {
	expiresAt := time.Now().Add(time.Minute * time.Duration(1)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": testUID.Hex(),
		"exp": expiresAt,
	})
	tokenString, _ := token.SignedString([]byte(testSecret))
	return tokenString
}

func TestRegisterPost(t *testing.T) {
	app.Init()

	e := httptest.New(t, app.Iris)

	removeTestUser()

	// wrong email
	e.POST("/register").
		WithJSON(bson.M{
			"email":    "wrongEmail.com",
			"password": testPassword,
		}).
		Expect().Status(httptest.StatusBadRequest).
		Body().Equal("Incorrect Email")

	// correct registration
	e.POST("/register").
		WithJSON(bson.M{
			"email":    testEmail,
			"password": testPassword,
		}).
		Expect().Status(iris.StatusOK).
		NoContent()

	// taken email
	e.POST("/register").
		WithJSON(bson.M{
			"email":    testEmail,
			"password": testPassword,
		}).
		Expect().Status(iris.StatusBadRequest).
		Body().Equal("Email Taken")

	removeTestUser()
}

func TestSigninPost(t *testing.T) {
	e := httptest.New(t, app.Iris)

	createTestUser()

	// correct credentials
	e.POST("/signin").
		WithJSON(bson.M{
			"email":    testEmail,
			"password": testPassword,
		}).
		Expect().Status(httptest.StatusOK)

	// non existing
	e.POST("/signin").
		WithJSON(bson.M{
			"email":    testEmail + "-NOPE",
			"password": testPassword,
		}).
		Expect().Status(httptest.StatusUnauthorized).
		Body().Equal("No Such User")

	// incorrect credentials
	e.POST("/signin").
		WithJSON(bson.M{
			"email":    testEmail,
			"password": testPassword + "-nope",
		}).
		Expect().Status(httptest.StatusUnauthorized).
		Body().Equal("Incorrect Credentials")

	// keepalive GET request with incorrect token
	e.GET("/keepalive").
		WithHeader("Authorization", "Bearer falseToken").
		Expect().Status(httptest.StatusUnauthorized)

	// keepalive GET request with correct token
	token := createTestToken()
	e.GET("/keepalive").
		WithHeader("Authorization", "Bearer "+token).
		Expect().Status(httptest.StatusOK)

	// account DELETE requests
	e.DELETE("/account").
		WithHeader("Authorization", "Bearer "+token).
		Expect().Status(httptest.StatusOK)

	// account should not exist anymore
	e.GET("/keepalive").
		WithHeader("Authorization", "Bearer "+token).
		Expect().Status(httptest.StatusUnauthorized).
		Body().Equal("No Such User")
}

func TestPasswordRecovery(t *testing.T) {
	e := httptest.New(t, app.Iris)

	createTestUser()

	// POST request to generate recovery code
	e.POST("/recover").
		WithJSON(bson.M{"email": testEmail}).
		Expect().Status(httptest.StatusOK).
		Body().Equal("SMTP account not configured.")

	var user User
	_ = app.Coll.Users.Find(bson.M{"email": testEmail}).One(&user)

	// PUT request without providing password field
	e.POST("/change").
		WithJSON(bson.M{
			"code": user.RecoveryCode,
		}).
		Expect().Status(httptest.StatusBadRequest)

	// PUT request with password field
	e.POST("/change").
		WithJSON(bson.M{
			"code":     user.RecoveryCode,
			"password": testPassword + "-new",
		}).
		Expect().Status(httptest.StatusOK)

	// correct credentials
	e.POST("/signin").
		WithJSON(bson.M{
			"email":    testEmail,
			"password": testPassword + "-new",
		}).
		Expect().Status(httptest.StatusOK)

	removeTestUser()
}

func TestApiData(t *testing.T) {
	e := httptest.New(t, app.Iris)

	createTestUser()
	token := createTestToken()

	// POST request with incorrect token
	e.POST("/api/data").
		WithHeader("Authorization", "Bearer falseToken").
		WithJSON(bson.M{"foo": "bar"}).
		Expect().Status(httptest.StatusUnauthorized)

	// correct POST request
	e.POST("/api/data").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(bson.M{
			"foo": "bar",
			"bar": "foo",
		}).
		Expect().Status(httptest.StatusOK)

	// GET request with incorrect token
	e.GET("/api/data").
		WithHeader("Authorization", "Bearer falseToken").
		Expect().Status(httptest.StatusUnauthorized)

	// correct GET request
	e.GET("/api/data").
		WithHeader("Authorization", "Bearer "+token).
		Expect().Status(httptest.StatusOK).
		JSON().Equal(bson.M{
		"foo": "bar",
		"bar": "foo",
	})

	// test DELETE requests
	e.DELETE("/api/data").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON([1]string{"bar"}).
		Expect().Status(httptest.StatusOK)

	e.GET("/api/data").
		WithHeader("Authorization", "Bearer "+token).
		Expect().Status(httptest.StatusOK).
		JSON().Equal(bson.M{
		"foo": "bar",
	})

	e.DELETE("/api/data").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(true).
		Expect().Status(httptest.StatusOK)

	e.GET("/api/data").
		WithHeader("Authorization", "Bearer "+token).
		Expect().Status(httptest.StatusOK).
		JSON().Equal(nil)

	removeTestUser()
	removeTestState()
}

func TestApiFile(t *testing.T) {
	e := httptest.New(t, app.Iris)

	createTestUser()
	token := createTestToken()

	// POST request with incorrect token
	e.POST("/api/file").
		WithHeader("Authorization", "Bearer falseToken").
		WithMultipart().WithFile("file", "golang.jpg").
		Expect().Status(httptest.StatusUnauthorized)

	// correct POST request
	e.POST("/api/file").
		WithHeader("Authorization", "Bearer "+token).
		WithMultipart().WithFile("file", "golang.jpg").
		Expect().Status(httptest.StatusOK)

	// incorrect GET request
	e.GET("/api/file/golang.nope").
		WithHeader("Authorization", "Bearer "+token).
		Expect().Status(httptest.StatusUnauthorized).
		Body().Equal("No Such File")

	// GET request with incorrect token
	e.GET("/api/file/golang.jpg").
		WithHeader("Authorization", "Bearer falseToken").
		WithMultipart().WithFile("file", "golang.jpg").
		Expect().Status(httptest.StatusUnauthorized)

	// correct GET request
	e.GET("/api/file/golang.jpg").
		WithHeader("Authorization", "Bearer "+token).
		Expect().Status(httptest.StatusOK).
		ContentType("image/jpeg")

	// test DELETE requests
	e.DELETE("/api/file").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(bson.M{
			"name": "golang.jpg",
		}).
		Expect().Status(httptest.StatusOK)

	e.GET("/api/file/golang.jpg").
		WithHeader("Authorization", "Bearer "+token).
		Expect().Status(httptest.StatusUnauthorized)

	removeTestUser()
	removeTestState()

	app.Coll.Files.Remove(testUID.Hex() + ":golang.jpg")
}
