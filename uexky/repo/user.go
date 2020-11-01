package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/go-redis/redis/v7"
	"gitlab.com/abyss.club/uexky/lib/errors"
	"gitlab.com/abyss.club/uexky/lib/postgres"
	librd "gitlab.com/abyss.club/uexky/lib/redis"
	"gitlab.com/abyss.club/uexky/lib/uid"
	"gitlab.com/abyss.club/uexky/uexky/entity"
)

type UserRepo struct {
	Redis *redis.Client
}

func userRedisKey(id uid.UID) string {
	return fmt.Sprintf("uid:%s", id.ToBase64String())
}

func (u *UserRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user User
	if err := db(ctx).Model(&user).Where("email = ?", email).Select(); err != nil {
		return nil, postgres.ErrHandlef(err, "GetUserByEmail(email=%v)", email)
	}
	return user.ToEntity(), nil
}

func (u *UserRepo) GetGuestByID(ctx context.Context, id uid.UID) (*entity.User, error) {
	var user User
	data, err := u.Redis.Get(userRedisKey(id)).Result()
	if err != nil {
		return nil, librd.ErrHandlef(err, "GetUser(id=%v)", id)
	}
	if err := json.Unmarshal([]byte(data), &user); err != nil {
		return nil, errors.Internal.Handlef(err, "unmarshal user: %s", data)
	}
	return user.ToEntity(), nil
}

func (u *UserRepo) GetByID(ctx context.Context, id uid.UID) (*entity.User, error) {
	user, err := u.GetGuestByID(ctx, id)
	if err == nil {
		return user, nil
	}
	if !errors.Is(err, errors.NotFound) {
		return nil, err
	}

	// err is redis.Nil, find in database
	var rUser User
	if err := db(ctx).Model(&rUser).Where("id = ?", id).Select(); err != nil {
		return nil, postgres.ErrHandlef(err, "GetUserByID(id=%v)", id)
	}
	return rUser.ToEntity(), nil
}

func (u *UserRepo) Insert(ctx context.Context, user *entity.User) (*entity.User, error) {
	rUser := NewUserFromEntity(user)
	if user.Role == entity.RoleGuest {
		rUser.CreatedAt = time.Now()
		rUser.UpdatedAt = time.Now()
		data, err := json.Marshal(&rUser)
		if err != nil {
			return nil, errors.BadParams.Handlef(err, "InsertUser(user=%+v)", user)
		}
		if _, err := u.Redis.Set(userRedisKey(user.ID), data, entity.GuestExpireTime).Result(); err != nil {
			return nil, librd.ErrHandlef(err, "InsertUser(user=%+v)", user)
		}
	} else {
		_, err := db(ctx).Model(rUser).Returning("*").Insert()
		if err != nil {
			return nil, postgres.ErrHandlef(err, "InsertUser(user=%+v)", user)
		}
	}
	return rUser.ToEntity(), nil
}

func (u *UserRepo) Update(ctx context.Context, user *entity.User) (*entity.User, error) {
	rUser := NewUserFromEntity(user)
	if user.Role == entity.RoleGuest {
		data, err := json.Marshal(&rUser)
		if err != nil {
			return nil, errors.BadParams.Handlef(err, "UpdateUser(user=%+v)", user)
		}
		if _, err := u.Redis.Set(userRedisKey(user.ID), data, entity.GuestExpireTime).Result(); err != nil {
			return nil, librd.ErrHandlef(err, "UpdateUser(user=%+v)", user)
		}
		return user, nil
	}
	q := db(ctx).Model(rUser).Where("id = ?", rUser.ID).
		Set("name = ?", rUser.Name).
		Set("role = ?", rUser.Role).
		Set("tags = ?", pg.Array(rUser.Tags)).
		Set("last_read_noti = ?", rUser.LastReadNoti).
		Returning("*")
	_, err := q.Update()
	if err != nil {
		return nil, errors.Wrapf(err, "UpdateUser(user=%+v)", user)
	}
	return rUser.ToEntity(), nil
}

func (u *UserRepo) ThreadSlice(ctx context.Context, user *entity.User, sq entity.SliceQuery) (*entity.ThreadSlice, error) {
	qf := func(prev *orm.Query) *orm.Query {
		return prev.Where("user_id = ?", user.ID)
	}
	return getThreadSlice(ctx, qf, &sq)
}

func (u *UserRepo) PostSlice(ctx context.Context, user *entity.User, sq entity.SliceQuery) (*entity.PostSlice, error) {
	qf := func(prev *orm.Query) *orm.Query {
		return prev.Where("user_id = ?", user.ID)
	}
	return getPostSlice(ctx, qf, &sq, true)
}
