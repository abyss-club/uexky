//+build wireinject

package cmd

import (
	"github.com/google/wire"

	"gitlab.com/abyss.club/uexky/graph"
	"gitlab.com/abyss.club/uexky/repo"
	"gitlab.com/abyss.club/uexky/service"
)

func InitResolver() graph.Resolver {
	wire.Build(service.NewService, graph.NewResolver, repo.NewRepository)
	return graph.Resolver{}
}
