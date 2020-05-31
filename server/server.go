package server

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/playground"
	"gitlab.com/abyss.club/uexky/graph"
	"gitlab.com/abyss.club/uexky/uexky"
	"gitlab.com/abyss.club/uexky/uexky/adapter"
)

type Server struct {
	Resolver  *graph.Resolver
	Service   *uexky.Service
	TxAdapter adapter.Tx
}

func (s *Server) Run() error {
	port := "8000"
	http.Handle("/", s.withDB(s.withUser(playground.Handler("GraphQL playground", "/query"))))
	http.Handle("/query", s.withDB(s.withUser(s.GraphQLHandler())))
	http.Handle("/auth", http.HandlerFunc(s.AuthHandler))
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	return http.ListenAndServe(":"+port, nil)
}
