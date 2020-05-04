package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"gitlab.com/abyss.club/uexky/graph/generated"
	"gitlab.com/abyss.club/uexky/uexky/entity"
)

func (r *queryResolver) UnreadNotiCount(ctx context.Context) (*entity.UnreadNotiCount, error) {
	return r.Service.GetUnreadNotiCount(ctx)
}

func (r *queryResolver) Notification(ctx context.Context, typeArg entity.NotiType, query entity.SliceQuery) (*entity.NotiSlice, error) {
	return r.Service.GetNotification(ctx, typeArg, query)
}

func (r *quotedNotiResolver) HasRead(ctx context.Context, obj *entity.QuotedNoti) (bool, error) {
	return r.Service.GetQuotedNotiHasRead(ctx, obj)
}

func (r *quotedNotiResolver) Thread(ctx context.Context, obj *entity.QuotedNoti) (*entity.Thread, error) {
	return r.Service.GetQuotedNotiThread(ctx, obj)
}

func (r *quotedNotiResolver) QuotedPost(ctx context.Context, obj *entity.QuotedNoti) (*entity.Post, error) {
	return r.Service.GetQuotedNotiPost(ctx, obj)
}

func (r *quotedNotiResolver) Post(ctx context.Context, obj *entity.QuotedNoti) (*entity.Post, error) {
	return r.Service.GetQuotedNotiPost(ctx, obj)
}

func (r *repliedNotiResolver) HasRead(ctx context.Context, obj *entity.RepliedNoti) (bool, error) {
	return r.Service.GetRepliedNotiHasRead(ctx, obj)
}

func (r *repliedNotiResolver) Thread(ctx context.Context, obj *entity.RepliedNoti) (*entity.Thread, error) {
	return r.Service.GetRepliedNotiThread(ctx, obj)
}

func (r *repliedNotiResolver) Repliers(ctx context.Context, obj *entity.RepliedNoti) ([]string, error) {
	return r.Service.GetRepliedNotiRepliers(ctx, obj)
}

func (r *systemNotiResolver) HasRead(ctx context.Context, obj *entity.SystemNoti) (bool, error) {
	return r.Service.GetSystemNotiHasRead(ctx, obj)
}

func (r *systemNotiResolver) Content(ctx context.Context, obj *entity.SystemNoti) (string, error) {
	return r.Service.GetSystemNotiContent(ctx, obj)
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// QuotedNoti returns generated.QuotedNotiResolver implementation.
func (r *Resolver) QuotedNoti() generated.QuotedNotiResolver { return &quotedNotiResolver{r} }

// RepliedNoti returns generated.RepliedNotiResolver implementation.
func (r *Resolver) RepliedNoti() generated.RepliedNotiResolver { return &repliedNotiResolver{r} }

// SystemNoti returns generated.SystemNotiResolver implementation.
func (r *Resolver) SystemNoti() generated.SystemNotiResolver { return &systemNotiResolver{r} }

type queryResolver struct{ *Resolver }
type quotedNotiResolver struct{ *Resolver }
type repliedNotiResolver struct{ *Resolver }
type systemNotiResolver struct{ *Resolver }
