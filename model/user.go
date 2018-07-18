package model

import (
	"context"
	"fmt"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky/mw"
	"gitlab.com/abyss.club/uexky/uuid64"
)

// aid generator for post, thread and anonymous id.
var aidGenerator = uuid64.Generator{Sections: []uuid64.Section{
	&uuid64.TimestampSection{Length: 6, Unit: time.Second, NoPadding: true},
	&uuid64.CounterSection{Length: 2, Unit: time.Second},
	&uuid64.RandomSection{Length: 1},
}}

const (
	nameLimit = 5
	tagLimit  = 15
	// ContextKeyUser for logged in user
	ContextKeyUser = mw.ContextKey("user")
)

// User for uexky
type User struct {
	ID    bson.ObjectId `json:"id" bson:"_id"`
	Email string        `json:"email" bson:"email"`
	Name  string        `json:"names" bson:"name"`
	Tags  []string      `json:"tags" bson:"tags"`
}

// GetUser by id (in context)
func GetUser(ctx context.Context) (*User, error) {
	user, err := requireSignIn(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByEmail ...
func GetUserByEmail(ctx context.Context, email string) (*User, error) {
	c := mw.GetMongo(ctx).C(colleUser)
	c.EnsureIndexKey("email")

	query := c.Find(bson.M{"email": email})
	count, err := query.Count()
	if err != nil {
		return nil, err
	}
	if count != 0 {
		var user *User
		if err := query.One(&user); err != nil {
			return nil, err
		}
		return user, nil
	}

	// New User
	user := &User{
		ID:    bson.NewObjectId(),
		Email: email,
	}
	if _, err := c.Upsert(bson.M{"email": email}, bson.M{"$set": user}); err != nil {
		return nil, err
	}
	return user, nil
}

func isNameUesd(ctx context.Context, name string) (bool, error) {
	c := mw.GetMongo(ctx).C(colleUser)
	c.EnsureIndexKey("name")

	count, err := c.Find(bson.M{"name": name}).Count()
	return count != 0, err
}

// SetName ...
func (a *User) SetName(ctx context.Context, name string) error {
	if a.Name != "" {
		return fmt.Errorf("You already have name '%v'", a.Name)
	}
	if used, err := isNameUesd(ctx, name); err != nil {
		return errors.Wrapf(err, "Check name '%s'", name)
	} else if used {
		return fmt.Errorf("This name is already in uesd")
	}

	c := mw.GetMongo(ctx).C(colleUser)
	if err := c.Update(bson.M{"_id": a.ID}, bson.M{
		"$set": bson.M{"name": name},
	}); err != nil {
		return err
	}
	a.Name = name
	return nil
}

// SyncTags ...
func (a *User) SyncTags(ctx context.Context, tags []string) error {
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

	c := mw.GetMongo(ctx).C(colleUser)
	if err := c.Update(bson.M{"_id": a.ID}, bson.M{
		"$set": bson.M{"tags": tagList},
	}); err != nil {
		return err
	}
	a.Tags = tagList
	return nil
}

type userAID struct {
	ObjectID    bson.ObjectId `bson:"_id"`
	UserID      bson.ObjectId `bson:"user_id"`
	ThreadID    string        `bson:"thread_id"`
	AnonymousID string        `bson:"anonymous_id"`
}

// AnonymousID ...
func (a *User) AnonymousID(ctx context.Context, threadID string, new bool) (string, error) {
	c := mw.GetMongo(ctx).C(colleUser)
	c.EnsureIndexKey("thread_id", "user_id")

	newAID := func() (string, error) {
		aid, err := aidGenerator.New()
		if err != nil {
			return "", err
		}
		uaid := userAID{
			ObjectID:    bson.NewObjectId(),
			UserID:      a.ID,
			ThreadID:    threadID,
			AnonymousID: aid,
		}
		if err := c.Insert(&uaid); err != nil {
			return "", err
		}
		return uaid.AnonymousID, nil
	}

	if new {
		return newAID()
	}
	query := c.Find(bson.M{"thread_id": threadID, "user_id": a.ID})
	if count, err := query.Count(); err != nil {
		return "", err
	} else if count == 0 {
		return newAID()
	}
	var aaid userAID
	if err := query.One(&aaid); err != nil {
		return "", err
	}
	return aaid.AnonymousID, nil
}

func requireSignIn(ctx context.Context) (*User, error) {
	email, ok := ctx.Value(mw.ContextKeyEmail).(string)
	if !ok {
		return nil, fmt.Errorf("Forbidden, no access token")
	}
	return GetUserByEmail(ctx, email)
}
