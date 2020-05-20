package server

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"gitlab.com/abyss.club/uexky/graph"
	"gitlab.com/abyss.club/uexky/graph/generated"
	"gitlab.com/abyss.club/uexky/uexky"
)

type Server struct {
	Resolver *graph.Resolver
}

func (s *Server) service() *uexky.Service {
	return s.Resolver.Service
}

func (s *Server) Run() error {
	port := "8000"
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: s.Resolver,
	}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)
	http.HandleFunc("/auth", s.AuthHandler)
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	return http.ListenAndServe(":"+port, nil)
}
