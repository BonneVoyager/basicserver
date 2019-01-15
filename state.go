package basicserver

import "github.com/globalsign/mgo/bson"

// State is an user's state data entity:
//
//    `ID` is user uid
//    `Data` is user data
//
type State struct {
	ID   bson.ObjectId `bson:"_id" json:"id"`
	Data bson.M        `bson:"data" json:"data"`
}
