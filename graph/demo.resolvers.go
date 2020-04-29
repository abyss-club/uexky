package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"gitlab.com/abyss.club/uexky/entity"
	"gitlab.com/abyss.club/uexky/graph/generated"
)

func (r *mutationResolver) SetName(ctx context.Context, id int, name *string) (*entity.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SetLevel(ctx context.Context, id int, level int) (*entity.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) User(ctx context.Context, name string) (*entity.User, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
