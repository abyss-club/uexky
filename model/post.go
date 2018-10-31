package model

import (
	"context"
	"fmt"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky/mw"
)

const quoteLimit = 3

// Post ...
type Post struct {
	ObjectID   bson.ObjectId `bson:"_id"`
	ID         string        `bson:"id"`
	Anonymous  bool          `bson:"anonymous"`
	Author     string        `bson:"author"`
	UserID     bson.ObjectId `bson:"user_id"`
	CreateTime time.Time     `bson:"create_time"`

	ThreadID string   `bson:"thread_id"`
	Content  string   `bson:"content"`
	Quotes   []string `bson:"quotes,omitempty"`
}

// PostInput ...
type PostInput struct {
	ThreadID  string
	Anonymous bool
	Content   string
	Quotes    *[]string
}

// ParsePost ...
func (pi *PostInput) ParsePost(ctx context.Context, user *User) (
	*Post, *Thread, []*Post, error,
) {
	if pi.Content == "" {
		return nil, nil, nil, errors.New("required params missed")
	}
	thread, err := FindThread(ctx, pi.ThreadID)
	if err != nil {
		return nil, nil, nil, err
	}

	post := &Post{
		ObjectID:   bson.NewObjectId(),
		Anonymous:  pi.Anonymous,
		UserID:     user.ID,
		CreateTime: time.Now(),

		ThreadID: pi.ThreadID,
		Content:  pi.Content,
	}

	postID, err := postIDGenerator.New()
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "gen post id")
	}
	post.ID = postID

	if pi.Anonymous {
		author, err := user.AnonymousID(ctx, pi.ThreadID, false)
		if err != nil {
			return nil, nil, nil, errors.Wrap(err, "get AnonymousID")
		}
		post.Author = author
	} else {
		if user.Name == "" {
			return nil, nil, nil, fmt.Errorf("Can't find name for user")
		}
		post.Author = user.Name
	}

	quotePosts := []*Post{}
	if pi.Quotes != nil {
		quotes := *(pi.Quotes)
		if len(quotes) > quoteLimit {
			return nil, nil, nil, fmt.Errorf("Count of Quotes can't greater than 5")
		}
		for _, r := range quotes {
			p, err := FindPost(ctx, r)
			if err != nil {
				return nil, nil, nil, errors.Wrap(err, "find quote posts")
			}
			quotePosts = append(quotePosts, p)
		}
		post.Quotes = quotes
	}
	return post, thread, quotePosts, nil
}

// NewPost ...
func NewPost(ctx context.Context, input *PostInput) (*Post, error) {
	user, err := requireSignIn(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "find sign info")
	}

	post, thread, quotes, err := input.ParsePost(ctx, user)
	if err != nil {
		return nil, err
	}

	m := mw.GetMongo(ctx)
	if err := m.C(collePost).Insert(post); err != nil {
		return nil, errors.Wrap(err, "insert post")
	}
	if err := m.C(colleThread).Update(
		bson.M{"id": post.ThreadID},
		bson.M{"$set": bson.M{"update_time": post.CreateTime}},
	); err != nil {
		return nil, errors.Wrapf(err, "update thread %s", post.ThreadID)
	}

	if err := TriggerNotifForPost(ctx, thread, post, quotes); err != nil {
		return nil, err
	}
	return post, nil
}

// FindPost ...
func FindPost(ctx context.Context, ID string) (*Post, error) {
	c := mw.GetMongo(ctx).C(collePost)
	c.EnsureIndexKey("id")
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

// QuotePosts ...
func (p *Post) QuotePosts(ctx context.Context) ([]*Post, error) {
	var quotes []*Post
	for _, id := range p.Quotes {
		post, err := FindPost(ctx, id)
		if err != nil {
			return nil, err
		}
		quotes = append(quotes, post)
	}
	return quotes, nil
}

// QuoteCount ...
func (p *Post) QuoteCount(ctx context.Context) (int, error) {
	c := mw.GetMongo(ctx).C(collePost)
	c.EnsureIndexKey("quotes")
	return c.Find(bson.M{"quotes": p.ID}).Count()
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
