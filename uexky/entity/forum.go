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
	Repo ForumRepo `wire:"-"` // TODO
}

func NewForumService(repo ForumRepo) ForumService {
	return ForumService{repo}
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

func (n *Thread) Replies(ctx context.Context, query SliceQuery) (*PostSlice, error) {
	panic(fmt.Errorf("not implemented"))
}

func (n *Thread) ReplyCount(ctx context.Context) (int, error) {
	panic(fmt.Errorf("not implemented"))
}

func (n *Thread) Catalog(ctx context.Context) ([]*ThreadCatalogItem, error) {
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

func (p *Post) Quotes(ctx context.Context) ([]*Post, error) {
	panic(fmt.Errorf("not implemented"))
}

func (p *Post) QuotedCount(ctx context.Context) (int, error) {
	panic(fmt.Errorf("not implemented"))
}

type Tag struct {
	Name      string   `json:"name"`
	IsMain    bool     `json:"isMain"`
	BelongsTo []string `json:"belongsTo"`
}
