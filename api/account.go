package api

import (
	"context"

	"github.com/nanozuki/uexky/model"
)

// AccountResolver for graphql
type AccountResolver struct {
	Account *model.Account
}

// Account resolve query 'account'
func (r *Resolver) Account(ctx context.Context) (*AccountResolver, error) {
	account, err := model.GetAccount(ctx)
	return &AccountResolver{account}, err
}

// AddAccount resolve mutation 'addAccount'
func (r *Resolver) AddAccount(ctx context.Context) (*AccountResolver, error) {
	account, err := model.NewAccount(ctx)
	return &AccountResolver{account}, err
}

// Token resolve account.token
func (ar *AccountResolver) Token(ctx context.Context) (string, error) {
	return ar.Account.Token, nil
}

// Names resolve account.names
func (ar *AccountResolver) Names(ctx context.Context) ([]string, error) {
	return ar.Account.Names, nil
}
