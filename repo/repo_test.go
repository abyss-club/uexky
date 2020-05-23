package repo

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-pg/pg/v9"
	red "github.com/go-redis/redis/v7"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // postgres migrate
	_ "github.com/golang-migrate/migrate/v4/source/file"       // migrate file source
	"gitlab.com/abyss.club/uexky/lib/config"
	"gitlab.com/abyss.club/uexky/lib/postgres"
	"gitlab.com/abyss.club/uexky/lib/redis"
)

func TestMain(m *testing.M) {
	if err := config.Load(""); err != nil {
		log.Fatalf("load config: %v", err)
	}
	fmt.Printf("run test in config: %#v\n", config.Get())
	os.Exit(m.Run())
}

func getRedis(t *testing.T) *red.Client {
	rc, err := redis.NewClient()
	if err != nil {
		t.Fatalf("connect redis: %v", err)
	}
	return rc
}

func getNewDBCtx(t *testing.T) context.Context {
	db, err := postgres.NewDB()
	if err != nil {
		t.Fatalf("get database: %v", err)
	}
	if err := rebuildDB(db); err != nil {
		t.Fatal(err)
	}
	txAdapter := &postgres.TxAdapter{DB: db}
	ctx := context.Background()
	return txAdapter.AttachDB(ctx)
}

func rebuildDB(db *pg.DB) error {
	if _, err := db.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;"); err != nil {
		return fmt.Errorf("drop db error: %w", err)
	}
	migrateFilesPath := "../migrates"
	path, err := filepath.Abs(migrateFilesPath)
	if err != nil {
		return fmt.Errorf("parse migration file path: %w", err)
	}
	source := fmt.Sprintf("file://%s", path)
	m, err := migrate.New(source, config.Get().PostgresURI)
	if err != nil {
		return fmt.Errorf("new migration: %w", err)
	}
	if err := m.Up(); err != nil {
		return fmt.Errorf("up migration: %w", err)
	}
	return nil
}
