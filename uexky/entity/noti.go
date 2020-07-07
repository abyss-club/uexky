// forum notificatoin: systemNoti repliedNoti quotedNoti

package entity

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

type NotiService struct {
	Repo NotiRepo
}

func (n *NotiService) GetUnreadNotiCount(ctx context.Context, user *User) (int, error) {
	return n.Repo.GetUserUnreadCount(ctx, user)
}

func (n *NotiService) GetNotification(ctx context.Context, user *User, query SliceQuery) (*NotiSlice, error) {
	panic(fmt.Errorf("not implemented"))
}

type Notification struct {
	Type      NotiType    `json:"type"`
	EventTime time.Time   `json:"eventTime"`
	Content   NotiContent `json:"content"`
}

func (n *Notification) HasRead(user *User) bool {
	panic(fmt.Errorf("not implemented"))
}

func ParseNotiType(s string) (NotiType, error) {
	t := NotiType(s)
	if !t.IsValid() {
		return NotiType(""), errors.Errorf("invalid noti type: %s", s)
	}
	return t, nil
}

type SendGroup string

const AllUser SendGroup = "all_user"

type SendTo struct {
	UserID    *int       `json:"send_to"`
	SendGroup *SendGroup `json:"send_to_group"`
}

func SendToUser(userID int) SendTo {
	return SendTo{UserID: &userID}
}

func SendToGourp(group SendGroup) SendTo {
	return SendTo{SendGroup: &group}
}

type SystemNotiContent struct {
	Title       string `json:"title"`
	ContentText string `json:"content"`
}
