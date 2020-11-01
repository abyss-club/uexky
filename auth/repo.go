package auth

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v7"
	"gitlab.com/abyss.club/uexky/lib/errors"
	librd "gitlab.com/abyss.club/uexky/lib/redis"
)

type Repo struct {
	Redis *redis.Client
}

func (r *Repo) SetCode(ctx context.Context, email string, code Code) error {
	_, err := r.Redis.Set(string(code), email, CodeExpire).Result()
	return librd.ErrHandlef(err, "SetCode(email=%s, code=%s)", email, code)
}

func (r *Repo) GetCodeEmail(ctx context.Context, code Code) (string, error) {
	email, err := r.Redis.Get(string(code)).Result()
	return email, librd.ErrHandlef(err, "GetCodeEmail(code=%s)", code)
}

func (r *Repo) DelCode(ctx context.Context, code Code) error {
	_, err := r.Redis.Del(string(code)).Result()
	return librd.ErrHandlef(err, "DelCode(code=%s)", code)
}

func (r *Repo) GetToken(ctx context.Context, tok string) (*Token, error) {
	data, err := r.Redis.Get(tok).Result()
	if err != nil {
		return nil, librd.ErrHandlef(err, "GetToken(tok=%s)", tok)
	}
	var token Token
	if err := json.Unmarshal([]byte(data), &token); err != nil {
		return nil, errors.Internal.Handlef(err, "GetToken(tok=%s) unmarshal json: %s", tok, data)
	}
	return &token, nil
}

func (r *Repo) SetToken(ctx context.Context, token *Token) error {
	data, err := json.Marshal(token)
	if err != nil {
		return errors.Permission.Errorf("SetToken(token=%+v), marshal json", token)
	}
	_, err = r.Redis.Set(token.Tok, data, TokenExpire).Result()
	return librd.ErrHandlef(err, "SetToken(token=%+v)", token)
}
