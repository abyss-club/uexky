package api

import (
	"context"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/nanozuki/uexky/model"
)

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

// Refers ...
func (pr *PostResolver) Refers(ctx context.Context) (*[]*PostResolver, error) {
	refers, err := pr.Post.ReferPosts(ctx)
	if err != nil {
		return nil, err
	}
	if len(refers) == 0 {
		return nil, nil
	}
	var rps []*PostResolver
	for _, p := range refers {
		rps = append(rps, &PostResolver{Post: p})
	}
	return &rps, nil
}

// PostInput ...
type PostInput struct {
	ThreadID string
	Author   *string
	Content  string
	Refers   *[]string
}
