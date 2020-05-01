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

func (r *queryResolver) Notification(ctx context.Context, typeArg string, query entity.SliceQuery) (*entity.NotiSlice, error) {
	return r.Service.GetNotification(ctx, typeArg, query)
}

func (r *systemNotiResolver) Content(ctx context.Context, obj *entity.SystemNoti) (string, error) {
	return obj.ContentText(ctx)
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// SystemNoti returns generated.SystemNotiResolver implementation.
func (r *Resolver) SystemNoti() generated.SystemNotiResolver { return &systemNotiResolver{r} }

type queryResolver struct{ *Resolver }
type systemNotiResolver struct{ *Resolver }
