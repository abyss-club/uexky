// user aggragate: user

package entity

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
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
		return nil, uerr.New(uerr.AuthError, "permission denied, no user found")
	}
	if err := user.RequirePermission(action); err != nil {
		return nil, errors.Wrapf(err, "RequirePermission(action=%v)", action)
	}
	return user, nil
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
		return "", errors.Wrapf(err, "TrySignInByEmail(email=%s)", email)
	}
	mail := newAuthMail(email, code)
	if err := s.Mail.SendEmail(ctx, mail); err != nil {
		return "", errors.Wrapf(err, "TrySignInByEmail(email=%s)", email)
	}
	return code, nil
}

func (s *UserService) SignInByCode(ctx context.Context, code string) (Token, error) {
	email, err := s.Repo.GetCodeEmail(ctx, code)
	if err != nil {
		return Token{}, errors.Wrapf(err, "SignInByCode(code=%s)", code)
	}
	tok := uid.RandomBase64Str(tokenLength)
	if err := s.Repo.SetToken(ctx, email, tok, tokenExpire); err != nil {
		return Token{}, errors.Wrapf(err, "SignInByCode(code=%s)", code)
	}
	if err := s.Repo.DelCode(ctx, code); err != nil {
		log.Error(err)
	}
	return Token{Tok: tok, Expire: tokenExpire}, nil
}

func (s *UserService) CtxWithUserByToken(ctx context.Context, tok string) (ct context.Context, isNew bool, err error) {
	email, err := s.Repo.GetTokenEmail(ctx, tok)
	if err != nil {
		return nil, false, errors.Wrapf(err, "CtxWithUserByToken(tok=%s)", tok)
	}
	if email == "" {
		user := &User{
			Role: RoleGuest,
		}
		return context.WithValue(ctx, userKey, user), false, nil
	}
	user, isNew, err := s.Repo.GetOrInsertUser(ctx, email)
	if err != nil {
		return nil, false, errors.Wrapf(err, "CtxWithUserByToken(tok=%s)", tok)
	}
	return context.WithValue(ctx, userKey, user), isNew, nil
}

func (s *UserService) BanUser(ctx context.Context, id int64) (bool, error) {
	if err := s.Repo.UpdateUser(ctx, id, &UserUpdate{Role: (*Role)(algo.NullString(string(RoleBanned)))}); err != nil {
		return false, errors.Wrapf(err, "BanUser(id=%v)", id)
	}
	return true, nil
}

type User struct {
	Email string   `json:"email"`
	Name  *string  `json:"name"`
	Role  Role     `json:"role"`
	Tags  []string `json:"tags"`

	Repo         UserRepo `json:"-"`
	ID           int64    `json:"-"`
	LastReadNoti uid.UID  `json:"-"`
}

func (u *User) RequirePermission(action Action) error {
	needRole := ActionRole[action]
	if u.Role.Value() < needRole.Value() {
		return uerr.New(uerr.PermissionError, "permission denied")
	}
	return nil
}

func (u *User) SetName(ctx context.Context, name string) error {
	if u.Name != nil {
		return uerr.New(uerr.ParamsError, "already have a name")
	}
	if err := u.Repo.UpdateUser(ctx, u.ID, &UserUpdate{Name: &name}); err != nil {
		return errors.Wrapf(err, "SetName(name=%s)", name)
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

func (u *User) SyncTags(ctx context.Context, tags []string) error {
	tagSet := purifyTags(tags, "")
	if err := u.Repo.UpdateUser(ctx, u.ID, &UserUpdate{Tags: tagSet}); err != nil {
		return errors.Wrapf(err, "SyncTags(user=%+v, tags=%v)", u, tags)
	}
	u.Tags = tagSet
	return nil
}

func (u *User) AddSubbedTag(ctx context.Context, tag string) error {
	tags := append(u.Tags, tag)
	err := u.SyncTags(ctx, tags)
	return errors.Wrapf(err, "AddSubbedTag(user=%+v, tag=%v)", u, tag)
}

func (u *User) DelSubbedTag(ctx context.Context, tag string) error {
	tagSet := purifyTags(u.Tags, tag)
	err := u.SyncTags(ctx, tagSet)
	return errors.Wrapf(err, "DelSubbedTag(user=%+v, tag=%v)", u, tag)
}

func (u *User) NotiReceivers() []Receiver {
	return []Receiver{SendToUser(u.ID), SendToGroup(AllUser)}
}
