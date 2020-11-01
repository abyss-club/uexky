package repo

import (
	"context"

	"github.com/go-pg/pg/v9/orm"
	"github.com/go-redis/redis/v7"
	"gitlab.com/abyss.club/uexky/lib/postgres"
	"gitlab.com/abyss.club/uexky/uexky/entity"
)

func NewRepo(r *redis.Client) *entity.Repo {
	return &entity.Repo{
		User:   &UserRepo{Redis: r},
		Thread: &ThreadRepo{Redis: r},
		Post:   &PostRepo{Redis: r},
		Tag:    &TagRepo{},
		Noti:   &NotiRepo{},
	}
}

type queryFunc func(prev *orm.Query) *orm.Query

func db(ctx context.Context) postgres.Session {
	return postgres.GetSessionFromContext(ctx)
}
