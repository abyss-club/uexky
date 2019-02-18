package resolver

import (
	"context"

	"gitlab.com/abyss.club/uexky-go/model"
	"gitlab.com/abyss.club/uexky-go/uexky"
)

// queries:

// Profile resolve query 'profile'
func (r *Resolver) Profile(ctx context.Context) (*UserResolver, error) {
	u := uexky.Pop(ctx)
	user, err := model.GetSignedInUser(u)
	if err != nil { // not login, return null user
		return &UserResolver{&model.User{}}, nil
	}
	return &UserResolver{user}, nil
}

// mutations:

// Auth ...
func (r *Resolver) Auth(
	ctx context.Context, args struct{ Email string },
) (bool, error) {
	u := uexky.Pop(ctx)
	if u.Auth.IsSignedIn() {
		return false, nil
	}

	authURL, err := authEmail(u, args.Email)
	if err != nil {
		return false, nil
	}
	if err := sendAuthMail(authURL, args.Email); err != nil {
		return false, err
	}
	return true, nil
}

// SetName ...
func (r *Resolver) SetName(
	ctx context.Context, args struct{ Name string },
) (*UserResolver, error) {
	u := uexky.Pop(ctx)
	user, err := model.GetSignedInUser(u)
	if err != nil {
		return nil, err
	}
	if err := user.SetName(u, args.Name); err != nil {
		return nil, err
	}
	return &UserResolver{User: user}, nil
}

// SyncTags ...
func (r *Resolver) SyncTags(
	ctx context.Context, args struct{ Tags []*string },
) (*UserResolver, error) {
	u := uexky.Pop(ctx)
	user, err := model.GetSignedInUser(u)
	if err != nil {
		return nil, err
	}
	tags := []string{}
	for _, t := range args.Tags {
		if t != nil {
			tags = append(tags, *t)
		}
	}
	if err := user.SyncTags(u, tags); err != nil {
		return nil, err
	}
	return &UserResolver{User: user}, nil
}

// AddSubbedTags ...
func (r *Resolver) AddSubbedTags(
	ctx context.Context, args struct{ Tags []string },
) (*UserResolver, error) {
	u := uexky.Pop(ctx)
	user, err := model.GetSignedInUser(u)
	if err != nil {
		return nil, err
	}
	if err := user.AddSubbedTags(u, args.Tags); err != nil {
		return nil, err
	}
	return &UserResolver{User: user}, nil
}

// DelSubbedTags ...
func (r *Resolver) DelSubbedTags(
	ctx context.Context, args struct{ Tags []string },
) (*UserResolver, error) {
	u := uexky.Pop(ctx)
	user, err := model.GetSignedInUser(u)
	if err != nil {
		return nil, err
	}
	if err := user.DelSubbedTags(u, args.Tags); err != nil {
		return nil, err
	}
	return &UserResolver{User: user}, nil
}

// BanUser ...
func (r *Resolver) BanUser(
	ctx context.Context, args struct{ PostID string },
) error {
	u := uexky.Pop(ctx)
	return model.BanUser(u, args.PostID)
}

// BlockPost ...
func (r *Resolver) BlockPost(
	ctx context.Context, args struct{ PostID string },
) error {
	u := uexky.Pop(ctx)
	return model.BlockPost(u, args.PostID)
}

// LockThread ...
func (r *Resolver) LockThread(
	ctx context.Context, args struct{ ThreadID string },
) error {
	u := uexky.Pop(ctx)
	return model.LockThread(u, args.ThreadID)
}

// BlockThread ...
func (r *Resolver) BlockThread(
	ctx context.Context, args struct{ ThreadID string },
) error {
	u := uexky.Pop(ctx)
	return model.BlockThread(u, args.ThreadID)
}

// EditTags ...
func (r *Resolver) EditTags(
	ctx context.Context,
	args struct {
		ThreadID string
		MainTag  string
		SubTags  []string
	},
) error {
	u := uexky.Pop(ctx)
	return model.EditTags(u, args.ThreadID, args.MainTag, args.SubTags)
}

// types:

// UserResolver for graphql
type UserResolver struct {
	User *model.User
}

// Email resolve user.email
func (ur *UserResolver) Email(ctx context.Context) (string, error) {
	return ur.User.Email, nil
}

// Name resolve user.name
func (ur *UserResolver) Name(ctx context.Context) (*string, error) {
	if ur.User.Name == "" {
		return nil, nil
	}
	return &(ur.User.Name), nil
}

// Tags ...
func (ur *UserResolver) Tags(ctx context.Context) (*[]string, error) {
	if len(ur.User.Tags) == 0 {
		return nil, nil
	}
	return &(ur.User.Tags), nil
}
