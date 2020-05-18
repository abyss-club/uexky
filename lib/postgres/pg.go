package postgres

import (
	"context"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
)

func NewDB(uri string) (*pg.DB, error) {
	opt, err := pg.ParseURL(uri)
	if err != nil {
		return nil, err
	}
	return pg.Connect(opt), nil
}

type Session interface {
	Exec(query interface{}, params ...interface{}) (pg.Result, error)
	Model(model ...interface{}) *orm.Query
	Query(model interface{}, query interface{}, params ...interface{}) (pg.Result, error)

	Insert(model ...interface{}) error
}

type contextKey int

const (
	sessionKey contextKey = 1
)

func ContextWithSession(ctx context.Context, session Session) context.Context {
	return context.WithValue(ctx, sessionKey, session)
}

func GetSessionFromContext(ctx context.Context) Session {
	session, ok := ctx.Value(sessionKey).(Session)
	if !ok || session == nil {
		panic("can't find postgres session")
	}
	return session
}
