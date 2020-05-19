package mail

import (
	"context"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky/config"
	"gitlab.com/abyss.club/uexky/uexky/entity"
)

type Adapter struct {
	mg mailgun.Mailgun
}

func NewAdapter() *Adapter {
	mg := mailgun.NewMailgun(config.Get().Server.Domain, config.Get().Mail.PrivateKey)
	return &Adapter{mg: mg}
}

func (a *Adapter) SendEmail(ctx context.Context, mail *entity.Mail) error {
	msg := a.mg.NewMessage(mail.From, mail.Subject, mail.Text, mail.To)
	msg.SetHtml(mail.HTML)
	resp, id, err := a.mg.Send(ctx, msg)
	return errors.Wrapf(err, "mailgun send mail error(%s): %s", id, resp)
}
