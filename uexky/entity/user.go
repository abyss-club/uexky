// user aggragate: user

package entity

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"gitlab.com/abyss.club/uexky/lib/algo"
	"gitlab.com/abyss.club/uexky/lib/config"
	"gitlab.com/abyss.club/uexky/lib/uerr"
	"gitlab.com/abyss.club/uexky/lib/uid"
	"gitlab.com/abyss.club/uexky/uexky/adapter"
)

type UserService struct {
	Repo UserRepo
	Mail adapter.MailAdapter
}

func (s *UserService) RequirePermission(ctx context.Context, action Action) (*User, error) {
	user, ok := ctx.Value(userKey).(*User)
	if !ok {
		return nil, errors.New("permission denied, no user found")
	}
	if err := user.RequirePermission(action); err != nil {
		return nil, err
	}
	return user, nil
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

func newAuthMail(email string, code Code) *adapter.Mail {
	srvCfg := &(config.Get().Server)
	authURL := code.SignInURL()
	return &adapter.Mail{
		From:    fmt.Sprintf("auth@%s", srvCfg.Domain),
		To:      email,
		Subject: "点击登入 Abyss!",
		Text:    fmt.Sprintf("点击此链接进入 Abyss：%s", authURL),
		HTML:    fmt.Sprintf(authEmailHTML, authURL),
	}
}

func (s *UserService) TrySignInByEmail(ctx context.Context, email string) (Code, error) {
	// TODO: validate email
	code := Code(uid.RandomBase64Str(codeLength))
	if err := s.Repo.SetCode(ctx, email, string(code), codeExpire); err != nil {
		return "", err
	}
	mail := newAuthMail(email, code)
	if err := s.Mail.SendEmail(ctx, mail); err != nil {
		return "", err
	}
	return code, nil
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

func (s *UserService) SignInByCode(ctx context.Context, code string) (Token, error) {
	email, err := s.Repo.GetCodeEmail(ctx, code)
	if err != nil {
		return Token{}, err
	}
	tok := uid.RandomBase64Str(tokenLength)
	if err := s.Repo.SetToken(ctx, email, tok, tokenExpire); err != nil {
		return Token{}, err
	}
	if err := s.Repo.DelCode(ctx, code); err != nil {
		log.Error(err)
	}
	return Token{Tok: tok, Expire: tokenExpire}, nil
}

type contextKey int

const (
	userKey contextKey = 1 + iota
)

func (s *UserService) CtxWithUserByToken(ctx context.Context, tok string) (context.Context, error) {
	email, err := s.Repo.GetTokenEmail(ctx, tok)
	if err != nil {
		return nil, err
	}
	if email == "" {
		user := &User{
			Role: RoleGuest,
		}
		return context.WithValue(ctx, userKey, user), nil
	}
	user, err := s.Repo.GetOrInsertUser(ctx, email)
	if err != nil {
		return nil, err
	}
	return context.WithValue(ctx, userKey, user), nil
}

func (s *UserService) BanUser(ctx context.Context, id int) (bool, error) {
	if err := s.Repo.UpdateUser(ctx, id, &UserUpdate{Role: (*Role)(algo.NullString(string(RoleBanned)))}); err != nil {
		return false, err
	}
	return true, nil
}

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

type User struct {
	Email string   `json:"email"`
	Name  *string  `json:"name"`
	Role  Role     `json:"role"`
	Tags  []string `json:"tags"`

	Repo UserRepo `json:"-"`
	ID   int      `json:"-"`
}

func (u *User) RequirePermission(action Action) error {
	needRole := ActionRole[action]
	if u.Role.Value() < needRole.Value() {
		return errors.New("permission denied")
	}
	return nil
}

func (u *User) SetName(ctx context.Context, name string) error {
	if u.Name != nil {
		return uerr.New(uerr.ParamsError, "already have a name")
	}
	if err := u.Repo.UpdateUser(ctx, u.ID, &UserUpdate{Name: &name}); err != nil {
		return err
	}
	u.Name = &name
	return nil
}

func purifyTags(includes []string, exclude string) []string {
	inMap := map[string]struct{}{}
	var tags []string
	for _, t := range includes {
		if _, ok := inMap[t]; !ok && t != exclude {
			tags = append(tags, t)
			inMap[t] = struct{}{}
		}
	}
	return tags
}

func (u *User) SyncTags(ctx context.Context, user *User, tags []string) error {
	tagSet := purifyTags(tags, "")
	if err := u.Repo.UpdateUser(ctx, u.ID, &UserUpdate{Tags: tagSet}); err != nil {
		return err
	}
	u.Tags = tagSet
	return nil
}

func (u *User) AddSubbedTag(ctx context.Context, user *User, tag string) error {
	tags := append(u.Tags, tag)
	return u.SyncTags(ctx, u, tags)
}

func (u *User) DelSubbedTag(ctx context.Context, user *User, tag string) error {
	tagSet := purifyTags(u.Tags, tag)
	if err := u.Repo.UpdateUser(ctx, u.ID, &UserUpdate{Tags: tagSet}); err != nil {
		return err
	}
	u.Tags = tagSet
	return nil
}
