package entity

import "context"

type NotiSearch struct {
	UserID int
	Type   NotiType
}

type NotiInsert struct {
	System  *SystemNoti
	Replied *RepliedNoti
	Quoted  *QuotedNoti
}

type NotiRepo interface {
	GetUserUnreadCount(ctx context.Context, user *User) (*UnreadNotiCount, error)
	GetNotiSlice(ctx context.Context, search *NotiSearch, query SliceQuery) (*NotiSlice, error)
	InsertNoti(ctx context.Context, insert NotiInsert) error
	UpdateReadID(ctx context.Context, userID int, nType NotiType, id int) error
}
