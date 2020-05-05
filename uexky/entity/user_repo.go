package entity

import "context"

type UserUpdate struct {
	Name *string
	Role *string
}

type UserRepo interface {
	SetCode(ctx context.Context, email string, code string, ex int) error
	GetCodeEmail(ctx context.Context, code string) (string, error)
	SetToken(ctx context.Context, email string, tok string, ex int) error
	GetTokenEmail(ctx context.Context, tok string) (string, error)

	GetOrInsertUser(ctx context.Context, email string) (*User, error)
	UpdateUser(ctx context.Context, id int, update *UserUpdate) error
}

type Mail struct {
	From    string
	To      string
	Subject string
	Text    string
	HTML    string
}

type MailService interface {
	SendEmail(ctx context.Context, mail *Mail) error
}
