package auth

import (
	"fmt"
	"net/http"
	"time"

	"gitlab.com/abyss.club/uexky/lib/config"
	"gitlab.com/abyss.club/uexky/lib/uid"
)

type Code string

func (c Code) SignInURL(redirectTo string) string {
	srvCfg := &(config.Get().Server)
	return fmt.Sprintf("%s://%s/auth/?code=%s&next=%s", srvCfg.Proto, srvCfg.APIDomain, c, redirectTo)
}

const (
	CodeLength  = 36
	CodeExpire  = 20 * time.Minute
	TokenLength = 24
	TokenExpire = 30 * time.Hour * 24
)

type Token struct {
	Tok  string   `json:"tok,omitempty"`
	User UserInfo `json:"user"`
}

type UserInfo struct {
	UserID  uid.UID `json:"user_id,omitempty"`
	Email   string  `json:"email,omitempty"`
	IsGuest bool    `json:"is_guest"`
}

func generateTok() string {
	return uid.RandomBase64Str(TokenLength)
}

func NewEmailToken(email string) *Token {
	return &Token{
		Tok: generateTok(),
		User: UserInfo{
			Email:   email,
			IsGuest: false,
		},
	}
}

func NewGuestToken() *Token {
	return &Token{
		Tok: generateTok(),
		User: UserInfo{
			UserID:  uid.NewUID(),
			IsGuest: true,
		},
	}
}

func (t Token) Cookie() *http.Cookie {
	cookie := &http.Cookie{
		Name:     "token",
		Value:    t.Tok,
		Path:     "/",
		MaxAge:   int(TokenExpire / time.Second),
		Domain:   config.Get().Server.Domain,
		HttpOnly: true,
	}
	if config.Get().Server.Proto == "https" {
		cookie.Secure = true
	}
	return cookie
}
