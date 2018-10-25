package resolver

import (
	"context"

	graphql "github.com/graph-gophers/graphql-go"
	"gitlab.com/abyss.club/uexky/model"
)

// queries:

// Post ...
func (r *Resolver) Post(
	ctx context.Context, args struct{ ID string },
) (*PostResolver, error) {
	post, err := model.FindPost(ctx, args.ID)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, nil
	}
	return &PostResolver{Post: post}, nil
}

// mutations:

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

// types:

// PostSliceResolver ...
type PostSliceResolver struct {
	posts     []*PostResolver
	sliceInfo *SliceInfoResolver
}

// Posts ...
func (psr *PostSliceResolver) Posts(ctx context.Context) ([]*PostResolver, error) {
	return psr.posts, nil
}

// SliceInfo ...
func (psr *PostSliceResolver) SliceInfo(ctx context.Context) (*SliceInfoResolver, error) {
	return psr.sliceInfo, nil
}

// PostResolver ...
type PostResolver struct {
	Post *model.Post
}

// ID ...
func (pr *PostResolver) ID(ctx context.Context) (string, error) {
	return pr.Post.ID, nil
}

// Anonymous ...
func (pr *PostResolver) Anonymous(ctx context.Context) (bool, error) {
	return pr.Post.Anonymous, nil
}

// Author ...
func (pr *PostResolver) Author(ctx context.Context) (string, error) {
	return pr.Post.Author, nil
}

// Content ...
func (pr *PostResolver) Content(ctx context.Context) (string, error) {
	return pr.Post.Content, nil
}

// CreateTime ...
func (pr *PostResolver) CreateTime(ctx context.Context) (graphql.Time, error) {
	return graphql.Time{Time: pr.Post.CreateTime}, nil
}

// Quotes ...
func (pr *PostResolver) Quotes(ctx context.Context) (*[]*PostResolver, error) {
	quotes, err := pr.Post.QuotePosts(ctx)
	if err != nil {
		return nil, err
	}
	if len(quotes) == 0 {
		return nil, nil
	}
	var rps []*PostResolver
	for _, q := range quotes {
		rps = append(rps, &PostResolver{Post: q})
	}
	return &rps, nil
}

// QuoteCount ...
func (pr *PostResolver) QuoteCount(ctx context.Context) (int32, error) {
	count, err := pr.Post.QuoteCount(ctx)
	return int32(count), err
}
