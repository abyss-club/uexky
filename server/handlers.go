package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"gitlab.com/abyss.club/uexky/graph/generated"
	"gitlab.com/abyss.club/uexky/lib/config"
)

func (s *Server) AuthHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	code := req.URL.Query().Get("code")
	token, err := s.service().SignInByCode(req.Context(), code)
	if err != nil {
		writeError(w, err)
		return
	}
	cookie := &http.Cookie{
		Name:     "token",
		Value:    token.Tok,
		Path:     "/",
		MaxAge:   int(token.Expire / time.Second),
		Domain:   config.Get().Server.Domain,
		HttpOnly: true,
	}
	if config.Get().Server.Proto == "https" {
		cookie.Secure = true
	}
	location := fmt.Sprintf("%s://%s", config.Get().Server.Proto, config.Get().Server.Domain)
	http.SetCookie(w, cookie)
	w.Header().Set("Location", location)
	w.Header().Set("Cache-Control", "no-cache, no-store")
	w.WriteHeader(http.StatusFound)
}

func (s *Server) GraphQLHandler() http.Handler {
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: s.Resolver,
	}))
	srv.SetErrorPresenter(func(ctx context.Context, err error) *gqlerror.Error {
		return nil
	})
	return srv
}
