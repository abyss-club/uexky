package mw

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky/mgmt"
	"gitlab.com/abyss.club/uexky/uuid64"
)

// 24 charactors Base64 token
var tokenGenerator = uuid64.Generator{Sections: []uuid64.Section{
	&uuid64.RandomSection{Length: 10},
	&uuid64.CounterSection{Length: 2, Unit: time.Millisecond},
	&uuid64.TimestampSection{Length: 7, Unit: time.Millisecond},
	&uuid64.RandomSection{Length: 5},
}}

const tokenCookieAge = 7 * 86400

func newTokenCookie(token string) *http.Cookie {
	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		Domain:   mgmt.Config.Domain.WEB,
		MaxAge:   tokenCookieAge,
		HttpOnly: true,
	}
	if mgmt.Config.Proto == "https" {
		cookie.Secure = true
	}
	return cookie
}

// User         Uexky
//  |--- code --->|
//  |<-- token ---|

func authCode(ctx context.Context, code string) (string, error) {
	rd := GetRedis(ctx)
	email, err := redis.String(rd.Do("GET", code))
	if err == redis.ErrNil {
		return "", errors.New("Invalid code")
	} else if err != nil {
		return "", errors.Wrap(err, "Get code from redis")
	}
	token, err := tokenGenerator.New()
	if err != nil {
		return "", errors.Wrap(err, "gen token")
	}
	if _, err := rd.Do("SET", token, email, "EX", tokenCookieAge); err != nil {
		return "", errors.Wrap(err, "set token to redis")
	}
	return token, nil
}

func refreshToken(ctx context.Context, token string) error {
	rd := GetRedis(ctx)
	_, err := rd.Do("EXPIRE", token, tokenCookieAge)
	return err
}

// AuthHandle ...
func AuthHandle(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	code := req.URL.Query().Get("code")
	if code == "" {
		httpError(w, http.StatusBadRequest, "缺乏必要信息")
		return
	}
	token, err := authCode(req.Context(), code)
	if err != nil {
		httpErrorf(w, http.StatusBadRequest, "验证信息错误，或已失效。 %v", err)
		return
	}

	GetRedis(req.Context()).Do("DEL", code) // delete after use
	cookie := newTokenCookie(token)
	http.SetCookie(w, cookie)
	w.Header().Set("Location", mgmt.WebURLPrefix())
	w.Header().Set("Cache-Control", "no-cache, no-store")
	w.WriteHeader(http.StatusFound)
}

// User         Uexky
//  |--- token -->|

func authToken(ctx context.Context, token string) (string, error) {
	email, err := redis.String(GetRedis(ctx).Do("GET", token))
	if err == redis.ErrNil {
		return "", nil
	} else if err != nil {
		return "", errors.Wrap(err, "Get token from redis")
	}
	return email, nil
}

// WithAuth ...
func WithAuth(handle httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		tokenCookie, err := req.Cookie("token")
		log.Printf("find token cookie %v", tokenCookie)
		if err != nil { // err must be ErrNoCookie, non-login user, do nothing
			handle(w, req, p)
			return
		}
		// refresh expire
		refreshToken(req.Context(), tokenCookie.Value)
		cookie := newTokenCookie(tokenCookie.Value)
		http.SetCookie(w, cookie)

		email, err := authToken(req.Context(), tokenCookie.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		if email != "" {
			log.Printf("Logged user %s", email)
			req = reqWithValue(req, ContextKeyEmail, email)
		}

		handle(w, req, p)
	}
}
