package basicserver

import "github.com/globalsign/mgo/bson"

// User is an user data entity:
//
//    `ID` is user uid
//    `Email` user email
//    `Password` encrypted password
//
type User struct {
	ID       bson.ObjectId `bson:"_id" json:"id"`
	Email    string        `bson:"email" json:"email"`
	Password string        `bson:"password" json:"password"`
}
