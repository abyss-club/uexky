package view

import (
	"net/http"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky/config"
	"gitlab.com/abyss.club/uexky/uexky"
	"gitlab.com/abyss.club/uexky/uuid64"
)

// AuthHandle ...
func AuthHandle(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	code := req.URL.Query().Get("code")
	if code == "" {
		httpError(w, http.StatusBadRequest, "缺乏必要信息")
		return
	}
	u := uexky.Pop(req.Context())
	token, err := authCode(u, code)
	if err != nil {
		httpErrorf(w, http.StatusBadRequest, "验证信息错误，或已失效。 %v", err)
		return
	}

	u.Redis.Do("DEL", code) // delete after use
	cookie := newTokenCookie(token)
	http.SetCookie(w, cookie)
	w.Header().Set("Location", config.WebURLPrefix())
	w.Header().Set("Cache-Control", "no-cache, no-store")
	w.WriteHeader(http.StatusFound)
}

// 24 charactors Base64 token
var tokenGenerator = uuid64.Generator{Sections: []uuid64.Section{
	&uuid64.RandomSection{Length: 10},
	&uuid64.CounterSection{Length: 2, Unit: time.Millisecond},
	&uuid64.TimestampSection{Length: 7, Unit: time.Millisecond},
	&uuid64.RandomSection{Length: 5},
}}

// User         Uexky
//  |--- code --->|
//  |<-- token ---|

func authCode(u *uexky.Uexky, code string) (string, error) {
	email, err := redis.String(u.Redis.Do("GET", code))
	if err == redis.ErrNil {
		return "", errors.New("Invalid code")
	} else if err != nil {
		return "", errors.Wrap(err, "Get code from redis")
	}
	token, err := tokenGenerator.New()
	if err != nil {
		return "", errors.Wrap(err, "gen token")
	}
	if _, err := u.Redis.Do("SET", token, email, "EX", tokenCookieAge); err != nil {
		return "", errors.Wrap(err, "set token to redis")
	}
	return token, nil
}

func newTokenCookie(token string) *http.Cookie {
	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		Domain:   config.Config.Domain.WEB,
		MaxAge:   tokenCookieAge,
		HttpOnly: true,
	}
	if config.Config.Proto == "https" {
		cookie.Secure = true
	}
	return cookie
}
