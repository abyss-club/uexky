package repo

import (
	"context"
	"fmt"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	log "github.com/sirupsen/logrus"
	"gitlab.com/abyss.club/uexky/lib/algo"
	"gitlab.com/abyss.club/uexky/lib/postgres"
	"gitlab.com/abyss.club/uexky/lib/uid"
	"gitlab.com/abyss.club/uexky/uexky/entity"
)

type ForumRepo struct {
	mainTags []string `wire:"-"`
}

func (f *ForumRepo) db(ctx context.Context) postgres.Session {
	return postgres.GetSessionFromContext(ctx)
}

func (f *ForumRepo) toEntityThread(t *Thread) *entity.Thread {
	thread := &entity.Thread{
		ID:        uid.UID(t.ID),
		CreatedAt: t.CreatedAt,
		Author: &entity.Author{
			UserID:    t.UserID,
			Anonymous: t.Anonymous,
			Author:    t.Author,
		},
		Title:   t.Title,
		Content: t.Content,
		MainTag: t.Tags[0],
		SubTags: t.Tags[1:],
		Blocked: t.Blocked,
		Locked:  t.Locked,

		Repo: f,
	}
	if thread.Blocked {
		thread.Content = entity.BlockedContent
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
		q.Where("id IN (SELECT id FROM thread WHERE ? && thread.tags)", pg.Array(search.Tags))
	}
	if search.UserID != nil {
		q.Where("user_id = ?", search.UserID)
	}
	applySlice := func(q *orm.Query, isAfter bool, cursor string) (*orm.Query, error) {
		if cursor != "" {
			c, err := uid.ParseUID(cursor)
			if err != nil {
				return nil, err
			}
			if !isAfter {
				q = q.Where("last_post_id > ?", c)
			} else {
				q = q.Where("last_post_id < ?", c)
			}
		}
		if !isAfter {
			return q.Order("last_post_id"), nil
		}
		return q.Order("last_post_id DESC"), nil
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

func (f *ForumRepo) GetThreadCatalog(ctx context.Context, id uid.UID) ([]*entity.ThreadCatalogItem, error) {
	var posts []Post
	q := f.db(ctx).Model(&posts).Column("id", "created_at").Where("thread_id=?", id).Order("id")
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

func (f *ForumRepo) GetAnonyID(ctx context.Context, userID int64, threadID uid.UID) (string, error) {
	var posts []Post
	q := f.db(ctx).Model(&posts).Column("author").
		Where("thread_id = ?", threadID).Where("anonymous = true").Order("id DESC").Limit(1)
	if err := q.Select(); err != nil {
		return "", err
	}
	if len(posts) > 0 {
		return posts[0].Author, nil
	}
	return uid.NewUID().ToBase64String(), nil
}

func (f *ForumRepo) InsertThread(ctx context.Context, thread *entity.Thread) error {
	log.Infof("InsertThread(%v)", thread)
	t := Thread{
		ID:         thread.ID,
		UserID:     thread.Author.UserID,
		Anonymous:  thread.Author.Anonymous,
		Author:     thread.Author.Author,
		Title:      thread.Title,
		Content:    thread.Content,
		LastPostID: thread.ID,
	}
	t.Tags = []string{thread.MainTag}
	t.Tags = append(t.Tags, thread.SubTags...)
	return f.db(ctx).Insert(&t)
}

func (f *ForumRepo) UpdateThread(ctx context.Context, id uid.UID, update *entity.ThreadUpdate) error {
	thread := Thread{}
	q := f.db(ctx).Model(&thread).Where("id = ?", id)
	if update.Blocked != nil {
		q.Set("blocked = ?", update.Blocked)
	}
	if update.Locked != nil {
		q.Set("locked = ?", update.Locked)
	}
	if update.MainTag != nil {
		tags := []string{*update.MainTag}
		tags = append(tags, update.SubTags...)
		q.Set("tags = ?", pg.Array(tags))
	}
	_, err := q.Update()
	return err
}

func (f *ForumRepo) toEntityPost(p *Post) *entity.Post {
	post := &entity.Post{
		ID:        uid.UID(p.ID),
		CreatedAt: p.CreatedAt,
		Author: &entity.Author{
			UserID:    p.UserID,
			Anonymous: p.Anonymous,
			Author:    p.Author,
		},
		Content: p.Content,

		Repo: f,
		Data: entity.PostData{
			ThreadID:   uid.UID(p.ThreadID),
			QuoteIDs:   make([]uid.UID, 0),
			QuotePosts: make([]*entity.Post, 0),
		},
	}
	var qids []uid.UID
	for _, pqid := range p.QuotedIDs {
		qids = append(qids, uid.UID(pqid))
	}
	post.Data.QuoteIDs = qids
	if p.Blocked != nil && *p.Blocked {
		post.Blocked = true
		post.Content = entity.BlockedContent
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
		q.Where("id = ANY(?)", pg.Array(search.IDs))
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
		if cursor != "" {
			c, err := uid.ParseUID(cursor)
			if err != nil {
				return nil, err
			}
			if isAfter != search.DESC {
				q = q.Where("id > ?", c)
			} else {
				q.Where("id < ?", c)
			}
		}
		if isAfter != search.DESC {
			return q.Order("id"), nil
		}
		return q.Order("id DESC"), nil
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
	log.Infof("InsertPost(%v)", post)
	newPost := &Post{
		ID:        post.ID,
		ThreadID:  post.Data.ThreadID,
		UserID:    post.Author.UserID,
		Anonymous: post.Author.Anonymous,
		Author:    post.Author.Author,
		Content:   post.Content,
		QuotedIDs: post.Data.QuoteIDs,
	}
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
		q.Set("blocked = ?", update.Blocked)
	}
	_, err := q.Update()
	if err != nil {
		return err
	}
	return nil
}

func (f *ForumRepo) GetTags(ctx context.Context, search *entity.TagSearch) ([]*entity.Tag, error) {
	type tag struct {
		Tag string `pg:"tag"`
	}
	var tags []tag
	var where, limit string
	if search.Text != "" {
		where = fmt.Sprintf("WHERE tag LIKE '%%%s%%'", search.Text)
	}
	if search.Limit != 0 {
		limit = fmt.Sprintf("LIMIT %v", search.Limit)
	}
	sql := fmt.Sprintf(`SELECT tag FROM (
		SELECT unnest(tags) as tag, max(created_at) as updated_at
		FROM thread group by tag
	) as tags %s ORDER BY updated_at DESC %s`, where, limit)
	if _, err := f.db(ctx).Query(&tags, sql); err != nil {
		return nil, err
	}
	mainTags, err := f.GetMainTags(ctx)
	if err != nil {
		return nil, err
	}
	var entities []*entity.Tag
	for _, t := range tags {
		entities = append(entities, &entity.Tag{
			Name:   t.Tag,
			IsMain: algo.InStrSlice(mainTags, t.Tag),
		})
	}
	return entities, nil
}

func (f *ForumRepo) GetMainTags(ctx context.Context) ([]string, error) {
	if f.mainTags != nil {
		return f.mainTags, nil
	}
	var tags []Tag
	if err := f.db(ctx).Model(&tags).Where("type = ?", "main").Select(); err != nil {
		return nil, dbErrWrap(err, "get main tags")
	}
	var mainTags []string
	for i := range tags {
		mainTags = append(mainTags, tags[i].Name)
	}
	f.mainTags = mainTags
	return mainTags, nil
}

func (f *ForumRepo) SetMainTags(ctx context.Context, tags []string) error {
	var mainTags []Tag
	tagType := "main"
	for _, t := range tags {
		mainTags = append(mainTags, Tag{
			Name:    t,
			TagType: &tagType,
		})
	}
	_, err := f.db(ctx).Model(&mainTags).Insert()
	return dbErrWrap(err, "insert main tags")
}
