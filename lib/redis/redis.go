package redis

import (
	red "github.com/go-redis/redis/v7"
	"gitlab.com/abyss.club/uexky/lib/config"
	"gitlab.com/abyss.club/uexky/lib/errors"
)

func NewClient() (*red.Client, error) {
	opt, err := red.ParseURL(config.Get().RedisURI)
	if err != nil {
		return nil, errors.BadParams.Handlef(err, "parse redis uri")
	}
	return red.NewClient(opt), nil
}

func ErrHandle(err error, a ...interface{}) error {
	if errors.Is(err, red.Nil) {
		return errors.NotFound.Handle(err, a...)
	}
	return errors.Redis.Handle(err, a...)
}

func ErrHandlef(err error, format string, a ...interface{}) error {
	if errors.Is(err, red.Nil) {
		return errors.NotFound.Handlef(err, format, a...)
	}
	return errors.Redis.Handlef(err, format, a...)
}
