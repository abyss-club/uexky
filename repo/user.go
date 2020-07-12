package repo

import (
	"context"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/go-redis/redis/v7"
	"gitlab.com/abyss.club/uexky/lib/postgres"
	"gitlab.com/abyss.club/uexky/lib/uid"
	"gitlab.com/abyss.club/uexky/uexky/entity"
)

type UserRepo struct {
	Redis *redis.Client
	Forum *ForumRepo
}

func (u *UserRepo) SetCode(ctx context.Context, email string, code string, ex time.Duration) error {
	_, err := u.Redis.Set(code, email, ex).Result()
	return err
}

func (u *UserRepo) GetCodeEmail(ctx context.Context, code string) (string, error) {
	return u.Redis.Get(code).Result()
}

func (u *UserRepo) DelCode(ctx context.Context, code string) error {
	_, err := u.Redis.Del(code).Result()
	return err
}

func (u *UserRepo) SetToken(ctx context.Context, email string, tok string, ex time.Duration) error {
	_, err := u.Redis.Set(tok, email, ex).Result()
	return err
}

func (u *UserRepo) GetTokenEmail(ctx context.Context, tok string) (string, error) {
	tok, err := u.Redis.Get(tok).Result()
	if err == redis.Nil {
		return "", nil
	}
	return tok, err
}

func (u *UserRepo) db(ctx context.Context) postgres.Session {
	return postgres.GetSessionFromContext(ctx)
}

func (u *UserRepo) toEntityUser(user *User, mainTags []string) *entity.User {
	entity := &entity.User{
		Email: user.Email,
		Name:  user.Name,
		Role:  entity.ParseRole(user.Role),
		Tags:  user.Tags,

		Repo:         u,
		ID:           user.ID,
		LastReadNoti: uid.UID(user.LastReadNoti),
	}
	// TODO: should in service level?
	if len(user.Tags) == 0 {
		entity.Tags = mainTags
	}
	return entity
}

func (u *UserRepo) GetOrInsertUser(ctx context.Context, email string) (*entity.User, bool, error) {
	var users []User
	if err := u.db(ctx).Model(&users).Where("email = ?", email).Select(); err != nil {
		return nil, false, err
	}
	mainTags, err := u.Forum.GetMainTags(ctx)
	if err != nil {
		return nil, false, err
	}
	if len(users) > 0 {
		return u.toEntityUser(&users[0], mainTags), false, nil
	}
	user := User{
		Email: email,
	}
	if _, err := u.db(ctx).Model(&user).Returning("*").Insert(); err != nil {
		return nil, false, err
	}
	return u.toEntityUser(&user, mainTags), true, nil
}

func (u *UserRepo) UpdateUser(ctx context.Context, id int64, update *entity.UserUpdate) error {
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
	return err
}
