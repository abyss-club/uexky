package resolver

import (
	"context"
	"time"

	"gitlab.com/abyss.club/uexky/model"
)

// queries:

// UnreadNotifCount ...
func (r *Resolver) UnreadNotifCount(
	ctx context.Context,
	args struct{ Type *string },
) (int32, error) {
	var typeStr model.NotifType
	if args.Type != nil {
		typeStr = model.NotifType(*(args.Type))
	}
	count, err := model.GetUnreadNotificationCount(ctx, typeStr)
	return int32(count), err
}

// Notification resolve query 'notification'
func (r *Resolver) Notification(ctx context.Context, args struct {
	Type  string
	Query *SliceQuery
}) (*NotifSliceResolver, error) {
	sq, err := args.Query.Parse(true)
	if err != nil {
		return nil, err
	}
	typeStr := model.NotifType(args.Type)

	notif, sliceInfo, err := model.GetNotificationByUser(ctx, typeStr, sq)
	if err != nil {
		return nil, err
	}

	nrs := []*NotifResolver{}
	for _, n := range notif {
		nrs = append(nrs, &NotifResolver{n})
	}
	sir := &SliceInfoResolver{sliceInfo}
	return &NotifSliceResolver{nrs, sir}, nil
}

// types:

// NotifSliceResolver ...
type NotifSliceResolver struct {
	notifSlice []*NotifResolver
	sliceInfo  *SliceInfoResolver
}

// Notif ...
func (nsr *NotifSliceResolver) Notif(ctx context.Context) ([]*NotifResolver, error) {
	return nsr.notifSlice, nil
}

// SliceInfo ...
func (nsr *NotifSliceResolver) SliceInfo(ctx context.Context) (*SliceInfoResolver, error) {
	return nsr.sliceInfo, nil
}

// NotifResolver ...
type NotifResolver struct {
	notif *model.Notification
}

// Type ...
func (nr *NotifResolver) Type(ctx context.Context) (string, error) {
	return string(nr.notif.Type), nil
}

// ReleaseTime ...
func (nr *NotifResolver) ReleaseTime(ctx context.Context) (time.Time, error) {
	return nr.notif.ReleaseTime, nil
}

// HasRead ...
func (nr *NotifResolver) HasRead(ctx context.Context) (bool, error) {
	user, err := model.GetUser(ctx)
	if err != nil {
		return false, nil
	}
	return nr.notif.ReleaseTime.After(user.ReadNotifTime), nil
}

// Content ...
func (nr *NotifResolver) Content(ctx context.Context) (string, error) {
	return nr.notif.GetContent(), nil
}
