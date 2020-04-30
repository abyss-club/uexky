package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"gitlab.com/abyss.club/uexky/uexky/types"
)

func (r *mutationResolver) Auth(ctx context.Context, email string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SetName(ctx context.Context, name string) (*types.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SyncTags(ctx context.Context, tags []*string) (*types.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddSubbedTag(ctx context.Context, tag string) (*types.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DelSubbedTag(ctx context.Context, tag string) (*types.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) BanUser(ctx context.Context, postID *string, threadID *string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) BlockPost(ctx context.Context, postID string) (*types.Post, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) LockThread(ctx context.Context, threadID string) (*types.Thread, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) BlockThread(ctx context.Context, threadID string) (*types.Thread, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) EditTags(ctx context.Context, threadID string, mainTag string, subTags []string) (*types.Thread, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Profile(ctx context.Context) (*types.User, error) {
	panic(fmt.Errorf("not implemented"))
}
