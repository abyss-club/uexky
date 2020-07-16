package postgres

import (
	"context"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gitlab.com/abyss.club/uexky/lib/config"
	"gitlab.com/abyss.club/uexky/lib/uerr"
)

func NewDB() (*pg.DB, error) {
	opt, err := pg.ParseURL(config.Get().PostgresURI)
	opt.PoolSize = 16
	if err != nil {
		return nil, uerr.Wrap(uerr.ParamsError, err, "parse postgres uri")
	}
	return pg.Connect(opt), nil
}

type Session interface {
	Exec(query interface{}, params ...interface{}) (pg.Result, error)
	Model(model ...interface{}) *orm.Query
	Query(model interface{}, query interface{}, params ...interface{}) (pg.Result, error)

	Insert(model ...interface{}) error
}

func GetSessionFromContext(ctx context.Context) Session {
	data := getData(ctx)
	return data.GetSession()
}

// TxAdapter implements adapter.Tx
type TxAdapter struct {
	DB *pg.DB
}

type contextKey int

const pgDataKey contextKey = 1

type PgContextData struct {
	DB    *pg.DB
	Tx    *pg.Tx
	Layer int
}

func getData(ctx context.Context) *PgContextData {
	data, ok := ctx.Value(pgDataKey).(*PgContextData)
	if !ok || data.DB == nil {
		panic("can't not find db in context")
	}
	return data
}

func (d *PgContextData) GetSession() Session {
	if d.Tx != nil && d.Layer > 0 {
		return d.Tx
	}
	return d.DB
}

func (tx *TxAdapter) AttachDB(ctx context.Context) context.Context {
	return context.WithValue(ctx, pgDataKey, &PgContextData{DB: tx.DB})
}

func (tx *TxAdapter) WithTx(ctx context.Context, fn func() error) error {
	data := getData(ctx)
	if data.Tx != nil {
		// already inside transaction, do nothing
		return fn()
	}

	transaction, err := data.DB.Begin()
	if err != nil {
		return uerr.Wrap(uerr.PostgresError, err, "begin transaction")
	}
	data.Tx = transaction

	if err := fn(); err != nil {
		rbErr := data.Tx.Rollback()
		if rbErr != nil {
			log.Error(rbErr)
			return errors.Wrapf(err, "rollback failed: %v", rbErr)
		}
		return err
	}

	if err := data.Tx.Commit(); err != nil {
		return uerr.Wrap(uerr.PostgresError, err, "commit transaction")
	}
	return nil
}
