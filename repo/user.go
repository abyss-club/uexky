package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/go-redis/redis/v7"
	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky/lib/postgres"
	"gitlab.com/abyss.club/uexky/lib/uerr"
	"gitlab.com/abyss.club/uexky/lib/uid"
	"gitlab.com/abyss.club/uexky/uexky/entity"
)

type UserRepo struct {
	Redis    *redis.Client
	MainTags *MainTag
}

func (u *UserRepo) SetCode(ctx context.Context, email string, code string) error {
	_, err := u.Redis.Set(code, email, entity.CodeExpire).Result()
	return redisErrWrapf(err, "SetCode(email=%s, code=%s)", email, code)
}

func (u *UserRepo) GetCodeEmail(ctx context.Context, code string) (string, error) {
	code, err := u.Redis.Get(code).Result()
	return code, redisErrWrapf(err, "GetCodeEmail(code=%s)", code)
}

func (u *UserRepo) DelCode(ctx context.Context, code string) error {
	_, err := u.Redis.Del(code).Result()
	return redisErrWrapf(err, "DelCode(code=%s)", code)
}

func (u *UserRepo) GetUserByID(ctx context.Context, id uid.UID) (*entity.User, error) {
	var user User
	data, err := u.Redis.Get(u.userRedisKey(id)).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, redisErrWrapf(err, "GetUserByID(id=%v)", id)
	}
	if err == nil {
		if err := json.Unmarshal([]byte(data), &user); err != nil {
			return nil, uerr.Wrapf(uerr.InternalError, err, "unmarshal user: %s", data)
		}
		return u.toEntityUser(&user), nil
	}

	// err is redis.Nil, find in database
	if err := u.db(ctx).Model(&user).Where("id = ?", id).Select(); err != nil {
		return nil, dbErrWrapf(err, "GetUserByID(id=%v)", id)
	}
	return u.toEntityUser(&user), nil

}

func (u *UserRepo) GetUserByAuthInfo(ctx context.Context, ai entity.AuthInfo) (*entity.User, error) {
	var user User
	if ai.IsGuest {
		if ai.UserID == 0 {
			return nil, uerr.New(uerr.ParamsError, "cannot get guest user without id")
		}
		data, err := u.Redis.Get(u.userRedisKey(ai.UserID)).Result()
		if err != nil {
			return nil, redisErrWrapf(err, "GetUserByAuthInfo(ai=%+v)", ai)
		}
		if err := json.Unmarshal([]byte(data), &user); err != nil {
			return nil, uerr.Wrapf(uerr.InternalError, err, "unmarshal user: %s", data)
		}
	} else {
		q := u.db(ctx).Model(&user)
		switch {
		case ai.UserID != 0:
			q = q.Where("id = ?", ai.UserID)
		case ai.Email != "":
			q = q.Where("email = ?", ai.Email)
		default:
			return nil, uerr.New(uerr.ParamsError, "cannot get signed user without id and email")
		}
		if err := q.Select(); err != nil {
			return nil, dbErrWrapf(err, "GetOrInsertUser.GetUser(ai=%+v)", ai)
		}
	}
	return u.toEntityUser(&user), nil

}

func (u *UserRepo) SetToken(ctx context.Context, token *entity.Token) error {
	data, err := json.Marshal(token)
	if err != nil {
		return uerr.Errorf(uerr.PermissionError, "SetToken(token=%+v), marshal json", token)
	}
	_, err = u.Redis.Set(token.Tok, data, token.Expire).Result()
	return redisErrWrapf(err, "SetToken(token=%+v)", token)
}

func (u *UserRepo) GetToken(ctx context.Context, tok string) (*entity.Token, error) {
	data, err := u.Redis.Get(tok).Result()
	if err != nil {
		return nil, redisErrWrapf(err, "GetToken(tok=%s)", tok)
	}
	var token entity.Token
	if err := json.Unmarshal([]byte(data), &token); err != nil {
		return nil, uerr.Wrapf(uerr.InternalError, err, "GetToken(tok=%s) unmarshal json: %s", tok, data)
	}
	return &token, nil
}

func (u *UserRepo) InsertUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	dbUser := User{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}
	if user.Role == entity.RoleGuest {
		dbUser.CreatedAt = time.Now()
		data, err := json.Marshal(&dbUser)
		if err != nil {
			return nil, uerr.Wrapf(uerr.ParamsError, err, "InsertUser(user=%+v)", user)
		}
		if _, err := u.Redis.Set(u.userRedisKey(user.ID), data, entity.TokenExpire).Result(); err != nil {
			return nil, redisErrWrapf(err, "InsertUser(user=%+v)", user)
		}
	} else {
		_, err := u.db(ctx).Model(&dbUser).Returning("*").Insert()
		if err != nil {
			return nil, dbErrWrapf(err, "InsertUser(user=%+v)", user)
		}
	}
	return u.toEntityUser(&dbUser), nil
}

func (u *UserRepo) UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	if user.Role == entity.RoleGuest {
		rUser := &User{
			ID:   user.ID,
			Role: entity.RoleGuest,
			Tags: user.Tags,
		}
		data, err := json.Marshal(&rUser)
		if err != nil {
			return nil, uerr.Wrapf(uerr.ParamsError, err, "InsertUser(user=%+v)", user)
		}
		if _, err := u.Redis.Set(u.userRedisKey(user.ID), data, entity.TokenExpire).Result(); err != nil {
			return nil, redisErrWrapf(err, "InsertUser(user=%+v)", user)
		}
		return user, nil
	}
	var rUser User
	q := u.db(ctx).Model(&rUser).Where("id = ?", user.ID).
		Set("name = ?", user.Name).
		Set("role = ?", user.Role).
		Set("tags = ?", pg.Array(user.Tags)).
		Returning("*")
	_, err := q.Update()
	if err != nil {
		return nil, errors.Wrapf(err, "InsertUser(user=%+v)", user)
	}
	return u.toEntityUser(&rUser), dbErrWrapf(err, "UpdateUser(user=%+v)", user)
}

func (u *UserRepo) db(ctx context.Context) postgres.Session {
	return postgres.GetSessionFromContext(ctx)
}

func (u *UserRepo) userRedisKey(id uid.UID) string {
	return fmt.Sprintf("uid:%s", id.ToBase64String())
}

func (u *UserRepo) toEntityUser(user *User) *entity.User {
	e := &entity.User{
		Email: user.Email,
		Name:  user.Name,
		Role:  user.Role,
		Tags:  user.Tags,

		Repo:         u,
		ID:           user.ID,
		LastReadNoti: user.LastReadNoti,
	}
	// TODO: should in service level?
	if len(user.Tags) == 0 {
		e.Tags = u.MainTags.Tags
	}
	if e.Role == "" {
		e.Role = entity.RoleNormal
	}
	return e
}
