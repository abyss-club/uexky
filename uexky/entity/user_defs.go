package entity

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"gitlab.com/abyss.club/uexky/lib/config"
)

type UserUpdate struct {
	Name *string
	Role *Role
	Tags []string
}

type UserRepo interface {
	SetCode(ctx context.Context, email string, code string, ex time.Duration) error
	GetCodeEmail(ctx context.Context, code string) (string, error)
	DelCode(ctx context.Context, code string) error
	SetToken(ctx context.Context, email string, tok string, ex time.Duration) error
	GetTokenEmail(ctx context.Context, tok string) (string, error)

	GetOrInsertUser(ctx context.Context, email string) (user *User, isNew bool, err error)
	UpdateUser(ctx context.Context, id int64, update *UserUpdate) error
}

const (
	codeLength  = 36
	codeExpire  = 20 * time.Minute
	tokenLength = 24
	tokenExpire = 30 * time.Hour * 24
)

const authEmailHTML = `<html>
	<head>
		<meta charset="utf-8">
		<title>点击登入 Abyss!</title>
	</head>
	<body>
		<p>点击 <a href="%s">此链接</a> 进入 Abyss</p>
	</body>
</html>
`

type Code string

func (c Code) SignInURL() string {
	srvCfg := &(config.Get().Server)
	return fmt.Sprintf("%s://%s/auth/?code=%s", srvCfg.Proto, srvCfg.APIDomain, c)
}

type Token struct {
	Tok    string
	Expire time.Duration
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

type contextKey int

const (
	userKey contextKey = 1 + iota
)

func ParseRole(s string) Role {
	if s == "" {
		return RoleNormal
	}
	r := Role(s)
	if !r.IsValid() {
		return RoleBanned
	}
	return r
}

func (r Role) Value() int {
	switch r {
	case RoleAdmin:
		return 100
	case RoleMod:
		return 10
	case RoleNormal:
		return 1
	case RoleGuest:
		return 0
	case RoleBanned:
		return -1
	default:
		return -10
	}
}

type Action string

const (
	ActionProfile     = Action("PROFILE")
	ActionBanUser     = Action("BAN_USER")
	ActionPromoteUser = Action("PROMOTE_USER")
	ActionBlockPost   = Action("BLOCK_POST")
	ActionLockThread  = Action("LOCK_THREAD")
	ActionBlockThread = Action("BLOCK_THREAD")
	ActionEditTag     = Action("EDIT_TAG")
	ActionEditSetting = Action("EDIT_SETTING")
	ActionPubPost     = Action("PUB_POST")
	ActionPubThread   = Action("PUB_THREAD")
)

var ActionRole = map[Action]Role{
	ActionProfile:     RoleBanned, // Because a user can only read the profile own by himself.
	ActionBanUser:     RoleMod,
	ActionPromoteUser: RoleAdmin,
	ActionBlockPost:   RoleMod,
	ActionLockThread:  RoleMod,
	ActionBlockThread: RoleMod,
	ActionEditTag:     RoleMod,
	ActionEditSetting: RoleAdmin,
}
