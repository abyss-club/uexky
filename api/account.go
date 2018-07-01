package api

import (
	"context"

	"gitlab.com/abyss.club/uexky/model"
)

// AccountResolver for graphql
type AccountResolver struct {
	Account *model.Account
}

// Email resolve account.token
func (ar *AccountResolver) Email(ctx context.Context) (string, error) {
	return ar.Account.Email, nil
}

// Names resolve account.names
func (ar *AccountResolver) Names(ctx context.Context) (*[]string, error) {
	if len(ar.Account.Names) == 0 {
		return nil, nil
	}
	return &(ar.Account.Names), nil
}

// Tags ...
func (ar *AccountResolver) Tags(ctx context.Context) (*[]string, error) {
	if len(ar.Account.Tags) == 0 {
		return nil, nil
	}
	return &(ar.Account.Tags), nil
}
