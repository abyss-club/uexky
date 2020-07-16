package postgres

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // postgres migrate
	_ "github.com/golang-migrate/migrate/v4/source/file"       // migrate file source
	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky/lib/config"
	"gitlab.com/abyss.club/uexky/lib/uerr"
)

// RebuildDB rebuild
func RebuildDB() error {
	if config.GetEnv() == config.ProdEnv {
		return uerr.New(uerr.ParamsError, "can't rebuild db in prod env")
	}
	db, err := NewDB()
	if err != nil {
		return errors.Wrap(err, "RebuildDB")
	}
	defer db.Close()
	if _, err := db.Exec("DROP SCHEMA IF EXISTS public CASCADE; CREATE SCHEMA public;"); err != nil {
		return uerr.Wrap(uerr.PostgresError, err, "drop db error")
	}
	m, err := GetMigrate()
	if err != nil {
		return errors.Wrap(err, "RebuildDB")
	}
	if err := m.Up(); err != nil {
		return uerr.Wrap(uerr.PostgresError, err, "migrate up")
	}
	return nil
}

func GetMigrate() (*migrate.Migrate, error) {
	source := fmt.Sprintf("file://%s", config.Get().MigrationFiles)
	m, err := migrate.New(source, config.Get().PostgresURI)
	return m, uerr.Wrap(uerr.PostgresError, err, "get migration helper")
}
