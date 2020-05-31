package adapter

import "context"

type Tx interface {
	AttachDB(ctx context.Context) context.Context
	Begin(ctx context.Context) error
	Commit(ctx context.Context) error
	Rollback(ctx context.Context, err error) error
}
