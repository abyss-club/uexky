package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/playground"
	"gitlab.com/abyss.club/uexky/graph"
	"gitlab.com/abyss.club/uexky/lib/config"
	"gitlab.com/abyss.club/uexky/uexky"
	"gitlab.com/abyss.club/uexky/uexky/adapter"
)

type Server struct {
	Resolver  *graph.Resolver
	Service   *uexky.Service
	TxAdapter adapter.Tx
}

func (s *Server) Run() error {
	srvCfg := config.Get().Server
	addr := fmt.Sprintf("%s:%v", srvCfg.Host, srvCfg.Port)
	http.Handle("/", s.withDB(s.withUser(playground.Handler("GraphQL playground", "/graphql"))))
	http.Handle("/graphql", s.withDB(s.withUser(s.withLimiter(s.GraphQLHandler()))))
	http.Handle("/auth/", http.HandlerFunc(s.AuthHandler))
	log.Printf("connect to http://%s/ for GraphQL playground", addr)
	return http.ListenAndServe(addr, nil)
}
