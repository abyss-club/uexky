package api

import (
	"context"

	"github.com/nanozuki/uexky/model"
)

// AccountResolver for graphql
type AccountResolver struct {
	Account *model.Account
}

// Token resolve account.token
func (ar *AccountResolver) Token(ctx context.Context) (string, error) {
	return ar.Account.Token, nil
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
