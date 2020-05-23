package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"gitlab.com/abyss.club/uexky/lib/uid"
	"gitlab.com/abyss.club/uexky/uexky/entity"
)

func (r *mutationResolver) PubThread(ctx context.Context, thread entity.ThreadInput) (*entity.Thread, error) {
	if err := r.TxAdapter.Begin(ctx); err != nil {
		return nil, err
	}
	newThread, err := r.Service.PubThread(ctx, thread)
	return newThread, r.TxAdapter.Rollback(ctx, err)
}

func (r *queryResolver) ThreadSlice(ctx context.Context, tags []string, query entity.SliceQuery) (*entity.ThreadSlice, error) {
	return r.Service.SearchThreads(ctx, tags, query)
}

func (r *queryResolver) Thread(ctx context.Context, id uid.UID) (*entity.Thread, error) {
	return r.Service.GetThreadByID(ctx, id)
}
