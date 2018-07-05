package model

import (
	"os"

	"github.com/globalsign/mgo"
	"gitlab.com/abyss.club/uexky/mgmt"
)

// MongoSession ...
var pkg struct {
	mongoSession *mgo.Session
	database     string
	mainTags     []string
}

// Init to Mongodb, write to mongoSession
func Init() error {
	dbUri, found := os.LookupEnv("MONGO_URI")
	if !found {
		dbUri = mgmt.Config.Mongo.URI
	}

	s, err := mgo.Dial(dbUri)
	if err != nil {
		return err
	}
	pkg.mongoSession = s
	pkg.database = mgmt.Config.Mongo.DB
	pkg.mainTags = mgmt.Config.MainTags
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

// MainTags ...
func MainTags() []string {
	return pkg.mainTags
}
