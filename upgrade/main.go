package main

import (
	"flag"
	"fmt"

	"github.com/go-pg/pg/v9"
	"gitlab.com/abyss.club/uexky/lib/errors"
)

var prevDB string
var newDB string

func init() {
	flag.StringVar(&prevDB, "prev", "", "db address of prev version")
	flag.StringVar(&newDB, "new", "", "db address of new version")
}

func main() {
	flag.Parse()
	fmt.Println("prev db: ", prevDB, "new db:", newDB)
	if prevDB == "" || newDB == "" {
		panic("you should specified database")
	}
	var migrator Migrator
	var err error
	migrator.PrevDB, err = connectDB(prevDB)
	if err != nil {
		panic(err)
	}
	migrator.NewDB, err = connectDB(newDB)
	if err != nil {
		panic(err)
	}
	if err := migrator.DoMigrate(); err != nil {
		panic(err)
	}
	fmt.Println("Migrate Succeed!")
}

func connectDB(url string) (*pg.DB, error) {
	opt, err := pg.ParseURL(url)
	opt.PoolSize = 16
	if err != nil {
		return nil, errors.BadParams.Handle(err, "parse postgres uri")
	}
	return pg.Connect(opt), nil
}
