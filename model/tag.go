package model

import (
	"context"
	"fmt"
	"log"
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
	MainTags    []string      `bson:"main_tags"`
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
		Key:        []string{"name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
	})
	for _, tag := range tagStrings {
		log.Printf("insert tag '%s'", tag)
		if _, err := c.Upsert(bson.M{"name": tag}, bson.M{
			"$set": bson.M{
				"name":         tag,
				"updated_time": time.Now(),
			},
			"$addToSet": bson.M{
				"main_tags": mainTag,
			},
		}); err != nil {
			return err
		}
	}
	return nil
}

// GetTagTree ...
func GetTagTree(ctx context.Context, query string) (*TagTree, error) {
	// try cache
	key := fmt.Sprintf("%s:%s", tagTreeCacheKey, query)
	tree := &TagTree{}
	if ok, err := mw.GetCache(ctx, key, tree); err != nil {
		return nil, err
	} else if ok {
		return tree, nil
	}

	tree = &TagTree{}
	for _, mTag := range mgmt.Config.MainTags {
		log.Printf("start fetch subTags for '%s'", mTag)
		newest, err := getNewestSubTags(ctx, mTag, query)
		if err != nil {
			return nil, err
		}
		tree.Nodes = append(tree.Nodes, &TagTreeNode{mTag, newest})
	}

	// set cache
	if err := mw.SetCache(ctx, key, tree, 600); err != nil {
		return nil, err
	}
	return tree, nil
}

func getNewestSubTags(ctx context.Context, mainTag string, query string) ([]string, error) {
	c := mw.GetMongo(ctx).C(colleTag)
	c.EnsureIndexKey("main_tags")

	var tags []*Tag
	find := bson.M{"main_tags": mainTag}
	if query != "" {
		find["name"] = bson.M{"$regex": query}
	}
	if err := c.Find(find).Sort("-updated_time", "-_id").Limit(10).All(&tags); err != nil {
		return nil, errors.Wrapf(err, "find newest sub tag for '%s'", mainTag)
	}
	var tagStrings []string
	for _, tag := range tags {
		tagStrings = append(tagStrings, tag.Name)
	}
	return tagStrings, nil
}
