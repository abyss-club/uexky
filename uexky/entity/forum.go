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
	AnonymousID *uid.UID
	UserName    *string
}

type Thread struct {
	ID        uid.UID   `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Anonymous bool      `json:"anonymous"`
	Title     *string   `json:"title"`
	Content   string    `json:"content"`
	Blocked   bool      `json:"blocked"`
	Locked    bool      `json:"locked"`

	Repo      ForumRepo `json:"-"`
	UserID    int       `json:"-"`
	AuthorObj Author    `json:"author"`

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

		Repo:   f.Repo,
		UserID: user.ID,

		mainTag: input.MainTag,
		subTags: input.SubTags,
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
	err = thread.EditTags(ctx, thread.mainTag, thread.subTags)
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
	if !n.Anonymous {
		return *n.AuthorObj.UserName
	}
	return n.AuthorObj.AnonymousID.ToBase64String()
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
	Content   string    `json:"content"`
	Blocked   bool      `json:"blocked"`

	Repo        ForumRepo `json:"-"`
	UserID      int       `json:"-"`
	ThreadID    uid.UID   `json:"-"`
	AuthorObj   Author    `json:"-"`
	QuotedPosts []*Post   `json:"-"`
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

		Repo:     f.Repo,
		UserID:   user.ID,
		ThreadID: input.ThreadID,
	}
	if input.Anonymous {
		aid := uid.NewUID()
		post.AuthorObj.AnonymousID = &aid // TODO
	} else {
		if user.Name == nil {
			return nil, errors.New("user name must be set")
		}
		post.AuthorObj.UserName = user.Name
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
	if !p.Anonymous {
		return *p.AuthorObj.UserName
	}
	return p.AuthorObj.AnonymousID.ToBase64String()
}

func (p *Post) Quotes(ctx context.Context) ([]*Post, error) {
	if p.QuotedPosts != nil {
		return p.QuotedPosts, nil
	}
	quotedPosts, err := p.Repo.GetPostQuotesPosts(ctx, p.ID)
	if err != nil {
		return nil, err
	}
	p.QuotedPosts = quotedPosts
	return p.QuotedPosts, nil
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
