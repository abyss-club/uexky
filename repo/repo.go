package repo

import (
	"context"
	"errors"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/go-redis/redis/v7"
	"gitlab.com/abyss.club/uexky/lib/postgres"
	"gitlab.com/abyss.club/uexky/lib/uerr"
	"gitlab.com/abyss.club/uexky/uexky/entity"
)

func NewRepo(r *redis.Client) *entity.Repo {
	return &entity.Repo{
		User:   &UserRepo{Redis: r},
		Thread: &ThreadRepo{Redis: r},
		Post:   &PostRepo{Redis: r},
		Tag:    &TagRepo{},
		Noti:   &NotiRepo{},
	}
}

type queryFunc func(prev *orm.Query) *orm.Query

func db(ctx context.Context) postgres.Session {
	return postgres.GetSessionFromContext(ctx)
}

func dbErrWrap(err error, a ...interface{}) error {
	errType := uerr.PostgresError
	if errors.Is(err, pg.ErrNoRows) {
		errType = uerr.NotFoundError
	}
	return uerr.Wrap(errType, err, a...)
}

func dbErrWrapf(err error, format string, a ...interface{}) error {
	errType := uerr.PostgresError
	if errors.Is(err, pg.ErrNoRows) {
		errType = uerr.NotFoundError
	}
	return uerr.Wrapf(errType, err, format, a...)
}

/*
func redisErrWrap(err error, a ...interface{}) error {
	errType := uerr.RedisError
	if errors.Is(err, redis.Nil) {
		errType = uerr.NotFoundError
	}
	return uerr.Wrap(errType, err, a...)
}
*/

func redisErrWrapf(err error, format string, a ...interface{}) error {
	errType := uerr.RedisError
	if errors.Is(err, redis.Nil) {
		errType = uerr.NotFoundError
	}
	return uerr.Wrapf(errType, err, format, a...)
}
