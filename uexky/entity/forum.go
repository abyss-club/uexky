// forum aggragate: thread, post, tags

package entity

import (
	"context"
	"errors"
	"time"

	"gitlab.com/abyss.club/uexky/lib/uid"
)

type ForumService struct {
	Repo ForumRepo
}

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

type Thread struct {
	ID        uid.UID   `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Anonymous bool      `json:"anonymous"`
	Title     *string   `json:"title"`
	Content   string    `json:"content"`
	MainTag   string    `json:"main_tag"`
	SubTags   []string  `json:"sub_tags"`
	Blocked   bool      `json:"blocked"`
	Locked    bool      `json:"locked"`

	Repo      ForumRepo `json:"-"`
	AuthorObj Author    `json:"-"`
}

func (f *ForumService) NewThread(ctx context.Context, user *User, input ThreadInput) (*Thread, error) {
	thread := &Thread{
		ID:        uid.NewUID(),
		CreatedAt: time.Now(),
		Anonymous: input.Anonymous,
		Title:     input.Title,
		Content:   input.Content,
		MainTag:   input.MainTag,
		SubTags:   input.SubTags,

		Repo: f.Repo,
		AuthorObj: Author{
			UserID: user.ID,
		},
	}
	if input.Anonymous {
		aid := uid.NewUID()
		thread.AuthorObj.AnonymousID = &aid
	} else {
		if user.Name == nil {
			return nil, errors.New("user name must be set")
		}
		thread.AuthorObj.UserName = user.Name
	}
	err := f.Repo.InsertThread(ctx, thread)
	if err != nil {
		return nil, err
	}
	return thread, err
}

func (f *ForumService) GetThreadByID(ctx context.Context, threadID uid.UID) (*Thread, error) {
	return f.Repo.GetThread(ctx, &ThreadSearch{ID: &threadID})
}

func (f *ForumService) GetUserThreads(ctx context.Context, user *User, query SliceQuery) (*ThreadSlice, error) {
	return f.Repo.GetThreadSlice(ctx, &ThreadsSearch{UserID: &user.ID}, query)
}

func (f *ForumService) SearchThreads(
	ctx context.Context, tags []string, query SliceQuery,
) (*ThreadSlice, error) {
	return f.Repo.GetThreadSlice(ctx, &ThreadsSearch{Tags: tags}, query)
}

func (n *Thread) Author() string {
	return n.AuthorObj.Name(n.Anonymous)
}

func (n *Thread) Replies(ctx context.Context, query SliceQuery) (*PostSlice, error) {
	return n.Repo.GetPostSlice(ctx, &PostsSearch{ThreadID: &n.ID}, query)
}

func (n *Thread) ReplyCount(ctx context.Context) (int, error) {
	return n.Repo.GetPostCount(ctx, &PostsSearch{ThreadID: &n.ID})
}

func (n *Thread) Catalog(ctx context.Context) ([]*ThreadCatalogItem, error) {
	return n.Repo.GetThreadCatelog(ctx, n.ID)
}

func (n *Thread) EditTags(ctx context.Context, mainTag string, subTags []string) error {
	update := &ThreadUpdate{MainTag: &mainTag, SubTags: subTags}
	err := n.Repo.UpdateThread(ctx, n.ID, update)
	if err != nil {
		return err
	}
	n.MainTag = mainTag
	n.SubTags = subTags
	return nil
}

func (n *Thread) Lock(ctx context.Context) error {
	locked := true
	if err := n.Repo.UpdateThread(ctx, n.ID, &ThreadUpdate{Locked: &locked}); err != nil {
		return err
	}
	n.Locked = true
	return nil
}

func (n *Thread) Block(ctx context.Context) error {
	blocked := true
	if err := n.Repo.UpdateThread(ctx, n.ID, &ThreadUpdate{Blocked: &blocked}); err != nil {
		return err
	}
	n.Blocked = true
	return nil
}

type PostData struct {
	ThreadID    uid.UID
	Author      Author
	QuotedIDs   []uid.UID
	QuotedPosts []*Post
}

type Post struct {
	ID        uid.UID   `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Anonymous bool      `json:"anonymous"`
	Content   string    `json:"content"`
	Blocked   bool      `json:"blocked"`

	Repo ForumRepo `json:"-"`
	Data PostData  `json:"-"`
}

type NewPostResponse struct {
	Post       *Post
	Thread     *Thread
	QuotedPost []*Post
}

func (f *ForumService) NewPost(ctx context.Context, user *User, input PostInput) (*NewPostResponse, error) {
	post := &Post{
		ID:        uid.NewUID(),
		CreatedAt: time.Now(),
		Anonymous: input.Anonymous,
		Content:   input.Content,

		Repo: f.Repo,
		Data: PostData{
			Author:   Author{UserID: user.ID},
			ThreadID: input.ThreadID,
		},
	}
	if input.Anonymous {
		aid := uid.NewUID()
		post.Data.Author.AnonymousID = &aid // TODO: find aid
	} else {
		if user.Name == nil {
			return nil, errors.New("user name must be set")
		}
		post.Data.Author.UserName = user.Name
	}
	err := f.Repo.InsertPost(ctx, post)
	return &NewPostResponse{Post: post}, err
}

func (f *ForumService) GetPostByID(ctx context.Context, postID uid.UID) (*Post, error) {
	return f.Repo.GetPost(ctx, &PostSearch{ID: &postID})
}

func (f *ForumService) GetUserPosts(ctx context.Context, user *User, query SliceQuery) (*PostSlice, error) {
	return f.Repo.GetPostSlice(ctx, &PostsSearch{UserID: &user.ID}, query)
}

func (p *Post) Author() string {
	return p.Data.Author.Name(p.Anonymous)
}

func (p *Post) Quotes(ctx context.Context) ([]*Post, error) {
	if p.Data.QuotedPosts != nil {
		return p.Data.QuotedPosts, nil
	}
	quotedPosts, err := p.Repo.GetPosts(ctx, &PostsSearch{IDs: p.Data.QuotedIDs})
	if err != nil {
		return nil, err
	}
	p.Data.QuotedPosts = quotedPosts
	return p.Data.QuotedPosts, nil
}

func (p *Post) QuotedCount(ctx context.Context) (int, error) {
	return p.Repo.GetPostQuotedCount(ctx, p.ID)
}

func (p *Post) Block(ctx context.Context) error {
	blocked := true
	if err := p.Repo.UpdatePost(ctx, p.ID, &PostUpdate{Blocked: &blocked}); err != nil {
		return err
	}
	p.Blocked = true
	return nil
}

func (f *ForumService) GetMainTags(ctx context.Context) ([]string, error) {
	return f.Repo.GetMainTags(ctx)
}

func (f *ForumService) GetRecommendedTags(ctx context.Context) ([]string, error) {
	return f.GetMainTags(ctx)
}

func (f *ForumService) SearchTags(ctx context.Context, query *string, limit *int) ([]*Tag, error) {
	search := &TagSearch{Text: query}
	if limit == nil {
		search.Limit = 10
	} else {
		search.Limit = *limit
	}
	return f.Repo.GetTags(ctx, search)
}
