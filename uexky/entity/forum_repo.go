package entity

import (
	"context"

	"gitlab.com/abyss.club/uexky/lib/uid"
)

type ThreadSearch struct {
	ID *uid.UID
}

type ThreadsSearch struct {
	UserID *int
	Tags   []string
}

type ThreadUpdate struct {
	MainTag *string
	SubTags []string
	Locked  *bool
	Blocked *bool
}

type PostSearch struct {
	ID *uid.UID
}
type PostsSearch struct {
	IDs      []uid.UID
	UserID   *int
	ThreadID *uid.UID
}

type PostUpdate struct {
	Blocked *bool
}

type TagSearch struct {
	Text     *string
	UserID   *int
	MainOnly bool
	Limit    int
}

type UserTagUpdate struct {
	AddTags []string
	DelTags []string
}

type ForumRepo interface {
	GetThread(ctx context.Context, search *ThreadSearch) (*Thread, error)
	GetThreadSlice(ctx context.Context, search *ThreadsSearch, query SliceQuery) (*ThreadSlice, error)
	GetThreadCatelog(ctx context.Context, id uid.UID) ([]*ThreadCatalogItem, error)
	GetThreadTags(ctx context.Context, id uid.UID) (main string, subs []string, err error)
	GetAnonyID(ctx context.Context, userID int, threadID uid.UID) (uid.UID, error)
	InsertThread(ctx context.Context, thread *Thread) error
	UpdateThread(ctx context.Context, id uid.UID, update *ThreadUpdate) error

	GetPost(ctx context.Context, search *PostSearch) (*Post, error)
	GetPosts(ctx context.Context, search *PostsSearch) ([]*Post, error)
	GetPostSlice(ctx context.Context, search *PostsSearch, query SliceQuery) (*PostSlice, error)
	GetPostCount(ctx context.Context, search *PostsSearch) (int, error)
	GetPostQuotesPosts(ctx context.Context, id uid.UID) ([]*Post, error)
	GetPostQuotedCount(ctx context.Context, id uid.UID) (int, error)
	InsertPost(ctx context.Context, post *Post) error
	UpdatePost(ctx context.Context, id uid.UID, update *PostUpdate) error

	GetTags(ctx context.Context, search *TagSearch) ([]*Tag, error)
	UpdateUserTags(ctx context.Context, userID int, update *UserTagUpdate) error
}
