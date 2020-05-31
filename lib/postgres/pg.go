package postgres

import (
	"context"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	log "github.com/sirupsen/logrus"
	"gitlab.com/abyss.club/uexky/lib/config"
	"gitlab.com/abyss.club/uexky/lib/uerr"
)

func NewDB() (*pg.DB, error) {
	opt, err := pg.ParseURL(config.Get().PostgresURI)
	opt.PoolSize = 20
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

func (tx *TxAdapter) Begin(ctx context.Context) error {
	data := getData(ctx)
	data.Layer++
	if data.Layer > 1 {
		return nil
	}
	var err error
	data.Tx, err = data.DB.Begin()
	if err != nil {
		return uerr.Errorf(uerr.DBError, "begin transaction: %w", err)
	}
	return nil
}

func (tx *TxAdapter) Commit(ctx context.Context) error {
	data := getData(ctx)
	if data.Tx == nil || data.Layer < 1 {
		return nil
	}
	data.Layer--
	if data.Layer > 0 {
		return nil
	}
	if err := data.Tx.Commit(); err != nil {
		return uerr.Errorf(uerr.DBError, "commit transaction: %w", err)
	}
	return nil
}

func (tx *TxAdapter) Rollback(ctx context.Context, err error) error {
	data := getData(ctx)
	if err == nil {
		return nil
	}
	if data.Layer < 1 {
		return err
	}
	data.Layer = 0
	rbErr := data.Tx.Rollback()
	if rbErr != nil {
		log.Error(err)
		return uerr.Errorf(uerr.DBError, "rollback transaction: %w", err)
	}
	return err
}
