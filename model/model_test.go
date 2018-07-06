package model

import (
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky/mgmt"
)

const testDB = "testing"

// this file only have common test tools
func prepTestDB() {
	mgmt.LoadConfig("")
	mgmt.Config.MainTags = []string{"MainA", "MainB", "MainC"}
	mgmt.ReplaceConfigByEnv()
	if err := Init(); err != nil {
		log.Fatal(errors.Wrap(err, "Connect to test db"))
	}
	session := pkg.mongoSession.Copy()
	if err := session.DB(testDB).DropDatabase(); err != nil {
		log.Fatal(errors.Wrap(err, "tear down"))
	}
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

func TestMain(m *testing.M) {
	prepTestDB()
	addMockUser()
	os.Exit(m.Run())
}
