package auth

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v7"
	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky/lib/uerr"
)

type Repo struct {
	Redis *redis.Client
}

func (r *Repo) SetCode(ctx context.Context, email string, code Code) error {
	_, err := r.Redis.Set(string(code), email, CodeExpire).Result()
	return redisErrWrapf(err, "SetCode(email=%s, code=%s)", email, code)
}

func (r *Repo) GetCodeEmail(ctx context.Context, code Code) (string, error) {
	email, err := r.Redis.Get(string(code)).Result()
	return email, redisErrWrapf(err, "GetCodeEmail(code=%s)", code)
}

func (r *Repo) DelCode(ctx context.Context, code Code) error {
	_, err := r.Redis.Del(string(code)).Result()
	return redisErrWrapf(err, "DelCode(code=%s)", code)
}

func (r *Repo) GetToken(ctx context.Context, tok string) (*Token, error) {
	data, err := r.Redis.Get(tok).Result()
	if err != nil {
		return nil, redisErrWrapf(err, "GetToken(tok=%s)", tok)
	}
	var token Token
	if err := json.Unmarshal([]byte(data), &token); err != nil {
		return nil, uerr.Wrapf(uerr.InternalError, err, "GetToken(tok=%s) unmarshal json: %s", tok, data)
	}
	return &token, nil
}

func (r *Repo) SetToken(ctx context.Context, token *Token) error {
	data, err := json.Marshal(token)
	if err != nil {
		return uerr.Errorf(uerr.PermissionError, "SetToken(token=%+v), marshal json", token)
	}
	_, err = r.Redis.Set(token.Tok, data, TokenExpire).Result()
	return redisErrWrapf(err, "SetToken(token=%+v)", token)
}

func redisErrWrapf(err error, format string, a ...interface{}) error {
	errType := uerr.RedisError
	if errors.Is(err, redis.Nil) {
		errType = uerr.NotFoundError
	}
	return uerr.Wrapf(errType, err, format, a...)
}
