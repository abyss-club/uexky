package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"gitlab.com/abyss.club/uexky/graph/generated"
	"gitlab.com/abyss.club/uexky/lib/uid"
	"gitlab.com/abyss.club/uexky/uexky/entity"
)

func (r *mutationResolver) Auth(ctx context.Context, email string) (bool, error) {
	return r.Service.SignInByEmail(ctx, email)
}

func (r *mutationResolver) SetName(ctx context.Context, name string) (*entity.User, error) {
	return r.Service.SetUserName(ctx, name)
}

func (r *mutationResolver) SyncTags(ctx context.Context, tags []*string) (*entity.User, error) {
	return r.Service.SyncUserTags(ctx, tags)
}

func (r *mutationResolver) AddSubbedTag(ctx context.Context, tag string) (*entity.User, error) {
	return r.Service.AddUserSubbedTag(ctx, tag)
}

func (r *mutationResolver) DelSubbedTag(ctx context.Context, tag string) (*entity.User, error) {
	return r.Service.DelUserSubbedTag(ctx, tag)
}

func (r *mutationResolver) BanUser(ctx context.Context, postID *uid.UID, threadID *uid.UID) (bool, error) {
	return r.Service.BanUser(ctx, postID, threadID)
}

func (r *mutationResolver) BlockPost(ctx context.Context, postID uid.UID) (*entity.Post, error) {
	return r.Service.BlockPost(ctx, postID)
}

func (r *mutationResolver) LockThread(ctx context.Context, threadID uid.UID) (*entity.Thread, error) {
	return r.Service.LockThread(ctx, threadID)
}

func (r *mutationResolver) BlockThread(ctx context.Context, threadID uid.UID) (*entity.Thread, error) {
	return r.Service.BlockThread(ctx, threadID)
}

func (r *mutationResolver) EditTags(ctx context.Context, threadID uid.UID, mainTag string, subTags []string) (*entity.Thread, error) {
	return r.Service.EditTags(ctx, threadID, mainTag, subTags)
}

func (r *queryResolver) Profile(ctx context.Context) (*entity.User, error) {
	return r.Service.Profile(ctx)
}

func (r *userResolver) Tags(ctx context.Context, obj *entity.User) ([]string, error) {
	return r.Service.GetUserTags(ctx, obj)
}

func (r *userResolver) Threads(ctx context.Context, obj *entity.User, query entity.SliceQuery) (*entity.ThreadSlice, error) {
	return r.Service.GetUserThreads(ctx, obj, query)
}

func (r *userResolver) Posts(ctx context.Context, obj *entity.User, query entity.SliceQuery) (*entity.PostSlice, error) {
	return r.Service.GetUserPosts(ctx, obj, query)
}

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type userResolver struct{ *Resolver }
