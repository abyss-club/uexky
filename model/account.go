package model

import (
	"context"
	"fmt"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky/uuid64"
)

// ContextKey ...
type ContextKey string

// ContextKeyToken ...
const ContextKeyToken = ContextKey("token")

// ContextLoggedInAccount ...
const ContextLoggedInAccount = ContextKey("loggedIn")

// aid generator for post, thread and anonymous id.
var aidGenerator = uuid64.Generator{Sections: []uuid64.Section{
	&uuid64.TimestampSection{Length: 6, Unit: time.Second, NoPadding: true},
	&uuid64.CounterSection{Length: 2, Unit: time.Second},
	&uuid64.RandomSection{Length: 1},
}}

const (
	nameLimit = 5
	tagLimit  = 15
)

// Account for uexky
type Account struct {
	ID    bson.ObjectId `json:"id" bson:"_id"`
	Email string        `json:"email" bson:"email"`
	Names []string      `json:"names" bson:"names"`
	Tags  []string      `json:"tags" bson:"tags"`
}

// GetAccount by id (in context)
func GetAccount(ctx context.Context) (*Account, error) {
	account, err := requireSignIn(ctx)
	if err != nil {
		return nil, err
	}
	return account, nil
}

// GetAccountByEmail ...
func GetAccountByEmail(ctx context.Context, email string) (*Account, error) {
	c, cs := Colle("accounts")
	defer cs()
	c.EnsureIndexKey("email")

	query := c.Find(bson.M{"email": email})
	count, err := query.Count()
	if err != nil {
		return nil, err
	}
	if count != 0 {
		var account *Account
		if err := query.One(&account); err != nil {
			return nil, err
		}
		return account, nil
	}

	// New Account
	account := &Account{
		ID:    bson.NewObjectId(),
		Email: email,
	}
	if _, err := c.Upsert(bson.M{"email": email}, bson.M{"$set": account}); err != nil {
		return nil, err
	}
	return account, nil
}

func isNameUesd(name string) (bool, error) {
	c, cs := Colle("accounts")
	defer cs()
	count, err := c.Find(bson.M{"names": name}).Count()
	return count != 0, err
}

// AddName ...
func (a *Account) AddName(ctx context.Context, name string) error {
	if len(a.Names) >= nameLimit {
		return fmt.Errorf("You already have %v names, cannot add more", len(a.Names))
	}
	if used, err := isNameUesd(name); err != nil {
		return errors.Wrapf(err, "Check name '%s'", name)
	} else if used {
		return fmt.Errorf("This name is already in uesd")
	}

	names := append(a.Names, name)
	c, cs := Colle("accounts")
	defer cs()
	if err := c.Update(bson.M{"_id": a.ID}, bson.M{
		"$set": bson.M{"names": names},
	}); err != nil {
		return err
	}
	a.Names = names
	return nil
}

// HaveName ...
func (a *Account) HaveName(name string) bool {
	for _, n := range a.Names {
		if n == name {
			return true
		}
	}
	return false
}

// SyncTags ...
func (a *Account) SyncTags(ctx context.Context, tags []string) error {
	tagSet := map[string]struct{}{}
	tagList := []string{}
	for _, tag := range tags {
		tagSet[tag] = struct{}{}
	}
	for tag := range tagSet {
		tagList = append(tagList, tag)
	}
	if len(tagList) > tagLimit {
		tagList = tagList[:tagLimit]
	}

	c, cs := Colle("accounts")
	defer cs()
	if err := c.Update(bson.M{"_id": a.ID}, bson.M{
		"$set": bson.M{"tags": tagList},
	}); err != nil {
		return err
	}
	a.Tags = tagList
	return nil
}

type accountAID struct {
	ObjectID    bson.ObjectId `bson:"_id"`
	AccountID   bson.ObjectId `bson:"account_id"`
	Token       string        `bson:"token"`
	ThreadID    string        `bson:"thread_id"`
	AnonymousID string        `bson:"anonymous_id"`
}

// AnonymousID ...
func (a *Account) AnonymousID(threadID string, new bool) (string, error) {
	c, cs := Colle("accounts_aid")
	c.EnsureIndexKey("thread_id", "account_id")
	defer cs()

	newAID := func() (string, error) {
		aid, err := aidGenerator.New()
		if err != nil {
			return "", err
		}
		aaid := accountAID{
			ObjectID:    bson.NewObjectId(),
			AccountID:   a.ID,
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
	query := c.Find(bson.M{"thread_id": threadID, "account_id": a.ID})
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
	accountID, ok := ctx.Value(ContextLoggedInAccount).(bson.ObjectId)
	if !ok {
		return nil, fmt.Errorf("Forbidden, no access token")
	}
	c, cs := Colle("accounts")
	defer cs()

	var account *Account
	if err := c.FindId(accountID).One(&account); err != nil {
		return nil, errors.Wrap(err, "Find account")
	}
	return account, nil
}
