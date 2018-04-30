package model

import (
	"context"
	"time"

	"github.com/globalsign/mgo/bson"
)

// Post ...
type Post struct {
	ObjectID   bson.ObjectId `bson:"_id"`
	ID         string        `bson:"id"`
	Anonymous  bool          `bson:"anonymous"`
	Author     string        `bson:"author"`
	Account    bson.ObjectId `bson:"account"`
	CreateTime time.Time     `bson:"creaate_time"`

	ThreadID bson.ObjectId `bson:"thread_id"`
	Content  string        `bson:"content"`
	Refers   []string      `bson:"refers"`
}

// ReferPosts ...
func (p *Post) ReferPosts(ctx context.Context) ([]*Post, error) {
	c, cs := Colle("posts")
	defer cs()

	var refers []*Post
	if err := c.Find(bson.M{"id": bson.M{"$in": p.Refers}}).All(&refers); err != nil {
		return nil, err
	}
	return refers, nil
}
