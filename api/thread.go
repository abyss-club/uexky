package api

import (
	"context"
	"time"

	"github.com/nanozuki/uexky/model"
)

// ThreadSliceResolver ...
type ThreadSliceResolver struct {
	threads   []*ThreadResolver
	sliceInfo *SliceInfoResolver
}

// ThreadResolver ...
type ThreadResolver struct {
	Thread *model.Thread
}

// ThreadSlice ...
func (ts *Resolver) ThreadSlice(ctx context.Context, tags []string, limit int, after string) (
	*ThreadSliceResolver, error,
) {
	sq := &model.SliceQuery{Limit: limit, After: after}
	threads, sliceInfo, err = model.GetThreadsByTags(ctx, tags, sq)
	if err != nil {
		return nil, err
	}

	var trs []*ThreadResolver
	for _, t := range threads {
		trs = append(trs, &ThreadResolver{Thread: t})
	}
	return &ThreadSliceResolver{threads: trs, sliceInfo: sliceInfo}, nil
}

// Threads ...
func (tsr *ThreadSliceResolver) Threads(ctx context.Context) ([]*ThreadResolver, error) {
	return tsr.threads, nil
}

// SliceInfo ...
func (tsr *ThreadSliceResolver) SliceInfo(ctx context.Context) (*SliceInfoResolver, error) {
	return tsr.sliceInfo, nil
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
func (tr *ThreadResolver) CreateTime(ctx context.Context) (time.Time, error) {
	return tr.Thread.CreateTime, nil
}

// MainTag ...
func (tr *ThreadResolver) MainTag(ctx context.Context) (string, error) {
	return tr.Thread.MainTag, nil
}

// SubTags ...
func (tr *ThreadResolver) SubTags(ctx context.Context) ([]string, error) {
	return tr.Thread.SubTags, nil
}

// Title ...
func (tr *ThreadResolver) Title(ctx context.Context) (string, error) {
	return tr.Thread.Title, nil
}

// PostResolver ...
type PostResolver struct {
	Post *model.Post
}

// Replies ...
func (tr *ThreadResolver) Replies(ctx context.Context, limit int, after string) (
	*PostResolver, *SliceInfo, error,
) {
	sq := &model.SliceQuery{Limit: limit, After: after}
	posts, sliceInfo, err := tr.Thread.GetReplies(ctx, sq)
	if err != nil {
		return nil, nil, err
	}

	var prs []*PostResolver
	for _, p := range posts {
		prs = append(prs, &PostResolver{Post: p})
	}
	return prs, sliceInfo, nil
}
