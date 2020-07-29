package entity

import (
	"context"
	"time"

	"gitlab.com/abyss.club/uexky/lib/uid"
)

// -- ForumRepo

type ThreadSearch struct {
	ID *uid.UID
}

type ThreadsSearch struct {
	UserID *uid.UID
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
	UserID   *uid.UID
	ThreadID *uid.UID
	DESC     bool
}

type PostUpdate struct {
	Blocked *bool
}

type TagSearch struct {
	Text  string
	Limit int
}

type UserTagUpdate struct {
	AddTags []string
	DelTags []string
}

const (
	BlockedContent       = "[此内容已被管理员屏蔽]"
	DuplicatedCheckRange = 3 * time.Minute
)

type ForumRepo interface {
	GetThread(ctx context.Context, search *ThreadSearch) (*Thread, error)
	GetThreadSlice(ctx context.Context, search *ThreadsSearch, query SliceQuery) (*ThreadSlice, error)
	GetThreadCatalog(ctx context.Context, id uid.UID) ([]*ThreadCatalogItem, error)
	GetAnonyID(ctx context.Context, userID uid.UID, threadID uid.UID) (string, error)
	InsertThread(ctx context.Context, thread *Thread) (*Thread, error)
	UpdateThread(ctx context.Context, thread *Thread) (*Thread, error)

	GetPost(ctx context.Context, search *PostSearch) (*Post, error)
	GetPosts(ctx context.Context, search *PostsSearch) ([]*Post, error)
	GetPostSlice(ctx context.Context, search *PostsSearch, query SliceQuery) (*PostSlice, error)
	GetPostCount(ctx context.Context, search *PostsSearch) (int, error)
	GetPostQuotedCount(ctx context.Context, id uid.UID) (int, error)
	InsertPost(ctx context.Context, post *Post) (*Post, error)
	UpdatePost(ctx context.Context, post *Post) (*Post, error)

	GetTags(ctx context.Context, search *TagSearch) ([]*Tag, error)
	GetMainTags(ctx context.Context) []string
	SetMainTags(ctx context.Context, tags []string) error

	CheckDuplicate(ctx context.Context, userID uid.UID, title, content string) error
}

type Author struct {
	UserID    uid.UID `json:"-"`
	Guest     bool    `json:"-"`
	Anonymous bool    `json:"anonymous"`
	Author    string  `json:"author"`
}

type PostData struct {
	ThreadID   uid.UID
	QuoteIDs   []uid.UID
	QuotePosts []*Post
}

type NewPostResponse struct {
	Post   *Post
	Thread *Thread
}
