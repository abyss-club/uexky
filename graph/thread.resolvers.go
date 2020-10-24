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

func (r *mutationResolver) PubThread(ctx context.Context, thread entity.ThreadInput) (*entity.Thread, error) {
	return r.Uexky.PubThread(ctx, thread)
}

func (r *queryResolver) ThreadSlice(ctx context.Context, tags []string, query entity.SliceQuery) (*entity.ThreadSlice, error) {
	return r.Uexky.SearchThreads(ctx, tags, query)
}

func (r *queryResolver) Thread(ctx context.Context, id uid.UID) (*entity.Thread, error) {
	return r.Uexky.GetThreadByID(ctx, id)
}

func (r *threadResolver) Replies(ctx context.Context, obj *entity.Thread, query entity.SliceQuery) (*entity.PostSlice, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *threadResolver) ReplyCount(ctx context.Context, obj *entity.Thread) (int, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *threadResolver) Catalog(ctx context.Context, obj *entity.Thread) ([]*entity.ThreadCatalogItem, error) {
	panic(fmt.Errorf("not implemented"))
}

// Thread returns generated.ThreadResolver implementation.
func (r *Resolver) Thread() generated.ThreadResolver { return &threadResolver{r} }

type threadResolver struct{ *Resolver }
