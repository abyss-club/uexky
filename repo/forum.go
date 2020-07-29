package repo

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math/rand"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/go-redis/redis/v7"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gitlab.com/abyss.club/uexky/lib/algo"
	"gitlab.com/abyss.club/uexky/lib/postgres"
	"gitlab.com/abyss.club/uexky/lib/uerr"
	"gitlab.com/abyss.club/uexky/lib/uid"
	"gitlab.com/abyss.club/uexky/uexky/entity"
)

type ForumRepo struct {
	Redis    *redis.Client
	MainTags *MainTag
}

func (f *ForumRepo) db(ctx context.Context) postgres.Session {
	return postgres.GetSessionFromContext(ctx)
}

func (f *ForumRepo) toEntityThread(t *Thread) *entity.Thread {
	thread := &entity.Thread{
		ID:        t.ID,
		CreatedAt: t.CreatedAt,
		Author: &entity.Author{
			UserID:    t.UserID,
			Guest:     t.Guest,
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
		return nil, dbErrWrapf(err, "GetThread(search=%+v)", search)
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
				return nil, errors.Wrapf(err, "GetThreadSlice(search=%+v) parse cursor", search)
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
		return nil, dbErrWrapf(err, "GetThreadSlice(search=%+v, query=%+v)", search, query)
	}

	sliceInfo := &entity.SliceInfo{HasNext: len(threads) > query.Limit}
	var entities []*entity.Thread
	dealSlice := func(i int, isFirst bool, isLast bool) {
		entities = append(entities, f.toEntityThread(&threads[i]))
		if isFirst {
			sliceInfo.FirstCursor = threads[i].LastPostID.ToBase64String()
		}
		if isLast {
			sliceInfo.LastCursor = threads[i].LastPostID.ToBase64String()
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
		return nil, dbErrWrapf(err, "GetThreadCatalog(id=%v)", id)
	}
	var cats []*entity.ThreadCatalogItem
	for i := range posts {
		cats = append(cats, &entity.ThreadCatalogItem{
			PostID:    posts[i].ID.ToBase64String(),
			CreatedAt: posts[i].CreatedAt,
		})
	}
	return cats, nil
}

func (f *ForumRepo) GetAnonyID(ctx context.Context, userID uid.UID, threadID uid.UID) (string, error) {
	var posts []Post
	q := f.db(ctx).Model(&posts).Column("author").
		Where("thread_id = ?", threadID).Where("anonymous = true").Order("id DESC").Limit(1)
	if err := q.Select(); err != nil {
		return "", dbErrWrapf(err, "GetAnonyID(userID=%v, threadID=%v", userID, threadID)
	}
	if len(posts) > 0 {
		return posts[0].Author, nil
	}
	return uid.NewUID().ToBase64String(), nil
}

func (f *ForumRepo) InsertThread(ctx context.Context, thread *entity.Thread) (*entity.Thread, error) {
	log.Infof("InsertThread(%v)", thread)
	t := Thread{
		ID:         thread.ID,
		UserID:     thread.Author.UserID,
		Guest:      thread.Author.Guest,
		Anonymous:  thread.Author.Anonymous,
		Author:     thread.Author.Author,
		Title:      thread.Title,
		Content:    thread.Content,
		LastPostID: thread.ID,
	}
	t.Tags = []string{thread.MainTag}
	t.Tags = append(t.Tags, thread.SubTags...)
	if _, err := f.db(ctx).Model(&t).Returning("*").Insert(); err != nil {
		return nil, dbErrWrapf(err, "InsertThread(thread=%+v)", thread)
	}
	return f.toEntityThread(&t), nil
}

func (f *ForumRepo) UpdateThread(ctx context.Context, t *entity.Thread) (*entity.Thread, error) {
	thread := Thread{}
	tags := []string{t.MainTag}
	tags = append(tags, t.SubTags...)
	q := f.db(ctx).Model(&thread).Where("id = ?", t.ID).
		Set("tags = ?", pg.Array(tags)).
		Set("blocked = ?", t.Blocked).
		Set("locked = ?", t.Locked)
	_, err := q.Returning("*").Update()
	return f.toEntityThread(&thread), dbErrWrapf(err, "UpdateThread(thread=%+v)", t)
}

func (f *ForumRepo) toEntityPost(p *Post) *entity.Post {
	post := &entity.Post{
		ID:        p.ID,
		CreatedAt: p.CreatedAt,
		Author: &entity.Author{
			UserID:    p.UserID,
			Guest:     p.Guest,
			Anonymous: p.Anonymous,
			Author:    p.Author,
		},
		Content: p.Content,
		Blocked: p.Blocked,

		Repo: f,
		Data: entity.PostData{
			ThreadID:   p.ThreadID,
			QuoteIDs:   p.QuotedIDs,
			QuotePosts: make([]*entity.Post, 0),
		},
	}
	if post.Blocked {
		post.Content = entity.BlockedContent
	}
	return post
}

func (f *ForumRepo) GetPost(ctx context.Context, search *entity.PostSearch) (*entity.Post, error) {
	var post Post
	if err := f.db(ctx).Model(&post).Where("id = ?", search.ID).Select(); err != nil {
		return nil, dbErrWrapf(err, "GetPost(search=%+v)", search)
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
		return nil, dbErrWrapf(err, "GetPosts(search=%+v)", search)
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
				return nil, errors.Wrapf(err, "GetPostSlice(search=%+v, query=%+v) parse cursor", search, query)
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
		return nil, dbErrWrapf(err, "GetPostSlice(search=%+v, query=%+v)", search, query)
	}

	sliceInfo := &entity.SliceInfo{HasNext: len(posts) > query.Limit}
	var entities []*entity.Post
	dealSlice := func(i int, isFirst bool, isLast bool) {
		entities = append(entities, f.toEntityPost(&posts[i]))
		if isFirst {
			sliceInfo.FirstCursor = posts[i].ID.ToBase64String()
		}
		if isLast {
			sliceInfo.LastCursor = posts[i].ID.ToBase64String()
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

func (f *ForumRepo) GetPostQuotedCount(ctx context.Context, id uid.UID) (int, error) {
	var count int
	_, err := f.db(ctx).Query(orm.Scan(&count), "SELECT count(*) FROM post WHERE ? = ANY(quoted_ids)", id)
	return count, dbErrWrapf(err, "GetPostQuotedCount(id=%v)", id)
}

func (f *ForumRepo) InsertPost(ctx context.Context, post *entity.Post) (*entity.Post, error) {
	log.Infof("InsertPost(%v)", post)
	newPost := &Post{
		ID:        post.ID,
		ThreadID:  post.Data.ThreadID,
		UserID:    post.Author.UserID,
		Guest:     post.Author.Guest,
		Anonymous: post.Author.Anonymous,
		Author:    post.Author.Author,
		Content:   post.Content,
		QuotedIDs: post.Data.QuoteIDs,
	}
	if _, err := f.db(ctx).Model(newPost).Returning("*").Insert(); err != nil {
		return nil, dbErrWrapf(err, "InsertPost.Insert(post=%+v)", post)
	}
	if _, err := f.db(ctx).Model((*Thread)(nil)).Set("last_post_id=?", post.ID).
		Where("id = ?", post.Data.ThreadID).Update(); err != nil {
		return nil, dbErrWrapf(err, "InsertPost.UpdateThread(post=%+v)", post)
	}
	return f.toEntityPost(newPost), nil
}

func (f *ForumRepo) UpdatePost(ctx context.Context, p *entity.Post) (*entity.Post, error) {
	post := Post{}
	q := f.db(ctx).Model(&post).Where("id = ?", p.ID).
		Set("blocked = ?", p.Blocked)
	_, err := q.Returning("*").Update()
	return f.toEntityPost(&post), dbErrWrapf(err, "UpdatePost(post=%+v)", p)
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
		return nil, dbErrWrapf(err, "GetTags(search=%+v)", search)
	}
	var entities []*entity.Tag
	for _, t := range tags {
		entities = append(entities, &entity.Tag{
			Name:   t.Tag,
			IsMain: algo.InStrSlice(f.MainTags.Tags, t.Tag),
		})
	}
	return entities, nil
}

func (f *ForumRepo) GetMainTags(ctx context.Context) []string {
	return f.MainTags.Tags
}

func (f *ForumRepo) SetMainTags(ctx context.Context, tags []string) error {
	return f.MainTags.SetMainTags(ctx, tags)
}

func (f *ForumRepo) CheckDuplicate(ctx context.Context, userID uid.UID, title, content string) error {
	msg := fmt.Sprintf("%s:%s:%s", userID.ToBase64String(), title, content)
	key := fmt.Sprintf("%x", sha256.Sum256([]byte(msg)))
	value := fmt.Sprintf("%v", rand.Int63())
	if _, err := f.Redis.SetNX(key, value, entity.DuplicatedCheckRange).Result(); err != nil {
		return redisErrWrapf(err, "CheckDuplicate, SetNX(%s, %s)", key, value)
	}
	got, err := f.Redis.Get(key).Result()
	if err != nil {
		return redisErrWrapf(err, "CheckDuplicate, Get(%s)", key)
	}
	if got != value { // value already exist
		return uerr.New(uerr.DuplicatedError, "content is duplicated in 5 minutes")
	}
	return nil
}

type MainTag struct {
	Tags []string
}

func NewMainTag(tx *postgres.TxAdapter) (*MainTag, error) {
	var tags []Tag
	if err := tx.DB.Model(&tags).Where("type = ?", "main").Select(); err != nil {
		return nil, dbErrWrap(err, "NewMainTag()")
	}
	var mainTags []string
	for i := range tags {
		mainTags = append(mainTags, tags[i].Name)
	}
	log.Debugf("NewMainTag = %+v", mainTags)
	return &MainTag{Tags: mainTags}, nil
}

func (mt *MainTag) SetMainTags(ctx context.Context, tags []string) error {
	var mainTags []Tag
	tagType := "main"
	for _, t := range tags {
		mainTags = append(mainTags, Tag{
			Name:    t,
			TagType: &tagType,
		})
	}
	db := postgres.GetSessionFromContext(ctx)
	if _, err := db.Model(&mainTags).Insert(); err != nil {
		return dbErrWrapf(err, "SetMainTags(tags=%v)", tags)
	}
	mt.Tags = tags
	return nil
}
