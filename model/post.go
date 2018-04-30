package model

import (
	"context"
	"fmt"
	"time"

	"github.com/globalsign/mgo/bson"
)

const referLimit = 5

// Post ...
type Post struct {
	ObjectID   bson.ObjectId `bson:"_id"`
	ID         string        `bson:"id"`
	Anonymous  bool          `bson:"anonymous"`
	Author     string        `bson:"author"`
	Account    string        `bson:"account"`
	CreateTime time.Time     `bson:"creaate_time"`

	ThreadID string   `bson:"thread_id"`
	Content  string   `bson:"content"`
	Refers   []string `bson:"refers"`
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

func isPostExist(postID string) (bool, error) {
	c, cs := Colle("posts")
	defer cs()

	if cnt, err := c.Find(bson.M{"id": postID}).Count(); err != nil {
		return false, err
	} else if cnt == 0 {
		return false, nil
	}
	return true, nil
}

// InsertPost ...
func (p *Post) InsertPost(ctx context.Context) error {
	account, err := requireSignIn(ctx)
	if err != nil {
		return err
	}

	p.ObjectID = bson.NewObjectId()
	if p.ID, err = postIDGenerator.New(); err != nil {
		return err
	}
	if p.Author == "" {
		p.Anonymous = true
		if p.Author, err = account.AnonymousID(p.ThreadID, false); err != nil {
			return err
		}
	} else {
		if !account.HaveName(p.Author) {
			return fmt.Errorf("Can't find name '%s'", p.Author)
		}
	}
	p.Author = account.Token
	p.CreateTime = time.Now()
	if len(p.Refers) > referLimit {
		return fmt.Errorf("Count of Refers can't greater than 5")
	}
	for _, r := range p.Refers {
		ok, err := isPostExist(p.ID)
		if err != nil {
			return err
		} else if !ok {
			return fmt.Errorf("Can't find post '%s'", r)
		}
	}

	c, cs := Colle("posts")
	defer cs()
	return c.Insert(p)
}
