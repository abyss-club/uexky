package redis

import (
	red "github.com/go-redis/redis/v7"
	"gitlab.com/abyss.club/uexky/lib/config"
	"gitlab.com/abyss.club/uexky/lib/uerr"
)

func NewClient() (*red.Client, error) {
	opt, err := red.ParseURL(config.Get().RedisURI)
	if err != nil {
		return nil, uerr.Wrapf(uerr.ParamsError, err, "parse redis uri")
	}
	return red.NewClient(opt), nil
}
