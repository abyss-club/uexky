package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"gitlab.com/abyss.club/uexky/uexky/types"
)

func (r *queryResolver) MainTags(ctx context.Context) ([]string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Recommended(ctx context.Context) ([]string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Tags(ctx context.Context, query *string, limit *int) ([]*types.Tag, error) {
	panic(fmt.Errorf("not implemented"))
}
