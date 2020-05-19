package repo

import (
	"context"

	"github.com/go-pg/pg/v9/orm"

	"gitlab.com/abyss.club/uexky/lib/algo"
	"gitlab.com/abyss.club/uexky/lib/postgres"
	"gitlab.com/abyss.club/uexky/lib/uid"
	"gitlab.com/abyss.club/uexky/uexky/entity"
)

type ForumRepo struct {
	mainTags []string `wire:"-"`
}

const blockedContent = "[此内容已被管理员屏蔽]" // TODO: should be service layer

func (f *ForumRepo) db(ctx context.Context) postgres.Session {
	return postgres.GetSessionFromContext(ctx)
}

func (f *ForumRepo) toEntityThread(t *Thread) *entity.Thread {
	thread := &entity.Thread{
		ID:        uid.UID(t.ID),
		CreatedAt: t.CreatedAt,
		Anonymous: t.Anonymous,
		Title:     t.Title,
		Content:   t.Content,
		MainTag:   t.Tags[0],
		SubTags:   t.Tags,
		Blocked:   t.Blocked,
		Locked:    t.Locked,

		Repo: f,
		AuthorObj: entity.Author{
			UserID: t.UserID,
		},
	}
	if t.Anonymous {
		thread.AuthorObj.AnonymousID = (*uid.UID)(t.AnonymousID)
	} else {
		thread.AuthorObj.UserName = t.UserName
	}
	if thread.Blocked {
		thread.Content = blockedContent
	}
	return thread
}

func (f *ForumRepo) GetThread(ctx context.Context, search *entity.ThreadSearch) (*entity.Thread, error) {
	thread := Thread{}
	q := f.db(ctx).Model(&thread).Where("id = ?", search.ID)
	if err := q.Select(); err != nil {
		return nil, err
	}
	return f.toEntityThread(&thread), nil
}

func (f *ForumRepo) GetThreadSlice(
	ctx context.Context, search *entity.ThreadsSearch, query entity.SliceQuery,
) (*entity.ThreadSlice, error) {
	var threads []Thread
	q := f.db(ctx).Model(&threads)
	if search.Tags != nil {
		q.Where("id IN (SELECT thread_id FROM threads_tags as tt WHERE tt.tag_name=ANY(?))", search.Tags)
	}
	if search.UserID != nil {
		q.Where("user_id = ?", search.UserID)
	}
	applySlice := func(q *orm.Query, isAfter bool, cursor string) (*orm.Query, error) {
		if cursor == "" {
			return q, nil
		}
		c, err := uid.ParseUID(cursor)
		if err != nil {
			return nil, err
		}
		lastPostID := Columns.Thread.LastPostID
		if !isAfter {
			return q.Where("? > ?", lastPostID, c).Order(lastPostID), nil
		}
		return q.Where("? < ?", lastPostID, c).Order("? DESC", lastPostID), nil
	}
	var err error
	q, err = applySliceQuery(applySlice, q, &query)
	if err != nil {
		return nil, err
	}
	if err := q.Select(); err != nil {
		return nil, err
	}

	sliceInfo := &entity.SliceInfo{HasNext: len(threads) > query.Limit}
	var entities []*entity.Thread
	dealSlice := func(i int, isFirst bool, isLast bool) {
		entities = append(entities, f.toEntityThread(&threads[i]))
		if isFirst {
			sliceInfo.FirstCursor = uid.UID(threads[i].LastPostID).ToBase64String()
		}
		if isLast {
			sliceInfo.LastCursor = uid.UID(threads[i].LastPostID).ToBase64String()
		}
	}
	dealSliceResult(dealSlice, &query, len(threads), query.Before != nil)

	return &entity.ThreadSlice{
		Threads:   entities,
		SliceInfo: sliceInfo,
	}, nil
}

func (f *ForumRepo) GetThreadCatelog(ctx context.Context, id uid.UID) ([]*entity.ThreadCatalogItem, error) {
	var posts []Post
	q := f.db(ctx).Model(&posts).Column("id", "created_at").Where("thread=?", id).Order("id")
	if err := q.Select(); err != nil {
		return nil, err
	}
	var cats []*entity.ThreadCatalogItem
	for i := range posts {
		cats = append(cats, &entity.ThreadCatalogItem{
			PostID:    uid.UID(posts[i].ID).ToBase64String(),
			CreatedAt: posts[i].CreatedAt,
		})
	}
	return cats, nil
}

func (f *ForumRepo) GetAnonyID(ctx context.Context, userID int, threadID uid.UID) (uid.UID, error) {
	var posts []Post
	q := f.db(ctx).Model(&posts).Column(Columns.Post.AnonymousID).
		Where("thread_id = ?", threadID).Where("anonymous = true").Order("id DESC").Limit(1)
	if err := q.Select(); err != nil {
		return uid.UID(0), err
	}
	if len(posts) > 0 {
		return uid.UID(*posts[0].AnonymousID), nil
	}
	return uid.NewUID(), nil
}

func (f *ForumRepo) InsertThread(ctx context.Context, thread *entity.Thread) error {
	t := Thread{
		ID:         int64(thread.ID),
		Anonymous:  thread.Anonymous,
		UserID:     thread.AuthorObj.UserID,
		Title:      thread.Title,
		Content:    thread.Content,
		LastPostID: int64(thread.ID),
	}
	if thread.Anonymous {
		t.AnonymousID = (*int64)(thread.AuthorObj.AnonymousID)
	} else {
		t.UserName = thread.AuthorObj.UserName
	}
	t.Tags = []string{thread.MainTag}
	t.Tags = append(t.Tags, thread.SubTags...)
	return f.db(ctx).Insert(&t)
}

func (f *ForumRepo) UpdateThread(ctx context.Context, id uid.UID, update *entity.ThreadUpdate) error {
	thread := Thread{}
	q := f.db(ctx).Model(&thread).Where("id = ?", id)
	if update.Blocked != nil {
		q.Set("blocked = ?Blocked", update)
	}
	if update.Locked != nil {
		q.Set("locked = ?Locked", update)
	}
	if update.MainTag != nil {
		tags := []string{*update.MainTag}
		tags = append(tags, update.SubTags...)
		q.Set("tags = ?", tags)
	}
	_, err := q.Update()
	return err
}

func (f *ForumRepo) toEntityPost(p *Post) *entity.Post {
	post := &entity.Post{
		ID:        uid.UID(p.ID),
		CreatedAt: p.CreatedAt,
		Anonymous: p.Anonymous,
		Content:   p.Content,

		Repo: f,
		Data: entity.PostData{
			ThreadID: uid.UID(p.ThreadID),
			Author: entity.Author{
				UserID:      p.UserID,
				AnonymousID: (*uid.UID)(p.AnonymousID),
				UserName:    p.UserName,
			},
		},
	}
	var qids []uid.UID
	for _, pqid := range p.QuotedIDs {
		qids = append(qids, uid.UID(pqid))
	}
	post.Data.QuoteIDs = qids
	if p.Blocked != nil && *p.Blocked {
		post.Blocked = true
		post.Content = blockedContent // TODO: move to service layer
	}

	return post
}

func (f *ForumRepo) GetPost(ctx context.Context, search *entity.PostSearch) (*entity.Post, error) {
	var post Post
	if err := f.db(ctx).Model(&post).Where("id = ?", search.ID).Select(); err != nil {
		return nil, err
	}
	return f.toEntityPost(&post), nil
}

func (f *ForumRepo) searchPostsQuery(ctx context.Context, search *entity.PostsSearch, posts *[]Post) *orm.Query {
	q := f.db(ctx).Model(posts)
	if search.IDs != nil {
		q.Where("id = ANY(?)", search.IDs)
	}
	if search.UserID != nil {
		q.Where("user_id = ?", search.UserID)
	}
	if search.ThreadID != nil {
		q.Where("thread_id = ?", search.ThreadID)
	}
	return q
}

func (f *ForumRepo) GetPosts(ctx context.Context, search *entity.PostsSearch) ([]*entity.Post, error) {
	var posts []Post
	q := f.searchPostsQuery(ctx, search, &posts)
	if err := q.Select(); err != nil {
		return nil, err
	}
	var ePosts []*entity.Post
	for i := range posts {
		ePosts = append(ePosts, f.toEntityPost(&posts[i]))
	}
	return ePosts, nil
}

func (f *ForumRepo) GetPostSlice(
	ctx context.Context, search *entity.PostsSearch, query entity.SliceQuery,
) (*entity.PostSlice, error) {
	var posts []Post
	q := f.searchPostsQuery(ctx, search, &posts)
	applySlice := func(q *orm.Query, isAfter bool, cursor string) (*orm.Query, error) {
		if cursor == "" {
			return q, nil
		}
		c, err := uid.ParseUID(cursor)
		if err != nil {
			return nil, err
		}
		if !isAfter {
			return q.Where("id < ?", c).Order("id DESC"), nil
		}
		return q.Where("id > ?", c).Order("id"), nil
	}
	var err error
	q, err = applySliceQuery(applySlice, q, &query)
	if err != nil {
		return nil, err
	}
	if err := q.Select(); err != nil {
		return nil, err
	}

	sliceInfo := &entity.SliceInfo{HasNext: len(posts) > query.Limit}
	var entities []*entity.Post
	dealSlice := func(i int, isFirst bool, isLast bool) {
		entities = append(entities, f.toEntityPost(&posts[i]))
		if isFirst {
			sliceInfo.FirstCursor = uid.UID(posts[i].ID).ToBase64String()
		}
		if isLast {
			sliceInfo.LastCursor = uid.UID(posts[i].ID).ToBase64String()
		}
	}
	dealSliceResult(dealSlice, &query, len(posts), query.Before != nil)
	return &entity.PostSlice{
		Posts:     entities,
		SliceInfo: sliceInfo,
	}, nil
}

func (f *ForumRepo) GetPostCount(ctx context.Context, search *entity.PostsSearch) (int, error) {
	var posts []Post
	q := f.searchPostsQuery(ctx, search, &posts)
	return q.Count()
}

func (f *ForumRepo) GetPostQuotesPosts(ctx context.Context, id uid.UID) ([]*entity.Post, error) {
	var posts []Post
	q := f.db(ctx).Model(&posts).Join("INNER JOIN posts_quotes ON post.id = posts_quotes.quoted_id").
		Where("posts_quotes.quoter_id = ?", id).Order("post.id")
	if err := q.Select(); err != nil {
		return nil, err
	}
	var ePosts []*entity.Post
	for i := range posts {
		ePosts = append(ePosts, f.toEntityPost(&posts[i]))
	}
	return ePosts, nil
}

func (f *ForumRepo) GetPostQuotedCount(ctx context.Context, id uid.UID) (int, error) {
	var count int
	_, err := f.db(ctx).Query(orm.Scan(&count), "SELECT count(*) FROM posts_quotes WHERE quoted_id=?", id)
	return count, err
}

func (f *ForumRepo) InsertPost(ctx context.Context, post *entity.Post) error {
	newPost := &Post{
		ID:        int64(post.ID),
		ThreadID:  int64(post.Data.ThreadID),
		Anonymous: post.Anonymous,
		UserID:    post.Data.Author.UserID,
		Content:   post.Content,
	}
	if post.Anonymous {
		newPost.AnonymousID = (*int64)(post.Data.Author.AnonymousID)
	} else {
		newPost.UserName = post.Data.Author.UserName
	}
	var qids []int64
	for _, pqid := range post.Data.QuoteIDs {
		qids = append(qids, int64(pqid))
	}
	newPost.QuotedIDs = qids
	if _, err := f.db(ctx).Model(newPost).Insert(); err != nil {
		return err
	}
	if _, err := f.db(ctx).Model((*Thread)(nil)).Set("last_post_id=?", post.ID).
		Where("id = ?", post.Data.ThreadID).Update(); err != nil {
		return err
	}
	return nil
}

func (f *ForumRepo) UpdatePost(ctx context.Context, id uid.UID, update *entity.PostUpdate) error {
	post := Post{}
	q := f.db(ctx).Model(&post).Where("id = ?", id)
	if update.Blocked != nil {
		q.Set("blocked = ?Blocked", update)
	}
	_, err := q.Update()
	if err != nil {
		return err
	}
	return nil
}

func (f *ForumRepo) GetMainTags(ctx context.Context) ([]string, error) {
	if f.mainTags != nil {
		return f.mainTags, nil
	}
	var tags []Tag
	if err := f.db(ctx).Model(&tags).Where("tag_type = main").Select(); err != nil {
		return nil, err
	}
	var mainTags []string
	for i := range tags {
		mainTags = append(mainTags, tags[i].Name)
	}
	f.mainTags = mainTags
	return mainTags, nil
}

func (f *ForumRepo) GetTags(ctx context.Context, search *entity.TagSearch) ([]*entity.Tag, error) {
	var tags []Tag
	q := f.db(ctx).Model(&tags).Limit(search.Limit)
	if search.Text != nil {
		q.Where("name LIKE '%?%'", *search.Text)
	}
	if search.UserID != nil {
		q.Where("name IN (SELECT tag_name FROM users_tags WHERE user_id=?)", *search.UserID)
	}
	if err := q.Select(); err != nil {
		return nil, err
	}
	mainTags, err := f.GetMainTags(ctx)
	if err != nil {
		return nil, err
	}
	var entities []*entity.Tag
	for _, t := range tags {
		entities = append(entities, &entity.Tag{
			Name:   t.Name,
			IsMain: algo.InStrSlice(mainTags, t.Name),
		})
	}
	return entities, nil
}
