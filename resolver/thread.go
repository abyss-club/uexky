package resolver

import (
	"context"

	graphql "github.com/graph-gophers/graphql-go"
	"gitlab.com/abyss.club/uexky/model"
)

// queries:

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
func (r *Resolver) Thread(
	ctx context.Context, args struct{ ID string },
) (*ThreadResolver, error) {
	th, err := model.FindThread(ctx, args.ID)
	if err != nil {
		return nil, err
	}
	if th == nil {
		return nil, nil
	}
	return &ThreadResolver{Thread: th}, nil
}

// mutations:

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

// types:

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

// ReplyCount ...
func (tr *ThreadResolver) ReplyCount(ctx context.Context) (int32, error) {
	count, err := tr.Thread.ReplyCount(ctx)
	return int32(count), err
}
