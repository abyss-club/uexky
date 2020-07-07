package entity

import "context"

type NotiSearch struct {
	UserID int
}

type NotiRepo interface {
	GetUserUnreadCount(ctx context.Context, user *User) (int, error)
	GetNotiSlice(ctx context.Context, search *NotiSearch, query SliceQuery) (*NotiSlice, error)
	InsertNoti(ctx context.Context, insert *Notification) error
	UpdateReadID(ctx context.Context, userID int, id int) error
}
