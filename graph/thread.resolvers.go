package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"gitlab.com/abyss.club/uexky/uexky/types"
)

func (r *mutationResolver) PubThread(ctx context.Context, thread types.ThreadInput) (*types.Thread, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ThreadSlice(ctx context.Context, tags []string, query types.SliceQuery) (*types.ThreadSlice, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Thread(ctx context.Context, id string) (*types.Thread, error) {
	panic(fmt.Errorf("not implemented"))
}
