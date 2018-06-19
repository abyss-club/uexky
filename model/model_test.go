package model

import (
	"log"
	"reflect"

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
)

const testDB = "testing"

var testMainTags = []string{"MainA", "MainB", "MainC"}

// this file only have common test tools
func dialTestDB() {
	if err := Init("localhost:27017", testDB, testMainTags); err != nil {
		log.Fatal(errors.Wrap(err, "Connect to test db"))
	}
}

func tearDown() {
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
