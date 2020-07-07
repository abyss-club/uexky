package entity

import (
	"context"

	"gitlab.com/abyss.club/uexky/lib/uid"
)

// -- ForumRepo

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

const BlockedContent = "[此内容已被管理员屏蔽]"

type ForumRepo interface {
	GetThread(ctx context.Context, search *ThreadSearch) (*Thread, error)
	GetThreadSlice(ctx context.Context, search *ThreadsSearch, query SliceQuery) (*ThreadSlice, error)
	GetThreadCatalog(ctx context.Context, id uid.UID) ([]*ThreadCatalogItem, error)
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
	GetMainTags(ctx context.Context) ([]string, error)
	SetMainTags(ctx context.Context, tags []string) error
}

// -- Author

type Author struct {
	UserID      int
	AnonymousID *uid.UID
	UserName    *string
}

func (a Author) Name(anonymous bool) string {
	if !anonymous {
		return *a.UserName
	}
	return a.AnonymousID.ToBase64String()
}

// -- Entity Extension

type PostData struct {
	ThreadID   uid.UID
	Author     Author
	QuoteIDs   []uid.UID
	QuotePosts []*Post
}

type NewPostResponse struct {
	Post   *Post
	Thread *Thread
}
