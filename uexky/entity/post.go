package entity

import (
	"context"
	"fmt"
	"time"

	"gitlab.com/abyss.club/uexky/lib/uerr"
	"gitlab.com/abyss.club/uexky/lib/uid"
)

type PostRepo interface {
	CheckIfDuplicate(ctx context.Context, userID uid.UID, content string) (bool, error)
	GetByID(ctx context.Context, id uid.UID) (*Post, error)

	Insert(ctx context.Context, post *Post) (*Post, error)
	Update(ctx context.Context, post *Post) (*Post, error)

	QuotedPosts(ctx context.Context, post *Post) ([]*Post, error)
	QuotedCount(ctx context.Context, post *Post) (int, error)
}

type Post struct {
	ID        uid.UID   `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Author    *Author   `json:"author"`
	Content   string    `json:"content"`
	Blocked   bool      `json:"blocked"`

	Data PostData `json:"-"`
}

type PostData struct {
	ThreadID   uid.UID
	QuoteIDs   []uid.UID
	QuotePosts []*Post
}

func (p Post) String() string {
	return fmt.Sprintf("<Post:%v:%s>", p.ID, p.ID.ToBase64String())
}

func NewPost(input *PostInput, user *User, thread *Thread, aid uid.UID) (*Post, error) {
	if thread.Locked {
		return nil, uerr.New(uerr.ParamsError, "thread has been locked")
	}
	post := &Post{
		ID:        uid.NewUID(),
		CreatedAt: time.Now(),
		Author: &Author{
			UserID:    user.ID,
			Guest:     user.Role == RoleGuest,
			Anonymous: input.Anonymous,
		},
		Content: input.Content,

		Data: PostData{
			ThreadID:   input.ThreadID,
			QuoteIDs:   input.QuoteIds,
			QuotePosts: make([]*Post, 0),
		},
	}
	if input.Anonymous {
		post.Author.Author = aid.ToBase64String()
		// if user.ID == thread.Author.UserID && thread.Author.Anonymous {
		// 	post.Author.Author = thread.Author.Author
		// } else {
		// 	aid, err := f.Repo.GetAnonyID(ctx, user.ID, thread.ID)
		// 	if err != nil {
		// 		return nil, errors.Wrapf(err, "NewPost(user=%+v, input=%+v)", user, input)
		// 	}
		// 	post.Author.Author = aid
		// }
	} else {
		if user.Name == nil {
			return nil, uerr.New(uerr.ParamsError, "user name must be set")
		}
		post.Author.Author = *user.Name
	}
	return post, nil
}

func (p *Post) Block() {
	p.Blocked = true
}
