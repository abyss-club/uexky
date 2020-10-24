package entity

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"gitlab.com/abyss.club/uexky/lib/config"
	"gitlab.com/abyss.club/uexky/lib/uid"
)

type UserUpdate struct {
	Name *string
	Role *Role
	Tags []string
}

type AuthInfo struct {
	UserID  uid.UID
	Email   string
	IsGuest bool
}

type UserRepo interface {
	SetCode(ctx context.Context, email string, code Code) error
	GetCodeEmail(ctx context.Context, code Code) (string, error)
	DelCode(ctx context.Context, code Code) error

	GetUserByID(ctx context.Context, id uid.UID) (*User, error)
	GetUserByAuthInfo(ctx context.Context, ai AuthInfo) (*User, error)
	SetToken(ctx context.Context, token *Token) error
	GetToken(ctx context.Context, tok string) (*Token, error)
	InsertUser(ctx context.Context, user *User) (*User, error)
	UpdateUser(ctx context.Context, user *User) (*User, error)
}

const (
	CodeLength  = 36
	CodeExpire  = 20 * time.Minute
	TokenLength = 24
	TokenExpire = 30 * time.Hour * 24
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

func (c Code) SignInURL(redirectTo string) string {
	srvCfg := &(config.Get().Server)
	return fmt.Sprintf("%s://%s/auth/?code=%s&next=%s", srvCfg.Proto, srvCfg.APIDomain, c, redirectTo)
}

type Token struct {
	Tok      string        `json:"tok,omitempty"`
	Expire   time.Duration `json:"expire,omitempty"`
	UserID   uid.UID       `json:"user_id,omitempty"`
	UserRole Role          `json:"user_role,omitempty"`
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
	ActionPubPost:     RoleGuest,
	ActionPubThread:   RoleGuest,
}
