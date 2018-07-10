package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/garyburd/redigo/redis"
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
	if _, err := rd.Do("SET", token, email, "EX", 600); err != nil {
		return "", errors.Wrap(err, "set token to redis")
	}
	return token, nil
}

// AuthHandle ...
func AuthHandle(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	code := req.URL.Query().Get("code")
	if code == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("缺乏必要信息"))
		return
	}
	token, err := authCode(req.Context(), code)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("验证信息错误，或已失效。 %v", err)))
		return
	}

	GetRedis(req.Context()).Do("DEL", code) // delete after use
	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		Domain:   mgmt.Config.Domain.WEB,
		MaxAge:   86400,
		HttpOnly: true,
	}
	if mgmt.Config.Proto == "https" {
		cookie.Secure = true
	}
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
		if err != nil { // err must be ErrNoCookie,  non-login user, do noting
			handle(w, req, p)
			return
		}

		email, err := authToken(req.Context(), tokenCookie.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		if email != "" {
			log.Printf("Logged user %s", email)
			ctx := context.WithValue(req.Context(), ContextKeyEmail, email)
			req = req.WithContext(ctx)
		}
		handle(w, req, p)
	}
}
