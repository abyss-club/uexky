package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"gitlab.com/abyss.club/uexky/graph/generated"
	"gitlab.com/abyss.club/uexky/uexky/types"
)

func (r *queryResolver) UnreadNotiCount(ctx context.Context) (*types.UnreadNotiCount, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Notification(ctx context.Context, typeArg string, query types.SliceQuery) (*types.NotiSlice, error) {
	panic(fmt.Errorf("not implemented"))
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
