package model

import (
	"context"
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky/api"
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
	c, cs := Colle(colleUser)
	defer cs()
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

func isNameUesd(name string) (bool, error) {
	c, cs := Colle(colleUser)
	defer cs()
	c.EnsureIndexKey("name")

	count, err := c.Find(bson.M{"name": name}).Count()
	return count != 0, err
}

// SetName ...
func (a *User) SetName(ctx context.Context, name string) error {
	if a.Name != "" {
		return fmt.Errorf("You already have name '%v'", a.Name)
	}
	if used, err := isNameUesd(name); err != nil {
		return errors.Wrapf(err, "Check name '%s'", name)
	} else if used {
		return fmt.Errorf("This name is already in uesd")
	}

	c, cs := Colle(colleUser)
	defer cs()
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

	c, cs := Colle(colleUser)
	defer cs()
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
func (a *User) AnonymousID(threadID string, new bool) (string, error) {
	c, cs := Colle(colleAID)
	c.EnsureIndexKey("thread_id", "user_id")
	defer cs()

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
	userID, ok := ctx.Value(api.ContextKeyLoggedInUser).(bson.ObjectId)
	if !ok {
		return nil, fmt.Errorf("Forbidden, no access token")
	}
	c, cs := Colle(colleUser)
	defer cs()

	var user *User
	if err := c.FindId(userID).One(&user); err != nil {
		return nil, errors.Wrap(err, "Find user")
	}
	return user, nil
}

// SignInUser ...
func SignInUser(code string) (string, error) {
	email, err := redis.String(RedisConn.Do("GET", code))
	if err == redis.ErrNil {
		return "", errors.New("Invalid code")
	} else if err != nil {
		return "", errors.Wrap(err, "Get code from redis")
	}
	user, err := GetUserByEmail(context.Background(), email)
	if err != nil {
		return "", errors.Wrap(err, "find user")
	}
	token, err := tokenGenerator.New()
	if err != nil {
		return "", errors.Wrap(err, "gen token")
	}
	if _, err := RedisConn.Do("SET", token, user.ID.Hex(), "EX", 600); err != nil {
		return "", errors.Wrap(err, "set token to redis")
	}
	return token, nil
}
