package server

import (
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"gitlab.com/abyss.club/uexky/graph/generated"
	"gitlab.com/abyss.club/uexky/lib/config"
)

func (s *Server) AuthHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	code := req.URL.Query().Get("code")
	token, err := s.Service.SignInByCode(req.Context(), code)
	if err != nil {
		writeError(w, err)
		return
	}
	location := fmt.Sprintf("%s://%s", config.Get().Server.Proto, config.Get().Server.Domain)
	http.SetCookie(w, token.Cookie())
	w.Header().Set("Location", location)
	w.Header().Set("Cache-Control", "no-cache, no-store")
	w.WriteHeader(http.StatusFound)
}

func (s *Server) GraphQLHandler() http.Handler {
	return handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: s.Resolver,
	}))
}
