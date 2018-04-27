package model

import (
	"context"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/nanozuki/uexky/model"
)

// Thread ...
type Thread struct {
	ObjectID   bson.ObjectId `bson:"_id"` // not display in front
	ID         string        `bson:"id"`
	Anonymous  bool          `bson:"anonymous"`
	Author     string        `bson:"author"`
	Account    bson.ObjectId `bson:"account"` // not display in front
	CreateTime time.Time     `bson:"created_time"`

	MainTag string   `bson:"main_tag"`
	SubTags []string `bson:"sub_tags"`
	Title   string   `bson:"title"`
	Content string   `bson:"content"`
}

// GetThreadsByTags ...
func GetThreadsByTags(ctx context.Context, tags []string, sq *SliceQuery) (
	[]*Threads, *SliceInfo, error,
) {
	c, cs := Colle("threads")
	defer cs()
	find := bson.M{"tags": bson.M{"$in": tags}}
	if idQry := sq.QueryObject(); idQry != nil {
		find["id"] = idQry
	}

	var threads []*model.Thread
	if err := c.Find(find).Sort("-id").Limit(sq.Limit).All(threads); err != nil {
		return nil, err
	}
	cnt := len(threads)
	return threads, &SliceInfo{FirstCursor: threads[0].ID, threads[cnt-1].ID}, nil
}

// GetReplies ...
func (t *Thread) GetReplies(ctx context.Context, sq *SliceQuery) ([]*Post, *SliceInfo, error) {
	c, cs := Colle("post")
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

// Post ...
type Post struct {
	ObjectID   bson.ObjectId `bson:"_id"`
	ID         string        `bson:"id"`
	Anonymous  string        `bson:"anonymous"`
	Author     string        `bson:"author"`
	Account    bson.ObjectId `bson:"account"`
	CreateTime time.Time     `bson:"creaate_time"`

	ThreadID bson.ObjectId `bson:"thread_id"`
	Content  string        `bson:"content"`
	Refers   []string      `bson:"refers"`
}
