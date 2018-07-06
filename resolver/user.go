package resolver

import (
	"context"

	"gitlab.com/abyss.club/uexky/model"
)

// UserResolver for graphql
type UserResolver struct {
	User *model.User
}

// Email resolve user.email
func (ur *UserResolver) Email(ctx context.Context) (string, error) {
	return ur.User.Email, nil
}

// Name resolve user.name
func (ur *UserResolver) Name(ctx context.Context) (*string, error) {
	if ur.User.Name == "" {
		return nil, nil
	}
	return &(ur.User.Name), nil
}

// Tags ...
func (ur *UserResolver) Tags(ctx context.Context) (*[]string, error) {
	if len(ur.User.Tags) == 0 {
		return nil, nil
	}
	return &(ur.User.Tags), nil
}
