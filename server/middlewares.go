package server

import (
	"context"
	"net/http"

	"gitlab.com/abyss.club/uexky/uexky"
)

// func (s *Server) withLog(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 	})
// }

func (s *Server) withUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie("token")
		if err != nil && err != http.ErrNoCookie {
			writeError(w, err)
			return
		}
		var tok string
		if tokenCookie != nil {
			tok = tokenCookie.Value
		}
		token, err := s.Resolver.Auth.GetToken(r.Context(), tok)
		if err != nil {
			writeError(w, err)
		}

		if token != nil {
			var ctx context.Context
			if token.User.IsGuest {
				ctx, err = s.Resolver.Uexky.AttachGuestUserToCtx(r.Context(), token.User.UserID)
			} else {
				ctx, err = s.Resolver.Uexky.AttachEmailUserToCtx(r.Context(), token.User.Email)
			}
			if err != nil {
				writeError(w, err)
				return
			}
			r = r.WithContext(ctx)
			http.SetCookie(w, token.Cookie())
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) withDB(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := s.TxAdapter.AttachDB(r.Context())
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (s *Server) withLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := uexky.AttachLimiter(r.Context(), 50)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
