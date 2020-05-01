//+build wireinject

package wire

import (
	"github.com/google/wire"

	"gitlab.com/abyss.club/uexky/graph"
	"gitlab.com/abyss.club/uexky/server"
	"gitlab.com/abyss.club/uexky/uexky"
	"gitlab.com/abyss.club/uexky/uexky/entity"
)

func InitServer() *server.Server {
	wire.Build(
		wire.Struct(new(server.Server), "*"),
		wire.Struct(new(graph.Resolver), "*"),
		wire.Struct(new(uexky.Service), "*"),
		wire.Struct(new(entity.UserService), "*"),
		wire.Struct(new(entity.ForumService), "*"),
		wire.Struct(new(entity.NotiService), "*"),
	)
	return &server.Server{}
}

// func InitResolver() graph.Resolver {
// 	wire.Build(service.NewService, graph.NewResolver, repo.NewRepository)
// 	return graph.Resolver{}
// }
