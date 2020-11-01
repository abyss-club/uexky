//+build wireinject

package uexky

import (
	"github.com/google/wire"
	"gitlab.com/abyss.club/uexky/adapter"
	"gitlab.com/abyss.club/uexky/lib/postgres"
	"gitlab.com/abyss.club/uexky/lib/redis"
	"gitlab.com/abyss.club/uexky/uexky/repo"
)

var repoSet = wire.NewSet(
	wire.Struct(new(postgres.TxAdapter), "*"),
	wire.Bind(new(adapter.Tx), new(*postgres.TxAdapter)),
	postgres.NewDB,
	repo.NewRepo,
)

var InfraSet = wire.NewSet(
	redis.NewClient,
)

var ServiceSet = wire.NewSet(
	repoSet,
	NewService,
)

func InitUexkyService() (*Service, error) {
	wire.Build(InfraSet, ServiceSet)
	return &Service{}, nil
}
