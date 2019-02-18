package basicserver

import (
	"math/rand"

	"github.com/globalsign/mgo/bson"
)

// User is an user data entity:
//
//    `ID` is user uid
//    `Email` user email
//    `Password` encrypted password
//
type User struct {
	ID           bson.ObjectId `bson:"_id" json:"id"`
	Email        string        `bson:"email"`
	Password     string        `bson:"password"`
	RecoveryCode string        `bson:"recovery_code"`
}

var codeLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var codeLettersLength = len(codeLetters)

// GenerateRecoveryCode creates a recovery code
func (user *User) GenerateRecoveryCode() string {
	b := make([]rune, 12)
	for i := range b {
		b[i] = codeLetters[rand.Intn(codeLettersLength)]
	}
	return string(b)
}
