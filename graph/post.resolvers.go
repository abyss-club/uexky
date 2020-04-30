package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"gitlab.com/abyss.club/uexky/graph/generated"
	"gitlab.com/abyss.club/uexky/uexky/types"
)

func (r *mutationResolver) PubPost(ctx context.Context, post types.PostInput) (*types.Post, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Post(ctx context.Context, id string) (*types.Post, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
