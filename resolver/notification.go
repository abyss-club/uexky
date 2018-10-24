package resolver

import (
	"context"

	graphql "github.com/graph-gophers/graphql-go"
	"gitlab.com/abyss.club/uexky/model"
)

// queries:

// UnreadNotiCount ...
func (r *Resolver) UnreadNotiCount(ctx context.Context) (*UnreadResolver, error) {
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
	*[]*SystemNotiResolver, error,
) {
	if nsr.notiType != model.NotiTypeSystem {
		return nil, nil
	}
	snrs := []*SystemNotiResolver{}
	for _, n := range nsr.notiSlice {
		snrs = append(snrs, &SystemNotiResolver{notiBaseResolver{
			notiType: nsr.notiType,
			noti:     n,
		}})
	}
	return &snrs, nil
}

// Replied ...
func (nsr *NotiSliceResolver) Replied(ctx context.Context) (
	*[]*RepliedNotiResolver, error,
) {
	if nsr.notiType != model.NotiTypeReplied {
		return nil, nil
	}
	rnrs := []*RepliedNotiResolver{}
	for _, n := range nsr.notiSlice {
		rnrs = append(rnrs, &RepliedNotiResolver{notiBaseResolver{
			notiType: nsr.notiType,
			noti:     n,
		}})
	}
	return &rnrs, nil
}

// Refered ...
func (nsr *NotiSliceResolver) Refered(ctx context.Context) (
	*[]*ReferedNotiResolver, error,
) {
	if nsr.notiType != model.NotiTypeRefered {
		return nil, nil
	}
	rnrs := []*ReferedNotiResolver{}
	for _, n := range nsr.notiSlice {
		rnrs = append(rnrs, &ReferedNotiResolver{notiBaseResolver{
			notiType: nsr.notiType,
			noti:     n,
		}})
	}
	return &rnrs, nil
}

// SliceInfo ...
func (nsr *NotiSliceResolver) SliceInfo(ctx context.Context) (*SliceInfoResolver, error) {
	return &SliceInfoResolver{nsr.sliceInfo}, nil
}

type notiBaseResolver struct {
	notiType model.NotiType
	noti     *model.NotiStore
}

// ID ...
func (n *notiBaseResolver) ID(ctx context.Context) (string, error) {
	return n.noti.ID, nil
}

// Type ...
func (n *notiBaseResolver) Type(ctx context.Context) (string, error) {
	return string(n.noti.Type), nil
}

// EventTime ...
func (n *notiBaseResolver) EventTime(ctx context.Context) (graphql.Time, error) {
	return graphql.Time{Time: n.noti.EventTime}, nil
}

// HasRead ...
func (n *notiBaseResolver) HasRead(ctx context.Context) (bool, error) {
	return n.noti.HasRead, nil
}

// SystemNotiResolver ...
type SystemNotiResolver struct {
	notiBaseResolver
}

// Title ...
func (n *SystemNotiResolver) Title(ctx context.Context) (string, error) {
	return n.noti.System.Title, nil
}

// Content ...
func (n *SystemNotiResolver) Content(ctx context.Context) (string, error) {
	return n.noti.System.Content, nil
}

// RepliedNotiResolver ...
type RepliedNotiResolver struct {
	notiBaseResolver
}

// Thread ...
func (n *RepliedNotiResolver) Thread(ctx context.Context) (*ThreadResolver, error) {
	thread, err := model.FindThread(ctx, n.noti.Replied.ThreadID)
	if err != nil {
		return nil, err
	}
	return &ThreadResolver{Thread: thread}, nil
}

// Repliers ...
func (n *RepliedNotiResolver) Repliers(ctx context.Context) ([]string, error) {
	return n.noti.Replied.Repliers, nil
}

// ReferedNotiResolver ...
type ReferedNotiResolver struct {
	notiBaseResolver
}

// Thread ...
func (n *ReferedNotiResolver) Thread(ctx context.Context) (*ThreadResolver, error) {
	thread, err := model.FindThread(ctx, n.noti.Refered.ThreadID)
	if err != nil {
		return nil, err
	}
	return &ThreadResolver{thread}, nil
}

// Post ...
func (n *ReferedNotiResolver) Post(ctx context.Context) (*PostResolver, error) {
	post, err := model.FindPost(ctx, n.noti.Refered.PostID)
	if err != nil {
		return nil, err
	}
	return &PostResolver{post}, nil
}

// Referers ...
func (n *ReferedNotiResolver) Referers(ctx context.Context) ([]string, error) {
	return n.noti.Refered.Referers, nil
}
