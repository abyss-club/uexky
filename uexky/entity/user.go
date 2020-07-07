// user aggragate: user

package entity

import (
	"context"
	"errors"
	"fmt"

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

type User struct {
	Email string   `json:"email"`
	Name  *string  `json:"name"`
	Role  Role     `json:"role"`
	Tags  []string `json:"tags"`

	Repo         UserRepo `json:"-"`
	ID           int      `json:"-"`
	LastReadNoti uid.UID  `json:"-"`
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

func (u *User) NotiReceivers() []Receiver {
	return []Receiver{SendToUser(u.ID), SendToGroup(AllUser)}
}
