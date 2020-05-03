// forum aggragate: thread, post, tags

package entity

import (
	"context"
	"fmt"
	"time"

	"gitlab.com/abyss.club/uexky/lib/uid"
)

type ForumRepo interface{}

type ForumService struct {
	Repo ForumRepo
}

type Thread struct {
	ID        uid.UID   `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Anonymous bool      `json:"anonymous"`
	Author    string    `json:"author"`
	Title     *string   `json:"title"`
	Content   string    `json:"content"`
	MainTag   string    `json:"mainTag"`
	SubTags   []string  `json:"subTags"`
	Blocked   bool      `json:"blocked"`
	Locked    bool      `json:"locked"`
}

func (f *ForumService) NewThread(ctx context.Context, input ThreadInput) (*Thread, error) {
	panic(fmt.Errorf("not implemented"))
}

func (f *ForumService) GetThreadByID(ctx context.Context, threadID uid.UID) (*Thread, error) {
	panic(fmt.Errorf("not implemented"))
}

func (f *ForumService) GetUserThreads(ctx context.Context, user *User, query SliceQuery) (*ThreadSlice, error) {
	panic(fmt.Errorf("not implemented"))
}

func (f *ForumService) SearchThreads(
	ctx context.Context, tags []string, query SliceQuery,
) (*ThreadSlice, error) {
	panic(fmt.Errorf("not implemented"))
}

func (n *Thread) Replies(ctx context.Context, query SliceQuery) (*PostSlice, error) {
	panic(fmt.Errorf("not implemented"))
}

func (n *Thread) ReplyCount(ctx context.Context) (int, error) {
	panic(fmt.Errorf("not implemented"))
}

func (n *Thread) Catalog(ctx context.Context) ([]*ThreadCatalogItem, error) {
	panic(fmt.Errorf("not implemented"))
}

func (n *Thread) EditTags(ctx context.Context, mainTag string, subTags []string) error {
	panic(fmt.Errorf("not implemented"))
}

func (n *Thread) Lock(ctx context.Context) error {
	panic(fmt.Errorf("not implemented"))
}

func (n *Thread) Block(ctx context.Context) error {
	panic(fmt.Errorf("not implemented"))
}

type Post struct {
	ID        uid.UID   `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Anonymous bool      `json:"anonymous"`
	Author    string    `json:"author"`
	Content   string    `json:"content"`
	Blocked   bool      `json:"blocked"`
}

type NewPostResponse struct {
	Post       *Post
	Thread     *Thread
	QuotedPost []*Post
}

func (f *ForumService) NewPost(ctx context.Context, input PostInput) (*NewPostResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (f *ForumService) GetPostByID(ctx context.Context, postID uid.UID) (*Post, error) {
	panic(fmt.Errorf("not implemented"))
}

func (f *ForumService) GetUserPosts(ctx context.Context, user *User, query SliceQuery) (*PostSlice, error) {
	panic(fmt.Errorf("not implemented"))
}

func (p *Post) Quotes(ctx context.Context) ([]*Post, error) {
	panic(fmt.Errorf("not implemented"))
}

func (p *Post) QuotedCount(ctx context.Context) (int, error) {
	panic(fmt.Errorf("not implemented"))
}

func (p *Post) Block(ctx context.Context) error {
	panic(fmt.Errorf("not implemented"))
}

type Tag struct {
	Name      string   `json:"name"`
	IsMain    bool     `json:"isMain"`
	BelongsTo []string `json:"belongsTo"`
}

func (f *ForumService) GetMainTags(ctx context.Context) ([]string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (f *ForumService) GetRecommendedTags(ctx context.Context) ([]string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (f *ForumService) SearchTags(ctx context.Context, query *string, limit *int) ([]*Tag, error) {
	panic(fmt.Errorf("not implemented"))
}

func (f *ForumService) GetUserTags(ctx context.Context, user *User) ([]string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (f *ForumService) SyncUserTags(ctx context.Context, user *User, tags []*string) error {
	panic(fmt.Errorf("not implemented"))
}

func (f *ForumService) AddUserSubbedTag(ctx context.Context, user *User, tag string) error {
	panic(fmt.Errorf("not implemented"))
}

func (f *ForumService) DelUserSubbedTag(ctx context.Context, user *User, tag string) error {
	panic(fmt.Errorf("not implemented"))
}
