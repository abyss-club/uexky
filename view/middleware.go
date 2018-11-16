package view

import (
	"log"
	"net/http"

	"github.com/gomodule/redigo/redis"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky/config"
	"gitlab.com/abyss.club/uexky/model"
	"gitlab.com/abyss.club/uexky/uexky"
)

// ----------------------------------------//
// add uexky object to ctx                 //
// ----------------------------------------//

func withUexky(handle httprouter.Handle) httprouter.Handle {
	pool := uexky.InitPool()
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		ctx, done := pool.Push(req.Context())
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
		// check if signed in
		u := uexky.Pop(req.Context())
		isSignedIn := false
		tokenCookie, err := req.Cookie("token")
		log.Printf("find token cookie %v", tokenCookie)
		if err == nil {
			isSignedIn = true
		} else if err != http.ErrNoCookie { // err must be ErrNoCookie
			httpError(w, http.StatusInternalServerError, err)
			return
		}

		// refresh cookie expire
		if isSignedIn {
			cookie, err := refreshToken(u, tokenCookie.Value)
			if err != nil {
				httpError(w, http.StatusInternalServerError, err)
				return
			}
			http.SetCookie(w, cookie)
		}

		// auth
		email := ""
		if isSignedIn {
			email, err = authToken(u, tokenCookie.Value)
			if err != nil {
				http.Error(w, err.Error(), http.StatusForbidden)
				return
			}
		}
		model.NewUexkyAuth(u, email)

		// flow
		ipHeader := config.Config.RateLimit.HTTPHeader
		if ipHeader == "" {
			uexky.NewMockFlow(u)
		} else {
			ip := req.Header.Get(ipHeader)
			uexky.NewUexkyFlow(u, ip, email)
		}

		handle(w, req, p)

		// after
		w.Header().Set("Flow-Remaining", u.Flow.Remaining())
	}
}

const (
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
		Domain:   config.Config.Domain.WEB,
		MaxAge:   tokenCookieAge,
		HttpOnly: true,
	}
	if config.Config.Proto == "https" {
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
