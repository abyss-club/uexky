package redis

import (
	red "github.com/go-redis/redis/v7"
	"gitlab.com/abyss.club/uexky/lib/config"
)

func NewClient() (*red.Client, error) {
	opt, err := red.ParseURL(config.Get().RedisURI)
	if err != nil {
		return nil, err
	}
	return red.NewClient(opt), nil
}
