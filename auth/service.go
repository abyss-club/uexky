package auth

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gitlab.com/abyss.club/uexky/lib/config"
	"gitlab.com/abyss.club/uexky/lib/uerr"
	"gitlab.com/abyss.club/uexky/lib/uid"
	"gitlab.com/abyss.club/uexky/uexky"
	"gitlab.com/abyss.club/uexky/uexky/adapter"
	"gitlab.com/abyss.club/uexky/uexky/entity"
)

type Service struct {
	R     R
	Repo  Repo
	Mail  adapter.MailAdapter
	Uexky *uexky.Service
}

type R struct {
	User entity.UserRepo
}

type Repo interface {
	SetCode(ctx context.Context, email string, code Code) error
	GetCodeEmail(ctx context.Context, code Code) (string, error)
	DelCode(ctx context.Context, code Code) error

	GetUserByAuthInfo(ctx context.Context, ai AuthInfo) (*entity.User, error)
	GetToken(ctx context.Context, tok string) (*Token, error)
	SetToken(ctx context.Context, token *Token) error
}

func (s *Service) TrySignInByEmail(ctx context.Context, email string, redirectTo string) (Code, error) {
	// TODO: validate email
	if redirectTo != "" && !strings.HasPrefix(redirectTo, "/") {
		return "", uerr.New(uerr.ParamsError, "invalid redirect target")
	}
	code := Code(uid.RandomBase64Str(CodeLength))
	if err := s.Repo.SetCode(ctx, email, code); err != nil {
		return "", errors.Wrapf(err, "TrySignInByEmail(email=%s)", email)
	}
	mail := newAuthMail(email, code, redirectTo)
	if err := s.Mail.SendEmail(ctx, mail); err != nil {
		return "", errors.Wrapf(err, "TrySignInByEmail(email=%s)", email)
	}
	return code, nil
}

func newAuthMail(email string, code Code, redirectTo string) *adapter.Mail {
	srvCfg := &(config.Get().Server)
	authURL := code.SignInURL(redirectTo)
	return &adapter.Mail{
		From:    fmt.Sprintf("auth@%s", srvCfg.Domain),
		To:      email,
		Subject: "点击登入 Abyss!",
		Text:    fmt.Sprintf("点击此链接进入 Abyss：%s", authURL),
		HTML:    fmt.Sprintf(authEmailHTML, authURL),
	}
}

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

// SignInByCode is only for signed in user
func (s *Service) SignInByCode(ctx context.Context, code string) (*Token, error) {
	user, email, err := s.signInByCode(ctx, Code(code))
	if err != nil {
		return nil, err
	}
	if user == nil {
		user := entity.NewSignedInUser(email)
		user, err = s.R.User.Insert(ctx, user)
		if err != nil {
			return nil, err
		}
		if err := s.Uexky.NewNotiOnNewUser(ctx, user); err != nil {
			log.Errorf("%+v", err)
		}
	}
	return s.SetToken(ctx, user, nil)
}

func (s *Service) signInByCode(ctx context.Context, code Code) (*entity.User, string, error) {
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

// CtxWithUserByToken add user to context by tok is for both signed user and guest user.
func (s *Service) CtxWithUserByToken(ctx context.Context, tok string) (context.Context, *Token, error) {
	user, token, err := s.signInByToken(ctx, tok)
	if err != nil {
		return nil, nil, err
	}
	if user == nil {
		// must be unsigned user
		user := entity.NewGuestUser()
		_, err = s.R.User.Insert(ctx, user)
		if err != nil {
			return nil, nil, err
		}
	}
	// no need check token is nil
	// if cannot find user or token, token here is nil, and will make a new one.
	token, err = s.SetToken(ctx, user, token)
	if err != nil {
		return nil, nil, err
	}
	return user.AttachContext(ctx), token, err
}

func (s *Service) signInByToken(ctx context.Context, tok string) (*entity.User, *Token, error) {
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
	user, err := s.Repo.GetUserByAuthInfo(ctx, AuthInfo{UserID: token.UserID, IsGuest: token.UserRole == entity.RoleGuest})
	return user, token, errors.Wrapf(err, "SignInByToken(tok=%s)", tok)
}

func (s *Service) SetToken(ctx context.Context, u *entity.User, prev *Token) (*Token, error) {
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
	err := s.Repo.SetToken(ctx, token)
	return token, errors.Wrap(err, "SetToken")
}
