//+build wireinject

package server

import (
	"github.com/google/wire"
	"gitlab.com/abyss.club/uexky/graph"
	"gitlab.com/abyss.club/uexky/uexky"
)

func InitProdServer() (*Server, error) {
	wire.Build(
		wire.Struct(new(Server), "*"),
		wire.Struct(new(graph.Resolver), "*"),
		uexky.ProdServiceSet,
	)
	return &Server{}, nil
}
