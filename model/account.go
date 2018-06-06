package model

import (
	"context"
	"fmt"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/nanozuki/uexky/uuid64"
)

// CtxTokenKey is contxt key for token
type CtxTokenKey struct{}

// 24 charactors Base64 token
var tokenGenerator = uuid64.Generator{Sections: []uuid64.Section{
	&uuid64.RandomSection{Length: 10},
	&uuid64.CounterSection{Length: 2, Unit: time.Millisecond},
	&uuid64.TimestampSection{Length: 7, Unit: time.Millisecond},
	&uuid64.RandomSection{Length: 5},
}}

// aid generator for post, thread and anonymous id.
var aidGenerator = uuid64.Generator{Sections: []uuid64.Section{
	&uuid64.TimestampSection{Length: 6, Unit: time.Second, NoPadding: true},
	&uuid64.CounterSection{Length: 2, Unit: time.Second},
	&uuid64.RandomSection{Length: 1},
}}

const (
	nameLimit = 5
)

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

// AddName ...
func (a *Account) AddName(ctx context.Context, name string) error {
	if len(a.Names) >= nameLimit {
		return fmt.Errorf("You already have %v names, cannot add more", len(a.Names))
	}
	a.Names = append(a.Names, name)
	c, cs := Colle("accounts")
	defer cs()
	if err := c.Update(bson.M{"token": a.Token}, bson.M{
		"$set": bson.M{"names": a.Names},
	}); err != nil {
		return err
	}
	return nil
}

// HaveName ...
func (a *Account) HaveName(name string) bool {
	for _, n := range a.Names {
		if name == n {
			return true
		}
	}
	return false
}

type accountAID struct {
	ObjectID    bson.ObjectId `bson:"_id"`
	Token       string        `bson:"token"`
	ThreadID    string        `bson:"thread_id"`
	AnonymousID string        `bson:"anonymous_id"`
}

// AnonymousID ...
func (a *Account) AnonymousID(threadID string, new bool) (string, error) {
	c, cs := Colle("accounts_aid")
	c.EnsureIndexKey("thread_id", "token")
	defer cs()

	newAID := func() (string, error) {
		aid, err := aidGenerator.New()
		if err != nil {
			return "", err
		}
		aaid := accountAID{
			ObjectID:    bson.NewObjectId(),
			Token:       a.Token,
			ThreadID:    threadID,
			AnonymousID: aid,
		}
		if err := c.Insert(&aaid); err != nil {
			return "", err
		}
		return aaid.AnonymousID, nil
	}

	if new {
		return newAID()
	}
	query := c.Find(bson.M{"thread_id": threadID, "token": a.Token})
	if count, err := query.Count(); err != nil {
		return "", err
	} else if count == 0 {
		return newAID()
	}
	var aaid accountAID
	if err := query.One(&aaid); err != nil {
		return "", err
	}
	return aaid.AnonymousID, nil
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
