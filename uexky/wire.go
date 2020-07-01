//+build wireinject

package uexky

import (
	"github.com/google/wire"
	"gitlab.com/abyss.club/uexky/lib/mail"
	"gitlab.com/abyss.club/uexky/lib/postgres"
	"gitlab.com/abyss.club/uexky/lib/redis"
	"gitlab.com/abyss.club/uexky/mocks"
	"gitlab.com/abyss.club/uexky/repo"
	"gitlab.com/abyss.club/uexky/uexky/adapter"
	"gitlab.com/abyss.club/uexky/uexky/entity"
)

var serviceSet = wire.NewSet(
	wire.Struct(new(Service), "*"),
	wire.Struct(new(entity.ForumService), "*"),
	wire.Struct(new(entity.UserService), "*"),
	wire.Struct(new(entity.NotiService), "*"),
)

var mailSet = wire.NewSet(
	wire.Bind(new(adapter.MailAdapter), new(*mail.Adapter)),
	mail.NewAdapter,
)

var repoSet = wire.NewSet(
	wire.Struct(new(postgres.TxAdapter), "*"),
	wire.Bind(new(adapter.Tx), new(*postgres.TxAdapter)),
	postgres.NewDB,

	redis.NewClient,

	wire.Struct(new(repo.ForumRepo), "*"),
	wire.Struct(new(repo.UserRepo), "*"),
	wire.Struct(new(repo.NotiRepo), "*"),
	wire.Bind(new(entity.ForumRepo), new(*repo.ForumRepo)),
	wire.Bind(new(entity.UserRepo), new(*repo.UserRepo)),
	wire.Bind(new(entity.NotiRepo), new(*repo.NotiRepo)),
)

var mockMailSet = wire.NewSet(
	wire.Bind(new(adapter.MailAdapter), new(*mocks.MailAdapter)),
	wire.Struct(new(mocks.MailAdapter), "*"),
)

var ProdServiceSet = wire.NewSet(
	serviceSet,
	repoSet,
	mailSet,
)

func InitProdService() (*Service, error) {
	wire.Build(ProdServiceSet)
	return &Service{}, nil
}

var DevServiceSet = wire.NewSet(
	serviceSet,
	repoSet,
	mockMailSet,
)

func InitDevService() (*Service, error) {
	wire.Build(DevServiceSet)
	return &Service{}, nil
}
