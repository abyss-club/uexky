package server

import (
	"fmt"
	"net/http"
	"time"

	"gitlab.com/abyss.club/uexky/config"
)

func (s *Server) AuthHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	code := req.URL.Query().Get("code")
	token, err := s.service().SignInByCode(req.Context(), code)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) //TODO: error type
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
