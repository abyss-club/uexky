package repo

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	red "github.com/go-redis/redis/v7"
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
	if err := postgres.RebuildDB(); err != nil {
		t.Fatal(err)
	}
	db, err := postgres.NewDB()
	if err != nil {
		t.Fatalf("get database: %v", err)
	}
	txAdapter := &postgres.TxAdapter{DB: db}
	ctx := context.Background()
	return txAdapter.AttachDB(ctx)
}
