package postgres

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // postgres migrate
	_ "github.com/golang-migrate/migrate/v4/source/file"       // migrate file source
	"gitlab.com/abyss.club/uexky/lib/config"
	"gitlab.com/abyss.club/uexky/lib/errors"
)

// RebuildDB rebuild
func RebuildDB() error {
	if config.GetEnv() == config.ProdEnv {
		return errors.BadParams.New("can't rebuild db in prod env")
	}
	db, err := NewDB()
	if err != nil {
		return errors.Wrap(err, "RebuildDB")
	}
	defer db.Close()
	if _, err := db.Exec("DROP SCHEMA IF EXISTS public CASCADE; CREATE SCHEMA public;"); err != nil {
		return errors.Postgres.Handle(err, "drop db error")
	}
	m, err := GetMigrate()
	if err != nil {
		return errors.Wrap(err, "RebuildDB")
	}
	if err := m.Up(); err != nil {
		return errors.Postgres.Handle(err, "migrate up")
	}
	return nil
}

func GetMigrate() (*migrate.Migrate, error) {
	source := fmt.Sprintf("file://%s", config.Get().MigrationFiles)
	m, err := migrate.New(source, config.Get().PostgresURI)
	return m, errors.Postgres.Handle(err, "get migration helper")
}
