//+build wireinject

package wire

import (
	"github.com/google/wire"
	"gitlab.com/abyss.club/uexky/graph"
	"gitlab.com/abyss.club/uexky/lib/mail"
	"gitlab.com/abyss.club/uexky/lib/postgres"
	"gitlab.com/abyss.club/uexky/lib/redis"
	"gitlab.com/abyss.club/uexky/repo"
	"gitlab.com/abyss.club/uexky/server"
	"gitlab.com/abyss.club/uexky/uexky"
	"gitlab.com/abyss.club/uexky/uexky/adapter"
	"gitlab.com/abyss.club/uexky/uexky/entity"
)

var prodRepoSet = wire.NewSet(
	wire.Struct(new(postgres.TxAdapter), "*"),
	wire.Bind(new(adapter.Tx), new(*postgres.TxAdapter)),
	postgres.NewDB,

	redis.NewClient,

	wire.Bind(new(adapter.MailAdapter), new(*mail.Adapter)),
	mail.NewAdapter,

	wire.Struct(new(repo.ForumRepo), "*"),
	wire.Struct(new(repo.UserRepo), "*"),
	wire.Struct(new(repo.NotiRepo), "*"),
	wire.Bind(new(entity.ForumRepo), new(*repo.ForumRepo)),
	wire.Bind(new(entity.UserRepo), new(*repo.UserRepo)),
	wire.Bind(new(entity.NotiRepo), new(*repo.NotiRepo)),
)

func InitServer() (*server.Server, error) {
	wire.Build(
		wire.Struct(new(server.Server), "*"),
		wire.Struct(new(graph.Resolver), "*"),
		wire.Struct(new(uexky.Service), "*"),
		wire.Struct(new(entity.ForumService), "*"),
		wire.Struct(new(entity.UserService), "*"),
		wire.Struct(new(entity.NotiService), "*"),
		prodRepoSet,
	)
	return &server.Server{}, nil
}
