package resolver

import (
	"context"

	"gitlab.com/abyss.club/uexky/model"
	"gitlab.com/abyss.club/uexky/mw"
	"gitlab.com/abyss.club/uexky/uexky"
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
func (r *Resolver) SetName(
	ctx context.Context, args struct{ Name string },
) (*UserResolver, error) {
	u := uexky.Pop(ctx)
	user, err := model.GetSignedInUser(u)
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
	if err := user.SyncTags(ctx, tags); err != nil {
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
	if err := user.AddSubbedTags(ctx, args.Tags); err != nil {
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
	if err := user.DelSubbedTags(ctx, args.Tags); err != nil {
		return nil, err
	}
	return &UserResolver{User: user}, nil
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
