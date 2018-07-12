package model

import (
	"context"
	"time"

	"github.com/globalsign/mgo/bson"
	"gitlab.com/abyss.club/uexky/mw"
)

// Tag ...
type Tag struct {
	ObjectID    bson.ObjectId `bson:"_id"`
	Name        string        `bson:"name"`
	Parent      bool          `bson:"parent"`
	UpdatedTime time.Time     `bson:"updated_time"`
}

// TagTree ...
type TagTree struct {
	Nodes []struct {
		MainTag string   `json:"main_tag"`
		SubTags []string `json:"sub_tags"`
	} `json:"nodes"`
}

const tagTreeCacheKey = "mw:tagtree"

// UpsertTags ...
func UpsertTags(ctx context.Context, tags []*Tag) error {
	c := mw.GetMongo(ctx).C(colleTag)
	for _, tag := range tags {
		if _, err := c.Upsert(bson.M{"name": tag.Name}, bson.M{"$set": tag}); err != nil {
			return err
		}
	}
	return nil
}

// GetTagTree ...
func GetTagTree(ctx context.Context) (*TagTree, error) {
	var tree *TagTree
	if ok, err := mw.GetCache(ctx, tagTreeCacheKey, tree); err != nil {
		return nil, err
	} else if ok {
		return tree, nil
	}
	return nil, nil // TODO
}
