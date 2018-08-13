package model

import (
	"context"
	"fmt"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky/mw"
)

const referLimit = 5

// Post ...
type Post struct {
	ObjectID   bson.ObjectId `bson:"_id"`
	ID         string        `bson:"id"`
	Anonymous  bool          `bson:"anonymous"`
	Author     string        `bson:"author"`
	UserID     bson.ObjectId `bson:"user_id"`
	CreateTime time.Time     `bson:"creaate_time"`

	ThreadID string   `bson:"thread_id"`
	Content  string   `bson:"content"`
	Refers   []string `bson:"refers,omitempty"`
}

// PostInput ...
type PostInput struct {
	ThreadID  string
	Anonymous bool
	Content   string
	Refers    *[]string
}

// NewPost ...
func NewPost(ctx context.Context, input *PostInput) (*Post, error) {
	user, err := requireSignIn(ctx)
	if err != nil {
		return nil, err
	}
	if exist, err := isThreadExist(ctx, input.ThreadID); err != nil {
		return nil, err
	} else if !exist {
		return nil, errors.Errorf("Thread %s is not exist", input.ThreadID)
	}
	if input.Content == "" {
		return nil, errors.New("required params missed")
	}
	post := &Post{
		ObjectID:   bson.NewObjectId(),
		Anonymous:  input.Anonymous,
		UserID:     user.ID,
		CreateTime: time.Now(),

		ThreadID: input.ThreadID,
		Content:  input.Content,
	}

	postID, err := postIDGenerator.New()
	if err != nil {
		return nil, err
	}
	post.ID = postID

	if input.Anonymous {
		author, err := user.AnonymousID(ctx, input.ThreadID, false)
		if err != nil {
			return nil, err
		}
		post.Author = author
	} else {
		if user.Name == "" {
			return nil, fmt.Errorf("Can't find name for user")
		}
		post.Author = user.Name
	}

	if input.Refers != nil {
		refers := *(input.Refers)
		if len(refers) > referLimit {
			return nil, fmt.Errorf("Count of Refers can't greater than 5")
		}
		for _, r := range refers {
			ok, err := isPostExist(ctx, r)
			if err != nil {
				return nil, err
			} else if !ok {
				return nil, fmt.Errorf("Can't find post '%s'", r)
			}
		}
		post.Refers = refers
	}

	m := mw.GetMongo(ctx)
	if err := m.C(collePost).Insert(post); err != nil {
		return nil, err
	}
	if err := m.C(colleThread).Update(
		bson.M{"thread_id": post.ThreadID},
		bson.M{"$set": bson.M{"update_time": post.CreateTime}},
	); err != nil {
		return nil, err
	}
	return post, nil
}

// FindPost ...
func FindPost(ctx context.Context, ID string) (*Post, error) {
	c := mw.GetMongo(ctx).C(collePost)
	query := c.Find(bson.M{"id": ID})
	if count, err := query.Count(); err != nil {
		return nil, err
	} else if count == 0 {
		return nil, nil
	}
	post := &Post{}
	if err := query.One(post); err != nil {
		return nil, err
	}
	return post, nil
}

// ReferPosts ...
func (p *Post) ReferPosts(ctx context.Context) ([]*Post, error) {
	var refers []*Post
	for _, id := range p.Refers {
		post, err := FindPost(ctx, id)
		if err != nil {
			return nil, err
		}
		refers = append(refers, post)
	}
	return refers, nil
}

func isPostExist(ctx context.Context, postID string) (bool, error) {
	c := mw.GetMongo(ctx).C(collePost)

	if cnt, err := c.Find(bson.M{"id": postID}).Count(); err != nil {
		return false, err
	} else if cnt == 0 {
		return false, nil
	}
	return true, nil
}
