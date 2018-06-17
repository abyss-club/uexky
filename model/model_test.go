package model

import (
	"log"

	"github.com/pkg/errors"
)

const testDB = "testing"

// this file only have common test tools
func dialTestDB() {
	if err := Dial("localhost:27017", testDB); err != nil {
		log.Fatal(errors.Wrap(err, "Connect to test db"))
	}
}

func tearDown() {
	session := mongoSession.Copy()
	if err := session.DB(testDB).DropDatabase(); err != nil {
		log.Fatal(errors.Wrap(err, "tear down"))
	}
}
