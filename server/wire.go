//+build wireinject

package server

import (
	"github.com/google/wire"
	"gitlab.com/abyss.club/uexky/auth"
	"gitlab.com/abyss.club/uexky/graph"
	"gitlab.com/abyss.club/uexky/lib/redis"
	"gitlab.com/abyss.club/uexky/uexky"
)

func InitProdServer() (*Server, error) {
	wire.Build(
		wire.Struct(new(Server), "*"),
		wire.Struct(new(graph.Resolver), "*"),
		uexky.ServiceSet,
		auth.ServiceSet,
		redis.NewClient,
	)
	return &Server{}, nil
}
