package model

import (
	"context"
	"fmt"

	"github.com/globalsign/mgo/bson"
	"github.com/nanozuki/uexky/uuid"
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
	Names []string      `json:"names" bson:"names"`
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

	account := &Account{bson.NewObjectId(), token, []string{}}
	c, cs := Colle("accounts")
	defer cs()
	if err := c.Insert(account); err != nil {
		return nil, err
	}
	return account, nil
}

// GetAccount by token
func GetAccount(ctx context.Context) (*Account, error) {
	account, err := requireSignIn(ctx)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func requireSignIn(ctx context.Context) (*Account, error) {
	token, ok := ctx.Value(CtxTokenKey{}).(string)
	if !ok || token == "" {
		return nil, fmt.Errorf("Forbidden, no access token")
	}

	c, cs := Colle("accounts")
	defer cs()

	query := c.Find(bson.M{"token": token})
	var account Account
	if count, err := query.Count(); err != nil {
		return nil, err
	} else if count == 0 {
		return nil, fmt.Errorf("Invalid token")
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
