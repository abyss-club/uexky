package mocks

import (
	"context"

	"gitlab.com/abyss.club/uexky/uexky/adapter"
)

type MailAdapter struct{}

func (a *MailAdapter) SendEmail(ctx context.Context, mail *adapter.Mail) error {
	return nil
}
