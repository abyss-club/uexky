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
	Redis *redis.Client
	Forum *ForumRepo
}

func (u *UserRepo) SetCode(ctx context.Context, email string, code string, ex time.Duration) error {
	_, err := u.Redis.Set(code, email, ex).Result()
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
		if ai.Email == "" {
			return nil, uerr.New(uerr.ParamsError, "cannot get signed user without email")
		}
		if err := u.db(ctx).Model(&user).Where("email = ?", ai.Email).Select(); err != nil {
			return nil, dbErrWrapf(err, "GetOrInsertUser.GetUser(ai=%+v)", ai)
		}
	}
	mainTags, err := u.Forum.GetMainTags(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "GetUserByAuthInfo(ai=%+v)", ai)
	}
	return u.toEntityUser(&user, mainTags), nil

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

func (u *UserRepo) InsertUser(ctx context.Context, user *entity.User, ex time.Duration) (*entity.User, error) {
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
		if _, err := u.Redis.Set(u.userRedisKey(user.ID), data, ex).Result(); err != nil {
			return nil, redisErrWrapf(err, "InsertUser(user=%+v)", user)
		}
	} else {
		if _, err := u.db(ctx).Model(&dbUser).Returning("*").Insert(); err != nil {
			return nil, dbErrWrapf(err, "InsertUser(user=%+v)", user)
		}
	}
	mainTags, err := u.Forum.GetMainTags(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "InsertUser(user=%+v)", user)
	}
	return u.toEntityUser(&dbUser, mainTags), nil
}

func (u *UserRepo) UpdateUser(ctx context.Context, id uid.UID, update *entity.UserUpdate) error {
	user := User{}
	q := u.db(ctx).Model(&user).Where("id = ?", id)
	if update.Name != nil {
		q.Set("name = ?", update.Name)
	}
	if update.Role != nil {
		q.Set("role = ?", update.Role)
	}
	if update.Tags != nil {
		q.Set("tags = ?", pg.Array(update.Tags))
	}
	_, err := q.Update()
	return dbErrWrapf(err, "UpdateUser(id=%v, update=%+v)", id, update)
}

func (u *UserRepo) db(ctx context.Context) postgres.Session {
	return postgres.GetSessionFromContext(ctx)
}

func (u *UserRepo) userRedisKey(id uid.UID) string {
	return fmt.Sprintf("uid:%s", id.ToBase64String())
}

func (u *UserRepo) toEntityUser(user *User, mainTags []string) *entity.User {
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
		e.Tags = mainTags
	}
	if e.Role == "" {
		e.Role = entity.RoleNormal
	}
	return e
}
