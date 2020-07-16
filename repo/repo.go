package repo

import (
	"errors"

	"github.com/go-pg/pg/v9"
	"github.com/go-redis/redis/v7"
	"gitlab.com/abyss.club/uexky/lib/uerr"
)

func dbErrWrap(err error, a ...interface{}) error {
	errType := uerr.DBError
	if errors.Is(err, pg.ErrNoRows) {
		errType = uerr.NotFoundError
	}
	return uerr.Wrap(errType, err, a...)
}

func dbErrWrapf(err error, format string, a ...interface{}) error {
	errType := uerr.DBError
	if errors.Is(err, pg.ErrNoRows) {
		errType = uerr.NotFoundError
	}
	return uerr.Wrapf(errType, err, format, a...)
}

//func redisErrWrap(err error, a ...interface{}) error {
//	errType := uerr.DBError
//	if errors.Is(err, redis.Nil) {
//		errType = uerr.NotFoundError
//	}
//	return uerr.Wrap(errType, err, a...)
//}

func redisErrWrapf(err error, format string, a ...interface{}) error {
	errType := uerr.DBError
	if errors.Is(err, redis.Nil) {
		errType = uerr.NotFoundError
	}
	return uerr.Wrapf(errType, err, format, a...)
}
