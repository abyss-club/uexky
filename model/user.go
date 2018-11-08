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
	&uuid64.RandomSection{Length: 1},
	&uuid64.CounterSection{Length: 2, Unit: time.Second},
	&uuid64.TimestampSection{Length: 6, Unit: time.Second, NoPadding: true},
}}

const (
	nameLimit = 5
	tagLimit  = 15
	// ContextKeyUser for logged in user
	ContextKeyUser = mw.ContextKey("user")
)

// User for uexky
type User struct {
	ID           bson.ObjectId `bson:"_id"`
	Email        string        `bson:"email"`
	Name         string        `bson:"name"`
	Tags         []string      `bson:"tags"`
	ReadNotiTime struct {
		System  time.Time `bson:"system"`
		Replied time.Time `bson:"replied"`
		Quoted  time.Time `bson:"quoted"`
	} `bson:"read_noti_time"`
}

// GetUser by id (in context)
func GetUser(ctx context.Context) (*User, error) {
	// TODO: duplicate
	user, err := requireSignIn(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByEmail ...
func GetUserByEmail(ctx context.Context, email string) (*User, error) {
	if err := mw.FlowCostQuery(ctx, 1); err != nil {
		return nil, err
	}
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

func isNameUsed(ctx context.Context, name string) (bool, error) {
	if err := mw.FlowCostQuery(ctx, 1); err != nil {
		return false, err
	}
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
	if used, err := isNameUsed(ctx, name); err != nil {
		return errors.Wrapf(err, "Check name '%s'", name)
	} else if used {
		return errors.New("This name is already in use")
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

// AddSubbedTags ...
func (a *User) AddSubbedTags(ctx context.Context, tags []string) error {
	c := mw.GetMongo(ctx).C(colleUser)
	if err := c.Update(bson.M{"_id": a.ID}, bson.M{"$addToSet": bson.M{
		"tags": bson.M{"$each": tags},
	}}); err != nil {
		return err
	}
	var user *User
	if err := c.FindId(a.ID).One(&user); err != nil {
		return err
	}
	a.Tags = user.Tags
	return nil
}

// DelSubbedTags ...
func (a *User) DelSubbedTags(ctx context.Context, tags []string) error {
	c := mw.GetMongo(ctx).C(colleUser)
	if err := c.Update(bson.M{"_id": a.ID}, bson.M{"$pull": bson.M{
		"tags": bson.M{"$in": tags},
	}}); err != nil {
		return err
	}
	var user *User
	if err := c.FindId(a.ID).One(&user); err != nil {
		return err
	}
	a.Tags = user.Tags
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
	c := mw.GetMongo(ctx).C(colleAID)
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

func (a *User) getReadNotiTime(t NotiType) time.Time {
	switch t {
	case NotiTypeSystem:
		return a.ReadNotiTime.System
	case NotiTypeReplied:
		return a.ReadNotiTime.Replied
	case NotiTypeQuoted:
		return a.ReadNotiTime.Quoted
	default:
		return time.Unix(0, 0)
	}
}

func (a *User) setReadNotiTime(ctx context.Context, t NotiType, time time.Time) error {
	c := mw.GetMongo(ctx).C(colleUser)
	if err := c.Update(bson.M{"_id": a.ID}, bson.M{"$set": bson.M{
		fmt.Sprintf("read_noti_time.%v", t): time,
	}}); err != nil {
		return err
	}

	switch t {
	case NotiTypeSystem:
		a.ReadNotiTime.System = time
	case NotiTypeReplied:
		a.ReadNotiTime.Replied = time
	case NotiTypeQuoted:
		a.ReadNotiTime.Quoted = time
	default:
		panic("Invalidate Notification Type")
	}
	return nil
}

func requireSignIn(ctx context.Context) (*User, error) {
	email, ok := ctx.Value(mw.ContextKeyEmail).(string)
	if !ok {
		return nil, fmt.Errorf("Forbidden, no access token")
	}
	return GetUserByEmail(ctx, email)
}
