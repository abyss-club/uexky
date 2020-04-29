package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"gitlab.com/abyss.club/uexky/graph/generated"
	"gitlab.com/abyss.club/uexky/service"
)

func (r *mutationResolver) SetName(ctx context.Context, id int, name *string) (*service.User, error) {
	return r.Service.SetName(ctx, id, name)
}

func (r *mutationResolver) SetLevel(ctx context.Context, id int, level int) (*service.User, error) {
	return r.Service.SetLevel(ctx, id, level)
}

func (r *queryResolver) User(ctx context.Context, id int) (*service.User, error) {
	return r.Service.GetUser(ctx, id)
}

func (r *userResolver) Friends(ctx context.Context, obj *service.User) ([]*service.User, error) {
	return obj.GetFriends(ctx)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
