package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"gitlab.com/abyss.club/uexky/graph/generated"
	"gitlab.com/abyss.club/uexky/lib/uid"
	"gitlab.com/abyss.club/uexky/uexky/entity"
)

func (r *mutationResolver) PubPost(ctx context.Context, post entity.PostInput) (*entity.Post, error) {
	if err := r.TxAdapter.Begin(ctx); err != nil {
		return nil, err
	}
	newPost, err := r.Service.PubPost(ctx, post)
	return newPost, r.TxAdapter.Rollback(ctx, err)
}

func (r *queryResolver) Post(ctx context.Context, id uid.UID) (*entity.Post, error) {
	return r.Service.GetPostByID(ctx, id)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
