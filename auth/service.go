package auth

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-redis/redis/v7"
	log "github.com/sirupsen/logrus"
	"gitlab.com/abyss.club/uexky/adapter"
	"gitlab.com/abyss.club/uexky/lib/config"
	"gitlab.com/abyss.club/uexky/lib/errors"
	"gitlab.com/abyss.club/uexky/lib/uid"
)

type Service struct {
	Repo *Repo
	Mail adapter.MailAdapter
}

// ---- Sign in/sign up by Email ----

func (s *Service) TrySignInByEmail(ctx context.Context, email string, redirectTo string) (Code, error) {
	// TODO: validate email
	if redirectTo != "" && !strings.HasPrefix(redirectTo, "/") {
		return "", errors.BadParams.New("invalid redirect target")
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
func (s *Service) SignInByCode(ctx context.Context, code Code) (*Token, error) {
	email, err := s.Repo.GetCodeEmail(ctx, code)
	if err != nil {
		return nil, errors.Wrapf(err, "SignInByCode(code=%s)", code)
	}
	token := NewEmailToken(email)
	if err := s.Repo.SetToken(ctx, token); err != nil {
		return nil, errors.Wrap(err, "SetToken")
	}
	if err := s.Repo.DelCode(ctx, code); err != nil {
		log.Error(errors.Wrap(err, "DelCode"))
	}
	return token, nil
}

// ---- Guest user ----

func (s *Service) SignInGuestUser(ctx context.Context) (*Token, error) {
	token := NewGuestToken()
	if err := s.Repo.SetToken(ctx, token); err != nil {
		return nil, errors.Wrap(err, "SetToken")
	}
	return token, nil
}

// ---- Regular apis ----

func (s *Service) GetToken(ctx context.Context, tok string) (*Token, error) {
	token, err := s.Repo.GetToken(ctx, tok)
	if err != nil {
		return nil, errors.Wrap(err, "GetToken")
	}
	// refresh ttl
	if err := s.Repo.SetToken(ctx, token); err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "GetToken")
	}
	return token, nil
}
