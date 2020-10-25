// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package uexky

import (
	"github.com/google/wire"
	"gitlab.com/abyss.club/uexky/lib/postgres"
	"gitlab.com/abyss.club/uexky/lib/redis"
	"gitlab.com/abyss.club/uexky/repo"
	"gitlab.com/abyss.club/uexky/uexky/adapter"
)

// Injectors from wire.go:

func InitUexkyService() (*Service, error) {
	db, err := postgres.NewDB()
	if err != nil {
		return nil, err
	}
	txAdapter := &postgres.TxAdapter{
		DB: db,
	}
	client, err := redis.NewClient()
	if err != nil {
		return nil, err
	}
	entityRepo := repo.NewRepo(client)
	service, err := NewService(txAdapter, entityRepo)
	if err != nil {
		return nil, err
	}
	return service, nil
}

// wire.go:

var repoSet = wire.NewSet(wire.Struct(new(postgres.TxAdapter), "*"), wire.Bind(new(adapter.Tx), new(*postgres.TxAdapter)), postgres.NewDB, redis.NewClient, repo.NewRepo)

var ServiceSet = wire.NewSet(
	repoSet,
	NewService,
)
