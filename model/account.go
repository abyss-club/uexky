package model

import (
	"net/http"

	"github.com/globalsign/mgo/bson"

	"github.com/CrowsT/uexky/api"
	"github.com/CrowsT/uexky/uuid"
)

var tokenGenerator = uuid.Generator{Sections: []uuid.Section{
	&uuid.RandomSection{Length: 10},
	&uuid.CounterSection{Length: 2},
	&uuid.TimestampSection{Length: 7},
	&uuid.RandomSection{Length: 5},
}}

// Account for uexky
type Account struct {
	ID    bson.ObjectId `json:"-" bson:"_id"`
	Token string        `json:"token" bson:"token"`
}

// NewAccount make a new account
func NewAccount() (*Account, error) {
	token, err := tokenGenerator.New()
	if err != nil {
		return nil, err
	}

	account := &Account{bson.NewObjectId(), token}
	session := mongoSession.Copy()
	accountColle := session.DB("test").C("accounts")
	if err := accountColle.Insert(account); err != nil {
		return nil, err
	}
	return account, nil
}

// GetAccount by token
func GetAccount(token string) (*Account, error) {
	accountColle := mongoSession.Copy().DB("test").C("accounts")
	query := accountColle.Find(bson.M{"token": token})
	var account Account
	if count, err := query.Count(); err != nil {
		return nil, err
	} else if count == 0 {
		return nil, &api.HTTPError{http.StatusNotFound, "Can't find User"}
	}
	if err := query.One(&account); err != nil {
		return nil, err
	}
	return &account, nil
}
