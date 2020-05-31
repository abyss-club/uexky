package postgres

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // postgres migrate
	_ "github.com/golang-migrate/migrate/v4/source/file"       // migrate file source
	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky/lib/config"
)

// RebuildDB rebuild
func RebuildDB() error {
	if config.GetEnv() == config.ProdEnv {
		return errors.New("can't rebuild db in prod env")
	}
	db, err := NewDB()
	if err != nil {
		return errors.Wrap(err, "get database")
	}
	defer db.Close()
	if _, err := db.Exec("DROP SCHEMA IF EXISTS public CASCADE; CREATE SCHEMA public;"); err != nil {
		return errors.Wrap(err, "drop db error")
	}
	m, err := GetMigrate()
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil {
		return errors.Wrap(err, "up migration")
	}
	return nil
}

func GetMigrate() (*migrate.Migrate, error) {
	source := fmt.Sprintf("file://%s", config.Get().MigrateFiles)
	m, err := migrate.New(source, config.Get().PostgresURI)
	return m, errors.Wrap(err, "new migration")
}
