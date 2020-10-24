package auth

import (
	"fmt"
	"net/http"
	"time"

	"gitlab.com/abyss.club/uexky/lib/config"
	"gitlab.com/abyss.club/uexky/lib/uid"
	"gitlab.com/abyss.club/uexky/uexky/entity"
)

type AuthInfo struct {
	UserID  uid.UID
	Email   string
	IsGuest bool
}

const (
	CodeLength  = 36
	CodeExpire  = 20 * time.Minute
	TokenLength = 24
	TokenExpire = 30 * time.Hour * 24
)

type Code string

func (c Code) SignInURL(redirectTo string) string {
	srvCfg := &(config.Get().Server)
	return fmt.Sprintf("%s://%s/auth/?code=%s&next=%s", srvCfg.Proto, srvCfg.APIDomain, c, redirectTo)
}

type Token struct {
	Tok      string        `json:"tok,omitempty"`
	Expire   time.Duration `json:"expire,omitempty"`
	UserID   uid.UID       `json:"user_id,omitempty"`
	UserRole entity.Role   `json:"user_role,omitempty"`
}

func (t Token) Cookie() *http.Cookie {
	cookie := &http.Cookie{
		Name:     "token",
		Value:    t.Tok,
		Path:     "/",
		MaxAge:   int(t.Expire / time.Second),
		Domain:   config.Get().Server.Domain,
		HttpOnly: true,
	}
	if config.Get().Server.Proto == "https" {
		cookie.Secure = true
	}
	return cookie
}

type TokenMsg struct {
	Email      *string `json:"email,omitempty"`
	UnsignedID *int64  `json:"unsigned_id,omitempty"`
}
