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

type Thread struct {
	ID        uid.UID   `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Anonymous bool      `json:"anonymous"`
	Author    string    `json:"author"`
	Title     *string   `json:"title"`
	Content   string    `json:"content"`
	Blocked   bool      `json:"blocked"`
	Locked    bool      `json:"locked"`
	Repo      ForumRepo `json:"-"`

	mainTag string
	subTags []string
}

func (f *ForumService) NewThread(ctx context.Context, user *User, input ThreadInput) (*Thread, error) {
	thread := &Thread{
		ID:        uid.NewUID(),
		CreatedAt: time.Now(),
		Anonymous: input.Anonymous,
		Title:     input.Title,
		Content:   input.Content,
		Repo:      f.Repo,

		mainTag: input.MainTag,
		subTags: input.SubTags,
	}
	if input.Anonymous {
		thread.Author = uid.NewUID().ToBase64String()
	} else {
		if user.Name == nil {
			return nil, errors.New("user name must be set")
		}
		thread.Author = *user.Name
	}
	err := f.Repo.InsertThread(ctx, user.ID, thread)
	if err != nil {
		return nil, err
	}
	err = thread.EditTags(ctx, thread.mainTag, thread.subTags)
	return thread, err
}

func (f *ForumService) GetThreadByID(ctx context.Context, threadID uid.UID) (*Thread, error) {
	return f.Repo.GetThread(ctx, &ThreadSearch{ID: &threadID})
}

func (f *ForumService) GetUserThreads(ctx context.Context, user *User, query SliceQuery) (*ThreadSlice, error) {
	return f.Repo.GetThreadSlice(ctx, &ThreadSearch{UserID: &user.ID}, query)
}

func (f *ForumService) SearchThreads(
	ctx context.Context, tags []string, query SliceQuery,
) (*ThreadSlice, error) {
	return f.Repo.GetThreadSlice(ctx, &ThreadSearch{Tags: tags}, query)
}

func (n *Thread) Replies(ctx context.Context, query SliceQuery) (*PostSlice, error) {
	return n.Repo.GetPostSlice(ctx, &PostSearch{ThreadID: &n.ID}, query)
}

func (n *Thread) ReplyCount(ctx context.Context) (int, error) {
	return n.Repo.GetPostCount(ctx, &PostSearch{ThreadID: &n.ID})
}

func (n *Thread) Catalog(ctx context.Context) ([]*ThreadCatalogItem, error) {
	return n.Repo.GetThreadCatelog(ctx, n.ID)
}

func (n *Thread) getTags(ctx context.Context) error {
	main, subs, err := n.Repo.GetThreadTags(ctx, n.ID)
	if err != nil {
		return err
	}
	n.mainTag = main
	n.subTags = subs
	return nil
}

func (n *Thread) MainTag(ctx context.Context) (string, error) {
	if n.mainTag != "" {
		return n.mainTag, nil
	}
	err := n.getTags(ctx)
	if err != nil {
		return "", err
	}
	return n.mainTag, nil
}

func (n *Thread) SubTags(ctx context.Context) ([]string, error) {
	if n.subTags != nil {
		return n.subTags, nil
	}
	err := n.getTags(ctx)
	if err != nil {
		return nil, err
	}
	return n.subTags, nil
}

func (n *Thread) EditTags(ctx context.Context, mainTag string, subTags []string) error {
	update := &ThreadUpdate{SubTags: subTags}
	if mainTag != "" {
		update.MainTag = &mainTag
	}
	err := n.Repo.UpdateThread(ctx, n.ID, update)
	if err != nil {
		return err
	}
	if mainTag != "" {
		n.mainTag = mainTag
	}
	if subTags != nil {
		n.subTags = subTags
	}
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

type Post struct {
	ID        uid.UID   `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Anonymous bool      `json:"anonymous"`
	Author    string    `json:"author"`
	Content   string    `json:"content"`
	Blocked   bool      `json:"blocked"`

	ThreadID uid.UID   `json:"-"`
	QuoteIDs []uid.UID `json:"-"`
	Repo     ForumRepo `json:"-"`
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

		ThreadID: input.ThreadID,
		Repo:     f.Repo,
	}
	if input.Anonymous {
		post.Author = uid.NewUID().ToBase64String()
	} else {
		if user.Name == nil {
			return nil, errors.New("user name must be set")
		}
		post.Author = *user.Name
	}
	err := f.Repo.InsertPost(ctx, user.ID, post)
	return &NewPostResponse{Post: post}, err
}

func (f *ForumService) GetPostByID(ctx context.Context, postID uid.UID) (*Post, error) {
	return f.Repo.GetPost(ctx, &PostSearch{ID: &postID})
}

func (f *ForumService) GetUserPosts(ctx context.Context, user *User, query SliceQuery) (*PostSlice, error) {
	return f.Repo.GetPostSlice(ctx, &PostSearch{UserID: &user.ID}, query)
}

func (p *Post) Quotes(ctx context.Context) ([]*Post, error) {
	if len(p.QuoteIDs) == 0 {
		return nil, nil
	}
	return p.Repo.GetPosts(ctx, &PostSearch{IDs: p.QuoteIDs})
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

type Tag struct {
	Name      string   `json:"name"`
	IsMain    bool     `json:"isMain"`
	BelongsTo []string `json:"belongsTo"`
}

func (f *ForumService) tagsToNames(tags []*Tag) []string {
	var names []string
	for _, t := range tags {
		names = append(names, t.Name)
	}
	return names
}

func (f *ForumService) GetMainTags(ctx context.Context) ([]string, error) {
	tags, err := f.Repo.GetTags(ctx, &TagSearch{MainOnly: true})
	if err != nil {
		return nil, err
	}
	return f.tagsToNames(tags), nil
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

func (f *ForumService) GetUserTags(ctx context.Context, user *User) ([]string, error) {
	if user.tags != nil {
		return user.tags, nil
	}
	tags, err := f.Repo.GetTags(ctx, &TagSearch{UserID: &user.ID})
	if err != nil {
		return nil, err
	}
	if len(tags) == 0 {
		user.tags = []string{}
	} else {
		user.tags = f.tagsToNames(tags)
	}
	return user.tags, nil
}

func (f *ForumService) SyncUserTags(ctx context.Context, user *User, tags []string) error {
	userTags, err := f.GetUserTags(ctx, user)
	if err != nil {
		return err
	}
	curTags := map[string]bool{}
	for _, t := range userTags {
		curTags[t] = true
	}
	nextTags := map[string]bool{}
	for _, t := range tags {
		nextTags[t] = true
	}

	var toAdd []string
	var toDel []string
	for t := range curTags {
		if !nextTags[t] {
			toDel = append(toDel, t)
		}
	}
	for t := range nextTags {
		if !curTags[t] {
			toAdd = append(toAdd, t)
		}
	}
	if err := f.Repo.UpdateUserTags(ctx, user.ID, &UserTagUpdate{AddTags: toAdd, DelTags: toDel}); err != nil {
		return err
	}
	user.tags = tags
	return nil
}

func (f *ForumService) AddUserSubbedTag(ctx context.Context, user *User, tag string) error {
	if err := f.Repo.UpdateUserTags(ctx, user.ID, &UserTagUpdate{AddTags: []string{tag}}); err != nil {
		return err
	}
	if user.tags != nil {
		user.tags = append(user.tags, tag)
	}
	return nil
}

func (f *ForumService) DelUserSubbedTag(ctx context.Context, user *User, tag string) error {
	if err := f.Repo.UpdateUserTags(ctx, user.ID, &UserTagUpdate{AddTags: []string{tag}}); err != nil {
		return err
	}
	if user.tags != nil {
		var next []string
		for _, t := range user.tags {
			if t != tag {
				next = append(next, t)
			}
		}
		user.tags = next
	}
	return nil
}
