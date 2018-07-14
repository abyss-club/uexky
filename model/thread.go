package model

import (
	"context"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
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

// NewThread init new thread and insert to db
func NewThread(ctx context.Context, input *ThreadInput) (*Thread, error) {
	user, err := requireSignIn(ctx)
	if err != nil {
		return nil, err
	}
	if !isMainTag(input.MainTag) {
		return nil, errors.Errorf("Can't set main tag '%s'", input.MainTag)
	}
	subTags := []string{}
	if input.SubTags != nil && len(*input.SubTags) != 0 {
		subTags = *input.SubTags
	}
	for _, tag := range subTags {
		if isMainTag(tag) {
			return nil, errors.Errorf("Can't set main tag to sub tags '%s'", tag)
		}
	}
	if input.Content == "" {
		return nil, errors.New("Can't post an empty thread")
	}

	thread := &Thread{
		ObjectID:   bson.NewObjectId(),
		Anonymous:  input.Anonymous,
		UserID:     user.ID,
		CreateTime: time.Now(),

		MainTag: input.MainTag,
		SubTags: subTags,
		Content: input.Content,
	}

	threadID, err := postIDGenerator.New()
	if err != nil {
		return nil, err
	}
	thread.ID = threadID

	if input.Anonymous {
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

	if input.Title != nil && *input.Title != "" {
		thread.Title = *input.Title
	}

	c := mw.GetMongo(ctx).C(colleThread)
	if err := c.Insert(thread); err != nil {
		return nil, err
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

	queryObj, err := sq.GenQueryByObjectID()
	if err != nil {
		return nil, nil, err
	}
	if len(mainTags) != 0 {
		queryObj["main_tag"] = bson.M{"$in": mainTags}
	}
	if len(subTags) != 0 {
		queryObj["sub_tags"] = bson.M{"$in": subTags}
	}

	var threads []*Thread
	if err := sq.Find(ctx, colleThread, queryObj, &threads); err != nil {
		return nil, nil, err
	}
	if len(threads) == 0 {
		return threads, &SliceInfo{}, nil
	}
	return threads, &SliceInfo{
		FirstCursor: threads[0].ObjectID.Hex(),
		LastCursor:  threads[len(threads)-1].ObjectID.Hex(),
	}, nil
}

// FindThread by id
func FindThread(ctx context.Context, ID string) (*Thread, error) {
	c := mw.GetMongo(ctx).C(colleThread)
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
	c := mw.GetMongo(ctx).C(colleThread)
	count, err := c.Find(bson.M{"id": threadID}).Count()
	if err != nil {
		return false, err
	}
	return count != 0, nil
}

// GetReplies ...
func (t *Thread) GetReplies(ctx context.Context, sq *SliceQuery) ([]*Post, *SliceInfo, error) {
	queryObj, err := sq.GenQueryByObjectID()
	if err != nil {
		return nil, nil, err
	}
	queryObj["thread_id"] = t.ID

	var posts []*Post
	if err := sq.Find(ctx, collePost, queryObj, &posts); err != nil {
		return nil, nil, err
	}
	if len(posts) == 0 {
		return posts, &SliceInfo{}, nil
	}
	si := &SliceInfo{
		FirstCursor: posts[0].ObjectID.Hex(),
		LastCursor:  posts[len(posts)-1].ObjectID.Hex(),
	}
	return posts, si, nil
}
