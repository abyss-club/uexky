package model

import (
	"github.com/globalsign/mgo"
)

// MongoSession ...
var mongoSession *mgo.Session
var database string

// Dial to Mongodb, write to mongoSession
func Dial(url, db string) error {
	s, err := mgo.Dial(url)
	if err != nil {
		return err
	}
	mongoSession = s
	database = db
	return nil
}

// Colle return collection by specified name
func Colle(collection string) (*mgo.Collection, func()) {
	session := mongoSession.Copy()
	colle := session.DB(database).C(collection)
	close := func() {
		session.Close()
	}
	return colle, close
}
