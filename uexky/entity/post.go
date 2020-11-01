package entity

import (
	"context"
	"fmt"
	"time"

	"gitlab.com/abyss.club/uexky/lib/errors"
	"gitlab.com/abyss.club/uexky/lib/uid"
)

type PostRepo interface {
	CheckIfDuplicated(ctx context.Context, userID uid.UID, content string) error
	GetByID(ctx context.Context, id uid.UID) (*Post, error)

	Insert(ctx context.Context, post *Post) (*Post, error)
	Update(ctx context.Context, post *Post) (*Post, error)

	QuotedPosts(ctx context.Context, post *Post) ([]*Post, error)
	QuotedCount(ctx context.Context, post *Post) (int, error)
}

type Post struct {
	ID        uid.UID   `json:"id"`
	ThreadID  uid.UID   `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
	Author    *Author   `json:"author"`
	QuoteIDs  []uid.UID `json:"-"`
	Content   string    `json:"content"`
	Blocked   bool      `json:"blocked"`
}

func (p Post) String() string {
	return fmt.Sprintf("<Post:%v:%s>", p.ID, p.ID.ToBase64String())
}

func NewPost(input *PostInput, user *User, thread *Thread, aid string) (*Post, error) {
	if thread.Locked {
		return nil, errors.BadParams.New("thread has been locked")
	}
	post := &Post{
		ID:        uid.NewUID(),
		ThreadID:  input.ThreadID,
		CreatedAt: time.Now(),
		Author: &Author{
			UserID:    user.ID,
			Guest:     user.Role == RoleGuest,
			Anonymous: input.Anonymous,
		},
		QuoteIDs: input.QuoteIds,
		Content:  input.Content,
	}
	if input.Anonymous {
		post.Author.Author = aid
	} else {
		if user.Name == nil {
			return nil, errors.BadParams.New("user name must be set")
		}
		post.Author.Author = *user.Name
	}
	return post, nil
}

func (p *Post) Block() {
	p.Blocked = true
}
