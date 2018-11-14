package model

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky/mgmt"
	"gitlab.com/abyss.club/uexky/mw"
)

const testDB = "testing"

var testCtx context.Context

func TestMain(m *testing.M) {
	ctx := prepTestDB()
	addMockUser(ctx)
	testCtx = ctx
	os.Exit(m.Run())
}

// this file only have common test tools
func prepTestDB() context.Context {
	mgmt.LoadConfig("")
	mgmt.Config.Mongo.DB = testDB
	mgmt.Config.MainTags = []string{"MainA", "MainB", "MainC"}
	mgmt.ReplaceConfigByEnv()

	mongo := mw.ConnectMongodb()
	if err := mongo.DB().DropDatabase(); err != nil {
		log.Fatal(errors.Wrap(err, "drop test dababase"))
	}
	ctx := context.WithValue(
		context.Background(), mw.ContextKeyMongo, mongo,
	)

	rd := mw.RedisPool.Get()
	ctx = context.WithValue(ctx, mw.ContextKeyRedis, rd)
	return ctx
}

func addMockUser(ctx context.Context) {
	log.Print("addMockUser!")
	users := []*User{
		&User{
			ID:    bson.NewObjectId(),
			Email: "0@mail.com",
			Name:  "test0",
			Tags:  []string{"动画"},
		},
		&User{
			ID:    bson.NewObjectId(),
			Email: "1@mail.com",
			Name:  "",
			Tags:  []string{},
		},
		&User{
			ID:    bson.NewObjectId(),
			Email: "2@mail.com",
			Name:  "",
			Tags:  []string{},
		},
	}

	c := u.Mongo.C(colleUser)
	for _, user := range users {
		if err := c.Insert(user); err != nil {
			log.Fatal(errors.Wrap(err, "gen mock users"))
		}
	}
	mockUsers = users
}

func ctxWithUser(u *User) context.Context {
	return context.WithValue(testCtx, mw.ContextKeyEmail, u.Email)
}

// compare functions:

var timeCmp = cmp.Comparer(func(l, r time.Time) bool {
	dur := l.Sub(r)
	fmt.Println("during", dur)
	return dur > -time.Millisecond && dur < time.Millisecond
})

func equal(x, y interface{}) bool {
	return cmp.Equal(x, y, cmpopts.EquateEmpty(), timeCmp)
}

func diff(x, y interface{}) string {
	return cmp.Diff(x, y, cmpopts.EquateEmpty(), timeCmp)
}
