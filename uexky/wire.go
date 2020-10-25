//+build wireinject

package uexky

import (
	"github.com/google/wire"
	"gitlab.com/abyss.club/uexky/lib/postgres"
	"gitlab.com/abyss.club/uexky/lib/redis"
	"gitlab.com/abyss.club/uexky/repo"
	"gitlab.com/abyss.club/uexky/uexky/adapter"
)

var repoSet = wire.NewSet(
	wire.Struct(new(postgres.TxAdapter), "*"),
	wire.Bind(new(adapter.Tx), new(*postgres.TxAdapter)),
	postgres.NewDB,
	redis.NewClient,
	repo.NewRepo,
)

var ServiceSet = wire.NewSet(
	repoSet,
	NewService,
)

func InitUexkyService() (*Service, error) {
	wire.Build(ServiceSet)
	return &Service{}, nil
}
