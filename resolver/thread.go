package resolver

import (
	"context"

	graphql "github.com/graph-gophers/graphql-go"
	"gitlab.com/abyss.club/uexky/model"
)

// ThreadSliceResolver ...
type ThreadSliceResolver struct {
	threads   []*ThreadResolver
	sliceInfo *SliceInfoResolver
}

// Threads ...
func (tsr *ThreadSliceResolver) Threads(ctx context.Context) ([]*ThreadResolver, error) {
	return tsr.threads, nil
}

// SliceInfo ...
func (tsr *ThreadSliceResolver) SliceInfo(ctx context.Context) (*SliceInfoResolver, error) {
	return tsr.sliceInfo, nil
}

// ThreadResolver ...
type ThreadResolver struct {
	Thread *model.Thread
}

// ID ...
func (tr *ThreadResolver) ID(ctx context.Context) (string, error) {
	return tr.Thread.ID, nil
}

// Anonymous ...
func (tr *ThreadResolver) Anonymous(ctx context.Context) (bool, error) {
	return tr.Thread.Anonymous, nil
}

// Author ...
func (tr *ThreadResolver) Author(ctx context.Context) (string, error) {
	return tr.Thread.Author, nil
}

// Content ...
func (tr *ThreadResolver) Content(ctx context.Context) (string, error) {
	return tr.Thread.Content, nil
}

// CreateTime ...
func (tr *ThreadResolver) CreateTime(ctx context.Context) (graphql.Time, error) {
	return graphql.Time{Time: tr.Thread.CreateTime}, nil
}

// MainTag ...
func (tr *ThreadResolver) MainTag(ctx context.Context) (string, error) {
	return tr.Thread.MainTag, nil
}

// SubTags ...
func (tr *ThreadResolver) SubTags(ctx context.Context) (*[]string, error) {
	if len(tr.Thread.SubTags) == 0 {
		return nil, nil
	}
	return &(tr.Thread.SubTags), nil
}

// Title ...
func (tr *ThreadResolver) Title(ctx context.Context) (*string, error) {
	if tr.Thread.Title == "" {
		return nil, nil
	}
	return &(tr.Thread.Title), nil
}

// Replies ...
func (tr *ThreadResolver) Replies(ctx context.Context, args struct {
	Query *SliceQuery
}) (
	*PostSliceResolver, error,
) {
	sq, err := args.Query.Parse(false)
	if err != nil {
		return nil, err
	}

	posts, sliceInfo, err := tr.Thread.GetReplies(ctx, sq)
	if err != nil {
		return nil, err
	}

	var prs []*PostResolver
	for _, p := range posts {
		prs = append(prs, &PostResolver{Post: p})
	}
	sir := &SliceInfoResolver{sliceInfo}
	return &PostSliceResolver{posts: prs, sliceInfo: sir}, nil
}

// CountOfReplies ...
func (tr *ThreadResolver) CountOfReplies(ctx context.Context) (int32, error) {
	count, err := tr.Thread.CountOfReplies(ctx)
	return int32(count), err
}
