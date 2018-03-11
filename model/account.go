package model

import (
	"context"
	"fmt"

	"github.com/CrowsT/uexky/uuid"
	"github.com/globalsign/mgo/bson"
)

// CtxTokenKey is contxt key for token
type CtxTokenKey struct{}

// 24 charactors Base64 token
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
func NewAccount(ctx context.Context) (*Account, error) {
	if err := requireNotSignIn(ctx); err != nil {
		return nil, err
	}

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
func GetAccount(ctx context.Context, token string) (*Account, error) {
	account, err := requireSignIn(ctx)
	if err != nil {
		return nil, err
	}

	if account.Token != token {
		return nil, fmt.Errorf("Forbidden")
	}
	return account, nil
}

func requireSignIn(ctx context.Context) (*Account, error) {
	token, ok := ctx.Value(CtxTokenKey{}).(string)
	if !ok || token == "" {
		return nil, fmt.Errorf("Forbidden, no access token")
	}

	session := mongoSession.Copy()
	defer session.Close()

	query := session.DB("test").C("accounts").Find(bson.M{"token": token})
	var account Account
	if count, err := query.Count(); err != nil {
		return nil, err
	} else if count == 0 {
		return nil, fmt.Errorf("Can't find User '%s'", token)
	}
	if err := query.One(&account); err != nil {
		return nil, err
	}
	return &account, nil
}

func requireNotSignIn(ctx context.Context) error {
	token, ok := ctx.Value(CtxTokenKey{}).(string)
	if ok && token != "" {
		return fmt.Errorf("You have already signed in")
	}
	return nil
}
