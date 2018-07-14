package model

import (
	"context"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky/mgmt"
	"gitlab.com/abyss.club/uexky/mw"
)

// Tag ...
type Tag struct {
	ObjectID    bson.ObjectId `bson:"_id"`
	Name        string        `bson:"name"`
	Parent      string        `bson:"parent"`
	UpdatedTime time.Time     `bson:"updated_time"`
}

// TagTree ...
type TagTree struct {
	Nodes []*TagTreeNode `json:"nodes"`
}

// TagTreeNode ...
type TagTreeNode struct {
	MainTag string   `json:"main_tag"`
	SubTags []string `json:"sub_tags"`
}

const tagTreeCacheKey = "mw:tagtree"

func isMainTag(tag string) bool {
	for _, mt := range mgmt.Config.MainTags {
		if mt == tag {
			return true
		}
	}
	return false
}

// UpsertTags after insert thread...
func UpsertTags(ctx context.Context, mainTag string, tagStrings []string) error {
	c := mw.GetMongo(ctx).C(colleTag)
	c.EnsureIndex(mgo.Index{
		Key:        []string{"parent", "name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
	})
	for _, tag := range tagStrings {
		if !isMainTag(mainTag) {
			return errors.Errorf("'%s' is not main tag", mainTag)
		}
		if _, err := c.Upsert(bson.M{"name": tag}, bson.M{"$set": &Tag{
			Name: tag, Parent: mainTag, UpdatedTime: time.Now(),
		}}); err != nil {
			return err
		}
	}
	return nil
}

// GetTagTree ...
func GetTagTree(ctx context.Context) (*TagTree, error) {
	// try cache
	var tree *TagTree
	if ok, err := mw.GetCache(ctx, tagTreeCacheKey, tree); err != nil {
		return nil, err
	} else if ok {
		return tree, nil
	}

	tree = &TagTree{}
	for _, mTag := range mgmt.Config.MainTags {
		newest, err := getNewestSubTags(ctx, mTag)
		if err != nil {
			return nil, err
		}
		tree.Nodes = append(tree.Nodes, &TagTreeNode{mTag, newest})
	}

	// set cache
	if err := mw.SetCache(ctx, tagTreeCacheKey, tree, 600); err != nil {
		return nil, err
	}
	return tree, nil
}

func getNewestSubTags(ctx context.Context, mainTag string) ([]string, error) {
	c := mw.GetMongo(ctx).C(colleTag)
	c.EnsureIndexKey("parent")

	var tags []*Tag
	if err := c.Find(bson.M{"parent": mainTag}).Sort("-updated_time").Limit(10).All(tags); err != nil {
		return nil, errors.Wrapf(err, "find newest sub tag for '%s'", mainTag)
	}
	var tagStrings []string
	for _, tag := range tags {
		tagStrings = append(tagStrings, tag.Name)
	}
	return tagStrings, nil
}
