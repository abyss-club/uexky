// user aggragate: user

package entity

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
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
		return nil, uerr.New(uerr.AuthError, "permission denied, no user found")
	}
	if err := user.RequirePermission(action); err != nil {
		return nil, errors.Wrapf(err, "RequirePermission(action=%v)", action)
	}
	return user, nil
}

// -- sign in/up by email

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
	code := Code(uid.RandomBase64Str(CodeLength))
	if err := s.Repo.SetCode(ctx, email, string(code)); err != nil {
		return "", errors.Wrapf(err, "TrySignInByEmail(email=%s)", email)
	}
	mail := newAuthMail(email, code)
	if err := s.Mail.SendEmail(ctx, mail); err != nil {
		return "", errors.Wrapf(err, "TrySignInByEmail(email=%s)", email)
	}
	return code, nil
}

func (s *UserService) SignInByCode(ctx context.Context, code string) (*User, string, error) {
	email, err := s.Repo.GetCodeEmail(ctx, code)
	if err != nil {
		return nil, "", errors.Wrapf(err, "SignInByCode(code=%s)", code)
	}
	user, err := s.Repo.GetUserByAuthInfo(ctx, AuthInfo{Email: email, IsGuest: false})
	if err != nil {
		if errors.Is(err, uerr.New(uerr.NotFoundError)) {
			return nil, email, nil
		}
		return nil, "", errors.Wrapf(err, "SignInByCode(code=%s)", code)
	}
	return user, email, nil
}

func (s *UserService) SignInByToken(ctx context.Context, tok string) (*User, *Token, error) {
	if tok == "" {
		return nil, nil, nil
	}
	token, err := s.Repo.GetToken(ctx, tok)
	if err != nil {
		if errors.Is(err, uerr.New(uerr.NotFoundError)) {
			return nil, nil, nil
		}
		return nil, nil, errors.Wrapf(err, "SignInByToken(tok=%s)", tok)
	}
	user, err := s.Repo.GetUserByAuthInfo(ctx, AuthInfo{UserID: token.UserID, IsGuest: token.UserRole == RoleGuest})
	return user, token, errors.Wrapf(err, "SignInByToken(tok=%s)", tok)
}

func (s *UserService) NewUser(ctx context.Context, user *User) (*User, error) {
	user, err := s.Repo.InsertUser(ctx, user)
	return user, errors.Wrapf(err, "InsertUser(user=%+v)", user)
}

func (s *UserService) BanUser(ctx context.Context, id uid.UID) (*User, error) {
	user, err := s.Repo.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, uerr.New(uerr.NotFoundError)) {
			return nil, nil // guest user is expired
		}
		return nil, errors.Wrapf(err, "BanUser(id=%v)", id)
	}
	user.Role = RoleBanned
	user, err = s.Repo.UpdateUser(ctx, user)
	if err != nil {
		return nil, errors.Wrapf(err, "BanUser(id=%v)", id)
	}
	return user, nil
}

type User struct {
	Email *string  `json:"email"`
	Name  *string  `json:"name"`
	Role  Role     `json:"role"`
	Tags  []string `json:"tags"`

	Repo         UserRepo `json:"-"`
	ID           uid.UID  `json:"-"`
	LastReadNoti uid.UID  `json:"-"`
}

func (u *User) RequirePermission(action Action) error {
	needRole := ActionRole[action]
	if u.Role.Value() < needRole.Value() {
		return uerr.New(uerr.PermissionError, "permission denied")
	}
	return nil
}

func (u *User) SetName(ctx context.Context, name string) (*User, error) {
	if u.Name != nil {
		return nil, uerr.New(uerr.ParamsError, "already have a name")
	}
	u.Name = algo.NullString(name)
	user, err := u.Repo.UpdateUser(ctx, u)
	if err != nil {
		return nil, errors.Wrapf(err, "SetName(name=%s)", name)
	}
	return user, nil
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

func (u *User) SyncTags(ctx context.Context, tags []string) (*User, error) {
	u.Tags = purifyTags(tags, "")
	user, err := u.Repo.UpdateUser(ctx, u)
	if err != nil {
		return nil, errors.Wrapf(err, "SyncTags(user=%+v, tags=%v)", u, tags)
	}
	return user, nil
}

func (u *User) AddSubbedTag(ctx context.Context, tag string) (*User, error) {
	tags := append(u.Tags, tag)
	user, err := u.SyncTags(ctx, tags)
	return user, errors.Wrapf(err, "AddSubbedTag(user=%+v, tag=%v)", u, tag)
}

func (u *User) DelSubbedTag(ctx context.Context, tag string) (*User, error) {
	tagSet := purifyTags(u.Tags, tag)
	user, err := u.SyncTags(ctx, tagSet)
	return user, errors.Wrapf(err, "DelSubbedTag(user=%+v, tag=%v)", u, tag)
}

func (u *User) NotiReceivers() []Receiver {
	return []Receiver{SendToUser(u.ID), SendToGroup(AllUser)}
}

func (u *User) SetToken(ctx context.Context, prev *Token) (*Token, error) {
	var token *Token
	if prev != nil {
		token = prev
	} else {
		token = &Token{
			Tok:      uid.RandomBase64Str(TokenLength),
			Expire:   TokenExpire,
			UserID:   u.ID,
			UserRole: u.Role,
		}
	}
	err := u.Repo.SetToken(ctx, token)
	return token, errors.Wrap(err, "SetToken")
}

func (u *User) AttachContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, userKey, u)
}
