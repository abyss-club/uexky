package adapter

import "context"

type Mail struct {
	From    string
	To      string
	Subject string
	Text    string
	HTML    string
}

type MailAdapter interface {
	SendEmail(ctx context.Context, mail *Mail) error
}
