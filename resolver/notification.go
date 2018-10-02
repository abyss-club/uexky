package resolver

import (
	"context"
	"time"

	"gitlab.com/abyss.club/uexky/model"
)

// NotificationSliceResolver ...
type NotificationSliceResolver struct {
	Query        bool
	Notification []*model.Notification
	SliceInfo    *model.SliceInfo
	LastReadTime time.Time
}

func (nsr *NotificationSliceResolver) doQuery(ctx context.Context) error {
	if nsr.Query {
		return nil
	}
	user, err := model.GetUser(ctx)
	if err != nil {
		return err
	}
	nsr.LastReadTime = user.ReadNotifTime

	notifs, slideInfo := GetNotificationByUser(ctx)
}

// UnreadCount ...
func (nsr *NotificationSliceResolver) UnreadCount(
	ctx context.Context,
) (int32, error) {
	count, err := model.GetUnreadNotificationCount(ctx)
	return int32(count), err
}

// Notification ...
func (nsr *NotificationSliceResolver) Notification() {
}

// SliceInfo ...
func (nsr *NotificationSliceResolver) SliceInfo() {
}

// NotificationResolver ...
type NotificationResolver struct {
}

// ReleaseTime ...
func ReleaseTime(ctx context.Context) {
}

// Type ...
func Type(ctx context.Context) {
}

// HasRead ...
func HasRead(ctx context.Context) {
}

// Content ...
func Content(ctx context.Context) {
}
