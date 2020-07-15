package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"gitlab.com/abyss.club/uexky/graph/generated"
	"gitlab.com/abyss.club/uexky/lib/config"
	"gitlab.com/abyss.club/uexky/lib/uerr"
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
	server := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: s.Resolver,
	}))
	server.SetErrorPresenter(func(ctx context.Context, err error) *gqlerror.Error {
		path := graphql.GetFieldContext(ctx).Path()
		message := err.Error()
		var code string
		switch {
		case errors.Is(err, uerr.New(uerr.ParamsError)):
			code = uerr.ParamsError.Code()
		case errors.Is(err, uerr.New(uerr.AuthError)):
			code = uerr.AuthError.Code()
		case errors.Is(err, uerr.New(uerr.PermissionError)):
			code = uerr.PermissionError.Code()
		case errors.Is(err, uerr.New(uerr.NotFoundError)):
			code = uerr.NotFoundError.Code()
		case errors.Is(err, uerr.New(uerr.DBError)):
			code = uerr.InternalError.Code()
		case errors.Is(err, uerr.New(uerr.InternalError)):
			code = uerr.InternalError.Code()
		default:
			code = uerr.InternalError.Code()
		}
		gerr := gqlerror.ErrorPathf(path, message)
		gerr.Extensions = map[string]interface{}{
			"code":       code,
			"stacktrace": strings.Split(fmt.Sprintf("%+v", err), "\n"),
		}
		return gerr
	})
	return server
}
