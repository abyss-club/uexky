package model

import (
	"context"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/globalsign/mgo/bson"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky/mgmt"
	"gitlab.com/abyss.club/uexky/mw"
)

const testDB = "testing"

var testCtx context.Context

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
	return ctx
}

var strSliceCmp = cmp.Comparer(func(l, r []string) bool {
	if len(l) == len(r) && len(l) == 0 {
		return true
	}
	return reflect.DeepEqual(l, r)
})

var threadSliceCmp = cmp.Comparer(func(l, r []*Thread) bool {
	if len(l) == len(r) && len(l) == 0 {
		return true
	}
	if len(l) != len(r) {
		return false
	}
	for i := 0; i < len(l); i++ {
		if !reflect.DeepEqual(l[i], r[i]) {
			return false
		}
	}
	return true
})

func addMockUser(ctx context.Context) {
	log.Print("addMockUser!")
	users := []*User{
		&User{bson.NewObjectId(), "0@mail.com", "test0", []string{"动画"}},
		&User{bson.NewObjectId(), "1@mail.com", "", []string{}},
		&User{bson.NewObjectId(), "2@mail.com", "", []string{}},
	}

	c := mw.GetMongo(ctx).C(colleUser)
	for _, user := range users {
		if err := c.Insert(user); err != nil {
			log.Fatal(errors.Wrap(err, "gen mock users"))
		}
	}
	mockUsers = users
}

func TestMain(m *testing.M) {
	ctx := prepTestDB()
	addMockUser(ctx)
	testCtx = ctx
	os.Exit(m.Run())
}
