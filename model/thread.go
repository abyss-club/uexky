package model

import (
	"context"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky/mgmt"
	"gitlab.com/abyss.club/uexky/mw"
	"gitlab.com/abyss.club/uexky/uuid64"
)

var postIDGenerator = uuid64.Generator{Sections: []uuid64.Section{
	&uuid64.TimestampSection{Length: 6, Unit: time.Second, NoPadding: true},
	&uuid64.CounterSection{Length: 2, Unit: time.Second},
	&uuid64.RandomSection{Length: 1},
}}

// Thread ...
type Thread struct {
	ObjectID   bson.ObjectId `bson:"_id"` // not display in front
	ID         string        `bson:"id"`
	Anonymous  bool          `bson:"anonymous"`
	Author     string        `bson:"author"`
	UserID     bson.ObjectId `bson:"user_id"` // not display in front
	CreateTime time.Time     `bson:"created_time"`
	UpdateTime time.Time     `bson:"update_time"`

	MainTag string   `bson:"main_tag"`
	SubTags []string `bson:"sub_tags"`
	Title   string   `bson:"title"`
	Content string   `bson:"content"`
}

// ThreadInput ...
type ThreadInput struct {
	Anonymous bool
	Content   string
	MainTag   string
	SubTags   *[]string
	Title     *string
}

// ParseThead ...
func (ti *ThreadInput) ParseThead(ctx context.Context, user *User) (*Thread, error) {
	if !isMainTag(ti.MainTag) {
		return nil, errors.Errorf("Can't set main tag '%s'", ti.MainTag)
	}
	subTags := []string{}
	if ti.SubTags != nil && len(*ti.SubTags) != 0 {
		subTags = *ti.SubTags
	}
	for _, tag := range subTags {
		if isMainTag(tag) {
			return nil, errors.Errorf("Can't set main tag to sub tags '%s'", tag)
		}
	}
	if ti.Content == "" {
		return nil, errors.New("Can't post an empty thread")
	}

	now := time.Now()
	thread := &Thread{
		ObjectID:   bson.NewObjectId(),
		Anonymous:  ti.Anonymous,
		UserID:     user.ID,
		CreateTime: now,
		UpdateTime: now,

		MainTag: ti.MainTag,
		SubTags: subTags,
		Content: ti.Content,
	}

	threadID, err := postIDGenerator.New()
	if err != nil {
		return nil, err
	}
	thread.ID = threadID

	if ti.Anonymous {
		author, err := user.AnonymousID(ctx, thread.ID, true)
		if err != nil {
			return nil, err
		}
		thread.Author = author
	} else {
		if user.Name == "" {
			return nil, errors.Errorf("Can't find name for user")
		}
		thread.Author = user.Name
	}

	if ti.Title != nil && *ti.Title != "" {
		thread.Title = *ti.Title
	}
	return thread, nil
}

// NewThread init new thread and insert to db
func NewThread(ctx context.Context, input *ThreadInput) (*Thread, error) {
	if err := mw.FlowCostMut(ctx, mgmt.Config.RateLimit.Cost.PubThread); err != nil {
		return nil, err
	}
	user, err := requireSignIn(ctx)
	if err != nil {
		return nil, err
	}

	thread, err := input.ParseThead(ctx, user)
	if err != nil {
		return nil, err
	}

	c := mw.GetMongo(ctx).C(colleThread)
	if err := c.Insert(thread); err != nil {
		return nil, err
	}

	// Set Tag info
	if err := UpsertTags(ctx, thread.MainTag, thread.SubTags); err != nil {
		return nil, errors.Wrap(err, "set tag info")
	}
	return thread, nil
}

// GetThreadsByTags ...
func GetThreadsByTags(ctx context.Context, tags []string, sq *SliceQuery) (
	[]*Thread, *SliceInfo, error,
) {
	mainTags := []string{}
	subTags := []string{}
	for _, tag := range tags {
		if isMainTag(tag) {
			mainTags = append(mainTags, tag)
		} else {
			subTags = append(subTags, tag)
		}
	}

	queryObj, err := sq.GenQueryByTime("update_time")
	if err != nil {
		return nil, nil, err
	}
	if len(mainTags) != 0 && len(subTags) != 0 {
		queryObj["$or"] = []bson.M{
			bson.M{"main_tag": bson.M{"$in": mainTags}},
			bson.M{"sub_tags": bson.M{"$in": subTags}},
		}
	} else if len(mainTags) != 0 {
		queryObj["main_tag"] = bson.M{"$in": mainTags}
	} else if len(subTags) != 0 {
		queryObj["sub_tags"] = bson.M{"$in": subTags}
	}

	c := mw.GetMongo(ctx).C(colleThread)
	c.EnsureIndexKey("main_tag")
	c.EnsureIndexKey("sub_tags")
	c.EnsureIndexKey("update_time")

	var threads []*Thread
	if err := sq.Find(ctx, colleThread, "update_time", queryObj, &threads); err != nil {
		return nil, nil, err
	}
	if len(threads) == 0 {
		return threads, &SliceInfo{}, nil
	}
	if !sq.Desc {
		ReverseSlice(threads)
	}
	return threads, &SliceInfo{
		FirstCursor: threads[0].genCursor(),
		LastCursor:  threads[len(threads)-1].genCursor(),
	}, nil
}

// FindThread by id
func FindThread(ctx context.Context, ID string) (*Thread, error) {
	if err := mw.FlowCostQuery(ctx, 1); err != nil {
		return nil, err
	}
	c := mw.GetMongo(ctx).C(colleThread)
	c.EnsureIndexKey("id")

	var th Thread
	query := c.Find(bson.M{"id": ID})
	if count, err := query.Count(); err != nil {
		return nil, err
	} else if count == 0 {
		return nil, errors.Errorf("Can't Find Thread '%v'", ID)
	}
	if err := query.One(&th); err != nil {
		return nil, err
	}
	return &th, nil
}

func isThreadExist(ctx context.Context, threadID string) (bool, error) {
	if err := mw.FlowCostQuery(ctx, 1); err != nil {
		return false, err
	}
	c := mw.GetMongo(ctx).C(colleThread)
	c.EnsureIndexKey("id")

	count, err := c.Find(bson.M{"id": threadID}).Count()
	if err != nil {
		return false, err
	}
	return count != 0, nil
}

// GetReplies ...
func (t *Thread) GetReplies(ctx context.Context, sq *SliceQuery) ([]*Post, *SliceInfo, error) {
	c := mw.GetMongo(ctx).C(collePost)
	c.EnsureIndexKey("thread_id")

	queryObj, err := sq.GenQueryByObjectID()
	if err != nil {
		return nil, nil, err
	}
	queryObj["thread_id"] = t.ID

	var posts []*Post
	if err := sq.Find(ctx, collePost, "_id", queryObj, &posts); err != nil {
		return nil, nil, err
	}
	if len(posts) == 0 {
		return posts, &SliceInfo{}, nil
	}
	if sq.Desc {
		ReverseSlice(posts)
	}
	si := &SliceInfo{
		FirstCursor: posts[0].ObjectID.Hex(),
		LastCursor:  posts[len(posts)-1].ObjectID.Hex(),
	}
	return posts, si, nil
}

// ReplyCount ...
func (t *Thread) ReplyCount(ctx context.Context) (int, error) {
	if err := mw.FlowCostQuery(ctx, 1); err != nil {
		return 0, err
	}
	c := mw.GetMongo(ctx).C(collePost)
	c.EnsureIndexKey("thread_id")

	return c.Find(bson.M{"thread_id": t.ID}).Count()
}

// return unix time of update time in millisecond(ms)
func (t *Thread) genCursor() string {
	return genTimeCursor(t.UpdateTime)
}
