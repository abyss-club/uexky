package mail

import (
	"context"

	"github.com/mailgun/mailgun-go/v4"
	"gitlab.com/abyss.club/uexky/adapter"
	"gitlab.com/abyss.club/uexky/lib/config"
	"gitlab.com/abyss.club/uexky/lib/errors"
)

type Adapter struct {
	mg mailgun.Mailgun
}

func NewAdapter() *Adapter {
	mg := mailgun.NewMailgun(config.Get().Mail.Domain, config.Get().Mail.PrivateKey)
	return &Adapter{mg: mg}
}

func (a *Adapter) SendEmail(ctx context.Context, mail *adapter.Mail) error {
	msg := a.mg.NewMessage(mail.From, mail.Subject, mail.Text, mail.To)
	msg.SetHtml(mail.HTML)
	resp, id, err := a.mg.Send(ctx, msg)
	return errors.Mailgun.Handlef(err, "mailgun send mail error(%s): %s", id, resp)
}
