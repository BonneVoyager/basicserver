package basicserver

import "github.com/globalsign/mgo/bson"

// State is an user's state data entity:
//
//    `ID` user uid
//    `Data` user data
//
type State struct {
	ID   bson.ObjectId `bson:"_id" json:"id"`
	Data bson.M        `bson:"data"`
}
