package model

import (
	"time"

	"github.com/globalsign/mgo/bson"
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
