package view

import (
	"log"
	"net/http"

	"github.com/gomodule/redigo/redis"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky/mgmt"
	"gitlab.com/abyss.club/uexky/model"
	"gitlab.com/abyss.club/uexky/uexky"
)

// ----------------------------------------//
// add uexky object to ctx                 //
// ----------------------------------------//

func withUexky(handle httprouter.Handle) httprouter.Handle {
	pool := uexky.InitPool()
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		ctx, done := pool.Push(req.Context(), nil, nil)
		defer done()
		req = req.WithContext(ctx)
		handle(w, req, p)
	}
}

// ----------------------------------------//
// add Auth and Flow object to uexky       //
// ----------------------------------------//

// Attach AuthInfo and Flow
func withAuthAndFlow(handle httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		tokenCookie, err := req.Cookie("token")
		ip := req.Header.Get(remoteIPHeader)
		log.Printf("find token cookie %v", tokenCookie)
		if err != nil { // err must be ErrNoCookie, non-login user, do nothing
			handle(w, req, p)
			return
		}
		u := uexky.Pop(req.Context())

		// refresh expire
		cookie, err := refreshToken(u, tokenCookie.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		http.SetCookie(w, cookie)
		email, err := authToken(u, tokenCookie.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		u.Auth = model.NewAuthInfo(u, email)
		u.Flow = uexky.NewFlow(u, ip, email)

		handle(w, req, p)
	}
}

const (
	remoteIPHeader = "Remote-IP"
	tokenCookieAge = 7 * 86400
)

func refreshToken(u *uexky.Uexky, token string) (*http.Cookie, error) {
	if _, err := u.Redis.Do("EXPIRE", token, tokenCookieAge); err != nil {
		return nil, err
	}
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
	return cookie, nil
}

// User         Uexky
//  |--- token -->|--> email
func authToken(u *uexky.Uexky, token string) (string, error) {
	email, err := redis.String(u.Redis.Do("GET", token))
	if err == redis.ErrNil {
		return "", nil
	} else if err != nil {
		return "", errors.Wrap(err, "Get token from redis")
	}
	return email, nil
}
