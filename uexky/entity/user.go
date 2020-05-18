// user aggragate: user

package entity

import (
	"context"
	"errors"
	"fmt"

	"gitlab.com/abyss.club/uexky/config"
	"gitlab.com/abyss.club/uexky/lib/uid"
)

type UserService struct {
	Repo UserRepo
	Mail MailService
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
	codeExpire  = 1200 // 20 minutes
	tokenLength = 24
	tokenExpire = 86400 * 30 // 30 days
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

func newAuthMail(email, code string) *Mail {
	srvCfg := &(config.Get().Server)
	authURL := fmt.Sprintf("%s://%s/auth/?code=%s", srvCfg.Proto, srvCfg.APIDomain, code)
	return &Mail{
		From:    fmt.Sprintf("auth@%s", srvCfg.Domain),
		To:      email,
		Subject: "点击登入 Abyss!",
		Text:    fmt.Sprintf("点击此链接进入 Abyss：%s", authURL),
		HTML:    fmt.Sprintf(authEmailHTML, authURL),
	}
}

func (s *UserService) TrySignInByEmail(ctx context.Context, email string) (bool, error) {
	// TODO: validate email
	code := uid.RandomBase64Str(codeLength)
	if err := s.Repo.SetCode(ctx, email, code, codeExpire); err != nil {
		return false, err
	}
	mail := newAuthMail(email, code)
	if err := s.Mail.SendEmail(ctx, mail); err != nil {
		return false, err
	}
	return true, nil
}

type Token struct {
	Tok          string
	ExpireSecond int
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
	return Token{Tok: tok, ExpireSecond: tokenExpire}, nil
}

func (s *UserService) CtxWithUserByToken(ctx context.Context, tok string) (context.Context, error) {
	email, err := s.Repo.GetTokenEmail(ctx, tok)
	if err != nil {
		return nil, err
	}
	if email == "" {
		user := &User{
			// Role: &RoleGuest,
		}
		return context.WithValue(ctx, userKey, user), nil
	}
	user, err := s.Repo.GetOrInsertUser(ctx, email)
	if err != nil {
		return nil, err
	}
	return context.WithValue(ctx, userKey, user), nil
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
	ActionBlockPost   = Action("BLOCK_POST")
	ActionLockThread  = Action("LOCK_THREAD")
	ActionBlockThread = Action("BLOCK_THREAD")
	ActionEditTag     = Action("EDIT_TAG")
	ActionEditSetting = Action("EDIT_SETTING")
	ActionPubPost     = Action("PUB_POST")
	ActionPubThread   = Action("PUB_THREAD")
)

var ActionRole = map[Action]Role{
	ActionProfile:     RoleNormal,
	ActionBanUser:     RoleMod,
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

	Repo         UserRepo     `json:"-"`
	ID           int          `json:"-"`
	LastReadNoti LastReadNoti `json:"-"`
}

func (u *User) RequirePermission(action Action) error {
	needRole := ActionRole[action]
	if u.Role.Value() < needRole.Value() {
		return errors.New("permission denied")
	}
	return nil
}

func (u *User) SetName(ctx context.Context, name string) error {
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

func (u *User) BanUser(ctx context.Context, id int) (bool, error) {
	banned := RoleBanned
	if err := u.Repo.UpdateUser(ctx, u.ID, &UserUpdate{Role: &banned}); err != nil {
		return false, err
	}
	u.Role = RoleBanned
	return true, nil
}
