package model

import (
	"context"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/nanozuki/uexky/uuid64"
	"github.com/pkg/errors"
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
	Account    string        `bson:"account"` // not display in front
	CreateTime time.Time     `bson:"created_time"`

	MainTag string   `bson:"main_tag"`
	SubTags []string `bson:"sub_tags"`
	Title   string   `bson:"title"`
	Content string   `bson:"content"`
}

// ThreadInput ...
type ThreadInput struct {
	Author  *string
	Content string
	MainTag string
	SubTags *[]string
	Title   *string
}

func isMainTags(tag string) bool {
	for _, mt := range pkg.mainTags {
		if mt == tag {
			return true
		}
	}
	return false
}

// NewThread init new thread and insert to db
func NewThread(ctx context.Context, input *ThreadInput) (*Thread, error) {
	account, err := requireSignIn(ctx)
	if err != nil {
		return nil, err
	}
	if !isMainTags(input.MainTag) {
		return nil, errors.Errorf("Can't set main tag '%s'", input.MainTag)
	}
	subTags := []string{}
	if len(*input.SubTags) != 0 {
		subTags = *input.SubTags
	}
	for _, tag := range subTags {
		if isMainTags(tag) {
			return nil, errors.Errorf("Can't set main tag to sub tags '%s'", tag)
		}
	}
	if input.Content == "" {
		return nil, errors.New("Can't post an empty thread")
	}

	thread := &Thread{
		ObjectID:   bson.NewObjectId(),
		Account:    account.Token,
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

	if input.Author == nil || *input.Author == "" {
		thread.Anonymous = true
		author, err := account.AnonymousID(thread.ID, true)
		if err != nil {
			return nil, err
		}
		thread.Author = author
	} else {
		if !account.HaveName(*input.Author) {
			return nil, errors.Errorf("Can't find name '%s'", thread.Author)
		}
		thread.Author = *input.Author
	}

	if input.Title != nil && *input.Title != "" {
		thread.Title = *input.Title
	}

	c, cs := Colle("threads")
	defer cs()
	if err := c.Insert(thread); err != nil {
		return nil, err
	}
	return thread, nil
}

// GetThreadsByTags ...
func GetThreadsByTags(ctx context.Context, tags []string, sq *SliceQuery) (
	[]*Thread, *SliceInfo, error,
) {
	c, cs := Colle("threads")
	defer cs()
	find := bson.M{"tags": bson.M{"$in": tags}}
	if idQry := sq.QueryObject(); idQry != nil {
		find["id"] = idQry
	}

	var threads []*Thread
	if err := c.Find(find).Sort("-id").Limit(sq.Limit).All(threads); err != nil {
		return nil, nil, err
	}
	return threads, &SliceInfo{
		FirstCursor: threads[0].ID,
		LastCursor:  threads[len(threads)].ID,
	}, nil
}

// FindThread by id
func FindThread(ctx context.Context, ID string) (*Thread, error) {
	c, cs := Colle("threads")
	defer cs()
	var th Thread
	query := c.Find(bson.M{"id": ID})
	if count, err := query.Count(); err != nil {
		return nil, err
	} else if count == 0 {
		return nil, nil
	}
	if err := query.One(&th); err != nil {
		return nil, err
	}
	return &th, nil
}

func isThreadExist(threadID string) (bool, error) {
	c, cs := Colle("threads")
	defer cs()
	count, err := c.Find(bson.M{"id": threadID}).Count()
	if err != nil {
		return false, err
	}
	return count != 0, nil
}

// GetReplies ...
func (t *Thread) GetReplies(ctx context.Context, sq *SliceQuery) ([]*Post, *SliceInfo, error) {
	c, cs := Colle("posts")
	defer cs()

	var posts []*Post
	find := bson.M{"thread_id": t.ID}
	if idQry := sq.QueryObject(); idQry != nil {
		find["id"] = idQry
	}

	if err := c.Find(find).Sort("id").Limit(sq.Limit).All(posts); err != nil {
		return nil, nil, err
	}
	cnt := len(posts)
	si := &SliceInfo{FirstCursor: posts[0].ID, LastCursor: posts[cnt-1].ID}
	return posts, si, nil
}
