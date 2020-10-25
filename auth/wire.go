//+build wireinject

package auth

import (
	"github.com/google/wire"
	"gitlab.com/abyss.club/uexky/adapter"
	"gitlab.com/abyss.club/uexky/lib/mail"
	"gitlab.com/abyss.club/uexky/lib/redis"
	"gitlab.com/abyss.club/uexky/mocks"
)

var mailSet = wire.NewSet(
	wire.Bind(new(adapter.MailAdapter), new(*mail.Adapter)),
	mail.NewAdapter,
)

var mockMailSet = wire.NewSet(
	wire.Bind(new(adapter.MailAdapter), new(*mocks.MailAdapter)),
	wire.Struct(new(mocks.MailAdapter), "*"),
)

var InfraSet = wire.NewSet(
	redis.NewClient,
)

var ServiceSet = wire.NewSet(
	mailSet,
	wire.Struct(new(Repo), "*"),
	wire.Struct(new(Service), "*"),
)

var MockServiceSet = wire.NewSet(
	mockMailSet,
	wire.Struct(new(Repo), "*"),
	wire.Struct(new(Service), "*"),
)

func InitAuthService() (*Service, error) {
	wire.Build(InfraSet, ServiceSet)
	return &Service{}, nil
}

func InitMockAuthService() (*Service, error) {
	wire.Build(InfraSet, MockServiceSet)
	return &Service{}, nil
}
