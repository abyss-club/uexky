//+build wireinject

package auth

import (
	"github.com/google/wire"
	"gitlab.com/abyss.club/uexky/lib/mail"
	"gitlab.com/abyss.club/uexky/mocks"
	"gitlab.com/abyss.club/uexky/repo"
	"gitlab.com/abyss.club/uexky/uexky"
	"gitlab.com/abyss.club/uexky/uexky/adapter"
	"gitlab.com/abyss.club/uexky/uexky/entity"
)

var mailSet = wire.NewSet(
	wire.Bind(new(adapter.MailAdapter), new(*mail.Adapter)),
	mail.NewAdapter,
)

var mockMailSet = wire.NewSet(
	wire.Bind(new(adapter.MailAdapter), new(*mocks.MailAdapter)),
	wire.Struct(new(mocks.MailAdapter), "*"),
)

var repoSet = wire.NewSet(
	wire.Bind(new(Repo), new(*RepoImpl)),
	wire.Struct(new(RepoImpl), "*"),
	wire.Bind(new(entity.UserRepo), new(*repo.UserRepo)),
	wire.Struct(new(repo.UserRepo), "*"),
	// redis.NewClient,
	wire.Struct(new(R), "*"),
)

var ServiceSet = wire.NewSet(
	repoSet,
	mailSet,
	uexky.ServiceSet,
	wire.Struct(new(Service), "*"),
)

var MockServiceSet = wire.NewSet(
	repoSet,
	mockMailSet,
	uexky.ServiceSet,
	wire.Struct(new(Service), "*"),
)

func InitAuthService() (*Service, error) {
	wire.Build(ServiceSet)
	return &Service{}, nil
}

func InitMockAuthService() (*Service, error) {
	wire.Build(MockServiceSet)
	return &Service{}, nil
}
