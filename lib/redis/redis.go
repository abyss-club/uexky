package redis

import (
	red "github.com/go-redis/redis/v7"
)

func NewClient(uri string) (*red.Client, error) {
	opt, err := red.ParseURL(uri)
	if err != nil {
		return nil, err
	}
	return red.NewClient(opt), nil
}
