package model

import (
	"github.com/globalsign/mgo"
)

// MongoSession ...
var pkg struct {
	mongoSession *mgo.Session
	database     string
	mainTags     []string
}

// Init to Mongodb, write to mongoSession
func Init(url, db string, mainTags []string) error {
	s, err := mgo.Dial(url)
	if err != nil {
		return err
	}
	pkg.mongoSession = s
	pkg.database = db
	pkg.mainTags = mainTags
	return nil
}

// Colle return collection by specified name
func Colle(collection string) (*mgo.Collection, func()) {
	session := pkg.mongoSession.Copy()
	colle := session.DB(pkg.database).C(collection)
	close := func() {
		session.Close()
	}
	return colle, close
}
