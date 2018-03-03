package model

import (
	"github.com/globalsign/mgo"
)

var mongoSession *mgo.Session

// Dial to Mongodb, write to mongoSession
func Dial(url string) error {
	s, err := mgo.Dial(url)
	if err != nil {
		return err
	}
	mongoSession = s
	return nil
}
