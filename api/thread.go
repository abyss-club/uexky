package api

import (
	"context"

	"github.com/globalsign/mgo/bson"
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
func (ts *Resolver) ThreadSlice(ctx context.Context, tags []string, first int, after string) (
	*ThreadSliceResolver, error,
) {
	session := model.MongoSession.Copy()
	defer session.Close()
	query := session.DB("test").C("threads").Find(bson.M{
		"tags": bson.M{"$in": tags},
		"_id":  bson.M{"$gt": bson.ObjectIdHex(after)},
	}).Sort("-_id").Limit(first)
	threads := []*model.Thread{}
	if err := query.All(threads); err != nil {
		return nil, err
	}

	var trs []*ThreadResolver
	for _, t := range threads {
		trs = append(trs, &ThreadResolver{Thread: t})
	}
	return &ThreadSliceResolver{
		threads: trs,
		sliceInfo: &SliceInfoResolver{
			firstCursor: threads[0].ID,
			lastCursor:  threads[len(threads)-1].ID, // TODO: empty list
		},
	}, nil
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

// Anonyumous ...
func (tr *ThreadResolver) Anonyumous(ctx context.Context) (bool, error) {
	return tr.Thread.Anonyumous, nil
}

// ID ...
func (tr *ThreadResolver) ID(ctx context.Context) (string, error) {
	return tr.Thread.ID, nil
}

// ID ...
func (tr *ThreadResolver) ID(ctx context.Context) (string, error) {
	return tr.Thread.ID, nil
}

// ID ...
func (tr *ThreadResolver) ID(ctx context.Context) (string, error) {
	return tr.Thread.ID, nil
}

// ID ...
func (tr *ThreadResolver) ID(ctx context.Context) (string, error) {
	return tr.Thread.ID, nil
}

// ID ...
func (tr *ThreadResolver) ID(ctx context.Context) (string, error) {
	return tr.Thread.ID, nil
}

// ID ...
func (tr *ThreadResolver) ID(ctx context.Context) (string, error) {
	return tr.Thread.ID, nil
}
