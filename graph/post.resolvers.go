package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"gitlab.com/abyss.club/uexky/graph/generated"
	"gitlab.com/abyss.club/uexky/lib/uid"
	"gitlab.com/abyss.club/uexky/uexky/entity"
)

func (r *mutationResolver) PubPost(ctx context.Context, post entity.PostInput) (*entity.Post, error) {
	return r.Uexky.PubPost(ctx, post)
}

func (r *postResolver) Quotes(ctx context.Context, obj *entity.Post) ([]*entity.Post, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *postResolver) QuotedCount(ctx context.Context, obj *entity.Post) (int, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Post(ctx context.Context, id uid.UID) (*entity.Post, error) {
	return r.Uexky.GetPostByID(ctx, id)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Post returns generated.PostResolver implementation.
func (r *Resolver) Post() generated.PostResolver { return &postResolver{r} }

type mutationResolver struct{ *Resolver }
type postResolver struct{ *Resolver }
