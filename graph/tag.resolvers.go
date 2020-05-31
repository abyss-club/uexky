package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"gitlab.com/abyss.club/uexky/uexky/entity"
)

func (r *queryResolver) MainTags(ctx context.Context) ([]string, error) {
	return r.Service.GetMainTags(ctx)
}

func (r *queryResolver) Recommended(ctx context.Context) ([]string, error) {
	return r.Service.GetRecommendedTags(ctx)
}

func (r *queryResolver) Tags(ctx context.Context, query *string, limit *int) ([]*entity.Tag, error) {
	return r.Service.SearchTags(ctx, query, limit)
}
