package resolver

import (
	"context"
	"time"

	"gitlab.com/abyss.club/uexky/model"
)

// queries:

// UnreadNotifCount ...
func (r *Resolver) UnreadNotifCount(ctx context.Context) (*UnreadResolver, error) {
	return &UnreadResolver{}, nil
}

// Notification resolve query 'notification'
func (r *Resolver) Notification(ctx context.Context, args struct {
	Type  string
	Query *SliceQuery
}) (*NotiSliceResolver, error) {
	sq, err := args.Query.Parse(true)
	if err != nil {
		return nil, err
	}
	notiType := model.NotiType(args.Type)

	noti, sliceInfo, err := model.GetNotification(ctx, notiType, sq)
	if err != nil {
		return nil, err
	}
	return &NotiSliceResolver{notiType, noti, sliceInfo}, nil
}

// types:

// UnreadResolver ...
type UnreadResolver struct{}

// System ...
func (ur *UnreadResolver) System(ctx context.Context) (int32, error) {
	count, err := model.GetUnreadNotificationCount(ctx, model.NotiTypeSystem)
	return int32(count), err
}

// Replied ...
func (ur *UnreadResolver) Replied(ctx context.Context) (int32, error) {
	count, err := model.GetUnreadNotificationCount(ctx, model.NotiTypeReplied)
	return int32(count), err
}

// Refered ...
func (ur *UnreadResolver) Refered(ctx context.Context) (int32, error) {
	count, err := model.GetUnreadNotificationCount(ctx, model.NotiTypeRefered)
	return int32(count), err
}

// NotiSliceResolver ...
type NotiSliceResolver struct {
	notiType  model.NotiType
	notiSlice []*model.NotiStore
	sliceInfo *model.SliceInfo
}

// System ...
func (nsr *NotiSliceResolver) System(ctx context.Context) (
	[]*SystemNotiResolver, error,
) {
	if nsr.notiType != model.NotiTypeSystem {
		return nil, nil
	}
	snrs := []*SystemNotiResolver{}
	for _, n := range nsr.notiSlice {
		snrs = append(snrs, &SystemNotiResolver{n.GetSystemNoti()})
	}
	return snrs, nil
}

// Replied ...
func (nsr *NotiSliceResolver) Replied(ctx context.Context) (
	[]*RepliedNotiResolver, error,
) {
	if nsr.notiType != model.NotiTypeReplied {
		return nil, nil
	}
	rnrs := []*RepliedNotiResolver{}
	for _, n := range nsr.notiSlice {
		rnrs = append(rnrs, &RepliedNotiResolver{n.GetRepliedNoti()})
	}
	return rnrs, nil
}

// Refered ...
func (nsr *NotiSliceResolver) Refered(ctx context.Context) (
	[]*ReferedNotiResolver, error,
) {
	if nsr.notiType != model.NotiTypeRefered {
		return nil, nil
	}
	rnrs := []*ReferedNotiResolver{}
	for _, n := range nsr.notiSlice {
		rnrs = append(rnrs, &ReferedNotiResolver{n.GetReferedNoti()})
	}
	return rnrs, nil
}

// SliceInfo ...
func (nsr *NotiSliceResolver) SliceInfo(ctx context.Context) (*SliceInfoResolver, error) {
	return &SliceInfoResolver{nsr.sliceInfo}, nil
}

// SystemNotiResolver ...
type SystemNotiResolver struct {
	noti *model.SystemNoti
}

// ID ...
func (n *SystemNotiResolver) ID(ctx context.Context) (string, error) {
	return n.noti.ID, nil
}

// Type ...
func (n *SystemNotiResolver) Type(ctx context.Context) (string, error) {
	return string(n.noti.Type), nil
}

// EventTime ...
func (n *SystemNotiResolver) EventTime(ctx context.Context) (time.Time, error) {
	return n.noti.EventTime, nil
}

// HasRead ...
func (n *SystemNotiResolver) HasRead(ctx context.Context) (bool, error) {
	return n.noti.HasRead, nil
}

// Title ...
func (n *SystemNotiResolver) Title(ctx context.Context) (string, error) {
	return n.noti.Title, nil
}

// Content ...
func (n *SystemNotiResolver) Content(ctx context.Context) (string, error) {
	return n.noti.Content, nil
}

// RepliedNotiResolver ...
type RepliedNotiResolver struct {
	noti *model.RepliedNoti
}

// ID ...
func (n *RepliedNotiResolver) ID(ctx context.Context) (string, error) {
	return n.noti.ID, nil
}

// Type ...
func (n *RepliedNotiResolver) Type(ctx context.Context) (string, error) {
	return string(n.noti.Type), nil
}

// EventTime ...
func (n *RepliedNotiResolver) EventTime(ctx context.Context) (time.Time, error) {
	return n.noti.EventTime, nil
}

// HasRead ...
func (n *RepliedNotiResolver) HasRead(ctx context.Context) (bool, error) {
	return n.noti.HasRead, nil
}

// Thread ...
func (n *RepliedNotiResolver) Thread(ctx context.Context) (string, error) {
	return n.noti.Thread, nil // TODO
}

// Repliers ...
func (n *RepliedNotiResolver) Repliers(ctx context.Context) ([]string, error) {
	return n.noti.Repliers, nil
}

// ReferedNotiResolver ...
type ReferedNotiResolver struct {
	noti *model.ReferedNoti
}

// ID ...
func (n *ReferedNotiResolver) ID(ctx context.Context) (string, error) {
	return n.noti.ID, nil
}

// Type ...
func (n *ReferedNotiResolver) Type(ctx context.Context) (string, error) {
	return string(n.noti.Type), nil
}

// EventTime ...
func (n *ReferedNotiResolver) EventTime(ctx context.Context) (time.Time, error) {
	return n.noti.EventTime, nil
}

// HasRead ...
func (n *ReferedNotiResolver) HasRead(ctx context.Context) (bool, error) {
	return n.noti.HasRead, nil
}

// Thread ...
func (n *ReferedNotiResolver) Thread(ctx context.Context) (string, error) {
	return n.noti.Thread, nil // TODO
}

// Post ...
func (n *ReferedNotiResolver) Post(ctx context.Context) (string, error) {
	return n.noti.Post, nil // TODO
}

// Referers ...
func (n *ReferedNotiResolver) Referers(ctx context.Context) (string, error) {
	return n.noti.Referers, nil
}
