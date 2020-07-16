package adapter

import "context"

type Tx interface {
	AttachDB(ctx context.Context) context.Context
	WithTx(ctx context.Context, fn func() error) error
}
