package model

import (
	"github.com/globalsign/mgo"
)

// MongoSession ...
var MongoSession *mgo.Session

// Dial to Mongodb, write to mongoSession
func Dial(url string) error {
	s, err := mgo.Dial(url)
	if err != nil {
		return err
	}
	MongoSession = s
	return nil
}
