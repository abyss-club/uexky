package resolver

import (
	"context"

	"gitlab.com/abyss.club/uexky/model"
	"gitlab.com/abyss.club/uexky/mw"
)

// Resolver for graphql
type Resolver struct{}

// Query:

// Profile resolve query 'profile'
func (r *Resolver) Profile(ctx context.Context) (*UserResolver, error) {
	user, err := model.GetUser(ctx) // TODO: not login, return nil
	return &UserResolver{user}, err
}

// ThreadSlice ...
func (r *Resolver) ThreadSlice(ctx context.Context, args struct {
	Tags  *[]string
	Query *SliceQuery
}) (
	*ThreadSliceResolver, error,
) {
	sq, err := args.Query.Parse(true)
	if err != nil {
		return nil, err
	}

	var tags []string
	if args.Tags != nil {
		tags = *(args.Tags)
	}

	threads, sliceInfo, err := model.GetThreadsByTags(ctx, tags, sq)
	if err != nil {
		return nil, err
	}

	var trs []*ThreadResolver
	for _, t := range threads {
		trs = append(trs, &ThreadResolver{Thread: t})
	}
	sir := &SliceInfoResolver{SliceInfo: sliceInfo}
	return &ThreadSliceResolver{threads: trs, sliceInfo: sir}, nil
}

// Thread ...
func (r *Resolver) Thread(ctx context.Context, args struct{ ID string }) (*ThreadResolver, error) {
	th, err := model.FindThread(ctx, args.ID)
	if err != nil {
		return nil, err
	}
	if th == nil {
		return nil, nil
	}
	return &ThreadResolver{Thread: th}, nil
}

// Post ...
func (r *Resolver) Post(ctx context.Context, args struct{ ID string }) (*PostResolver, error) {
	post, err := model.FindPost(ctx, args.ID)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, nil
	}
	return &PostResolver{Post: post}, nil
}

// Tags ...
func (r *Resolver) Tags(ctx context.Context, args struct{ Query string }) (*TagResolver, error) {
	return nil, nil // TODO
}

// Mutation:

// Auth ...
func (r *Resolver) Auth(ctx context.Context, args struct{ Email string }) (bool, error) {
	_, ok := ctx.Value(mw.ContextKeyEmail).(string)
	if ok {
		return false, nil
	}

	authURL, err := authEmail(ctx, args.Email)
	if err != nil {
		return false, nil
	}
	if err := sendAuthMail(authURL, args.Email); err != nil {
		return false, err
	}
	return true, nil
}

// SetName ...
func (r *Resolver) SetName(ctx context.Context, args struct{ Name string }) (*UserResolver, error) {
	user, err := model.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	if err := user.SetName(ctx, args.Name); err != nil {
		return nil, err
	}
	return &UserResolver{User: user}, nil
}

// SyncTags ...
func (r *Resolver) SyncTags(
	ctx context.Context, args struct{ Tags []*string },
) (*UserResolver, error) {
	user, err := model.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	tags := []string{}
	for _, t := range args.Tags {
		if t != nil {
			tags = append(tags, *t)
		}
	}
	if err := user.SyncTags(ctx, tags); err != nil {
		return nil, err
	}
	return &UserResolver{User: user}, nil
}

// PubThread ...
func (r *Resolver) PubThread(
	ctx context.Context,
	args struct{ Thread *model.ThreadInput },
) (
	*ThreadResolver, error,
) {
	thread, err := model.NewThread(ctx, args.Thread)
	if err != nil {
		return nil, err
	}
	return &ThreadResolver{Thread: thread}, nil
}

// PubPost ...
func (r *Resolver) PubPost(
	ctx context.Context,
	args struct{ Post *model.PostInput },
) (
	*PostResolver, error,
) {
	post, err := model.NewPost(ctx, args.Post)
	if err != nil {
		return nil, err
	}
	return &PostResolver{Post: post}, nil
}
