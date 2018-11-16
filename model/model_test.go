package model

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky/config"
	"gitlab.com/abyss.club/uexky/uexky"
)

const testDB = "testing"

var mockUsers []*User
var uexkyPool *uexky.Pool
var mu []*uexky.Uexky // mock uexky

func TestMain(m *testing.M) {
	prepTestDB()
	addMockUser()
	os.Exit(m.Run())
}

// this file only have common test tools
func prepTestDB() {
	config.LoadConfig("")
	config.Config.Mongo.DB = testDB
	config.Config.MainTags = []string{"MainA", "MainB", "MainC"}
	config.ReplaceConfigByEnv()

	uexkyPool = uexky.InitPool()
	u := uexkyPool.NewUexky()
	defer u.Close()

	if err := u.Mongo.DB().DropDatabase(); err != nil {
		log.Fatal(errors.Wrap(err, "drop test dababase"))
	}
}

func addMockUser() {
	uTemp := uexkyPool.NewUexky()
	defer uTemp.Close()
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

	c := uTemp.Mongo.C(colleUser)
	for _, user := range users {
		if err := c.Insert(user); err != nil {
			log.Fatal(errors.Wrap(err, "gen mock users"))
		}
	}
	mockUsers = users
	mu = []*uexky.Uexky{}
	for _, user := range users {
		u := uexkyPool.NewUexky()
		NewUexkyAuth(u, user.Email)
		uexky.NewMockFlow(u)
		mu = append(mu, u)
	}
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
