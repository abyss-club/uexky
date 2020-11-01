package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"gitlab.com/abyss.club/uexky/graph/generated"
	"gitlab.com/abyss.club/uexky/uexky/entity"
)

func (r *queryResolver) UnreadNotiCount(ctx context.Context) (int, error) {
	return r.Uexky.GetUnreadNotiCount(ctx)
}

func (r *queryResolver) Notification(ctx context.Context, query entity.SliceQuery) (*entity.NotiSlice, error) {
	return r.Uexky.GetNotifications(ctx, query)
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
