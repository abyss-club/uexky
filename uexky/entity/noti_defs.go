package entity

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky/lib/uid"
)

//-- NotiRepo

type NotiRepo interface {
	GetUserUnreadCount(ctx context.Context, user *User) (int, error)
	GetNotiByKey(ctx context.Context, user *User, key string) (*Notification, error)
	GetNotiSlice(ctx context.Context, user *User, query SliceQuery) (*NotiSlice, error)
	InsertNoti(ctx context.Context, notification *Notification) error
	UpdateNotiContent(ctx context.Context, noti *Notification) error
	UpdateReadID(ctx context.Context, userID int64, id uid.UID) error
}

// -- NotiType

func ParseNotiType(s string) (NotiType, error) {
	t := NotiType(s)
	if !t.IsValid() {
		return NotiType(""), errors.Errorf("invalid noti type: %s", s)
	}
	return t, nil
}

// -- Receiver

type SendGroup string

const AllUser SendGroup = "all_user"

type Receiver string

func SendToUser(userID int64) Receiver {
	return Receiver(fmt.Sprintf("u:%v", userID))
}

func SendToGroup(group SendGroup) Receiver {
	return Receiver(fmt.Sprintf("g:%v", group))
}
