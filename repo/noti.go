package repo

import (
	"context"
	"fmt"

	"gitlab.com/abyss.club/uexky/lib/postgres"
	"gitlab.com/abyss.club/uexky/uexky/entity"
)

type NotiRepo struct{}

func (n *NotiRepo) db(ctx context.Context) postgres.Session {
	return postgres.GetSessionFromContext(ctx)
}

func (n *NotiRepo) GetUserUnreadCount(ctx context.Context, user *entity.User) (int, error) {
	panic(fmt.Errorf("not implemented"))
}

func (n *NotiRepo) GetNotiSlice(
	ctx context.Context, search *entity.NotiSearch, query entity.SliceQuery,
) (*entity.NotiSlice, error) {
	panic(fmt.Errorf("not implemented"))
}

func (n *NotiRepo) InsertNoti(ctx context.Context, insert *entity.Notification) error {
	panic(fmt.Errorf("not implemented"))
}

func (n *NotiRepo) UpdateReadID(ctx context.Context, userID int, id int) error {
	panic(fmt.Errorf("not implemented"))
}
