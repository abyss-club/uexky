package mocks

import (
	"context"

	"gitlab.com/abyss.club/uexky/adapter"
)

type MailAdapter struct {
	LastMail *adapter.Mail `wire:"-"`
}

func (a *MailAdapter) SendEmail(ctx context.Context, mail *adapter.Mail) error {
	a.LastMail = mail
	return nil
}
