package repo

import (
	"context"
	"errors"

	pg "github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"gitlab.com/abyss.club/uexky/lib/uid"
	"gitlab.com/abyss.club/uexky/uexky/entity"
)

type Forum struct {
	DB *pg.DB
}

func NewForum(uri string) (*Forum, error) {
	option, err := pg.ParseURL(uri)
	if err != nil {
		return nil, err
	}
	return &Forum{
		DB: pg.Connect(option),
	}, nil
}

const blockedContent = "[此内容已被管理员屏蔽]" // TODO: should be service layer

func (f *Forum) toEntityThread(t *Thread) *entity.Thread {
	thread := &entity.Thread{
		ID:        uid.UID(t.ID),
		CreatedAt: t.CreatedAt,
		Anonymous: t.Anonymous,
		Title:     t.Title,
		Content:   t.Content,
		Blocked:   t.Blocked,
		Locked:    t.Locked,

		Repo:   f,
		UserID: t.UserID,
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

func (f *Forum) GetThread(ctx context.Context, search *entity.ThreadSearch) (*entity.Thread, error) {
	thread := Thread{}
	q := f.DB.Model(&thread).Where("id = ?", search.ID)
	if err := q.Select(); err != nil {
		return nil, err
	}
	return f.toEntityThread(&thread), nil
}

func (f *Forum) GetThreadSlice(
	ctx context.Context, search *entity.ThreadsSearch, query entity.SliceQuery,
) (*entity.ThreadSlice, error) {
	var threads []Thread
	q := f.DB.Model(&threads)
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
	var eThreads []*entity.Thread
	dealSlice := func(i int, isFirst bool, isLast bool) {
		eThreads = append(eThreads, f.toEntityThread(&threads[i]))
		if isFirst {
			sliceInfo.FirstCursor = uid.UID(threads[i].LastPostID).ToBase64String()
		}
		if isLast {
			sliceInfo.LastCursor = uid.UID(threads[i].LastPostID).ToBase64String()
		}
	}
	dealSliceResult(dealSlice, &query, len(threads), query.Before != nil)

	return &entity.ThreadSlice{
		Threads:   eThreads,
		SliceInfo: sliceInfo,
	}, nil
}

func (f *Forum) GetThreadCatelog(ctx context.Context, id uid.UID) ([]*entity.ThreadCatalogItem, error) {
	var posts []Post
	q := f.DB.Model(&posts).Column("id", "created_at").Where("thread=?", id).Order("id")
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

func (f *Forum) GetThreadTags(ctx context.Context, id uid.UID) (main string, subs []string, err error) {
	var tags []Tag
	if err := f.DB.Model(&tags).
		Join("INNER JOIN threads_tags ON threads_tags.tag_name = tag.name").
		Where("threads_tags.threads_id = ?", id).Select(); err != nil {
		return "", nil, err
	}
	for _, t := range tags {
		if t.IsMain {
			main = t.Name
		} else {
			subs = append(subs, t.Name)
		}
	}
	return
}

func (f *Forum) GetAnonyID(ctx context.Context, userID int, threadID uid.UID) (uid.UID, error) {
	aid := AnonymouId{
		ThreadID:    int64(threadID),
		UserID:      userID,
		AnonymousID: int64(uid.NewUID()),
	}
	q := f.DB.Model(&aid).OnConflict("(thread_id, user_id) DO UPDATE").
		Set("updated_at = now()").Returning("*")
	if _, err := q.Insert(); err != nil {
		return 0, err
	}
	return uid.UID(aid.AnonymousID), nil
}

func (f *Forum) InsertThread(ctx context.Context, thread *entity.Thread) error {
	t := Thread{
		ID:         int64(thread.ID),
		Anonymous:  thread.Anonymous,
		UserID:     thread.UserID,
		Title:      thread.Title,
		Content:    thread.Content,
		LastPostID: int64(thread.ID),
	}
	if thread.Anonymous {
		t.AnonymousID = (*int64)(thread.AuthorObj.AnonymousID)
	} else {
		t.UserName = thread.AuthorObj.UserName
	}
	return f.DB.Insert(&t)
}

func (f *Forum) UpdateThread(ctx context.Context, id uid.UID, update *entity.ThreadUpdate) error {
	if update.Blocked != nil || update.Locked != nil {
		thread := Thread{}
		q := f.DB.Model(&thread).Where("id = ?", id)
		if update.Blocked != nil {
			q.Set("blocked = ?Blocked", update)
		}
		if update.Locked != nil {
			q.Set("locked = ?Locked", update)
		}
		_, err := q.Update()
		if err != nil {
			return err
		}
	}
	if update.MainTag != nil || update.SubTags != nil {
		if !(update.MainTag != nil && update.SubTags != nil) {
			return errors.New("must specify both main tag and sub tags")
		}
		if err := f.setThreadTags(ctx, id, update, false); err != nil {
			return err
		}
	}
	return nil
}

func (f *Forum) setThreadTags(_ context.Context, id uid.UID, update *entity.ThreadUpdate, isNew bool) error {
	tags := update.SubTags
	if update.MainTag != nil {
		tags = append(tags, *update.MainTag)
	}
	// validate params
	var mainTags []Tag
	if err := f.DB.Model(&mainTags).Where("is_main = true").Where("name = ANY(?)", tags).Select(); err != nil {
		return err
	}
	if len(mainTags) != 1 {
		return errors.New("one and only one main tag should be specified")
	}
	if !isNew {
		delTags := []ThreadsTag{}
		if _, err := f.DB.Model(&delTags).Where("thread_id = ?", id).Delete(); err != nil {
			return err
		}
	}
	for _, ut := range tags {
		isMain := ut == *update.MainTag
		if _, err := f.DB.Model(&Tag{Name: ut, IsMain: isMain}).
			OnConflict("(thread_id, tag_name) DO UPDATE").
			Set("update_at = now()").Insert(); err != nil {
			return err
		}
		if _, err := f.DB.Model(&ThreadsTag{ThreadID: int64(id), TagName: ut}).
			OnConflict("(thread_id, tag_name) DO UPDATE").
			Set("update_at = now()").Insert(); err != nil {
			return err
		}
		if !isMain {
			if _, err := f.DB.Model(&TagsMainTag{Name: ut, BelongsTo: *update.MainTag}).
				OnConflict("(name, belongs_to) DO UPDATE").
				Set("update_at = now()").Insert(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (f *Forum) toEntityPost(p *Post) *entity.Post {
	post := &entity.Post{
		ID:        uid.UID(p.ID),
		CreatedAt: p.CreatedAt,
		Anonymous: p.Anonymous,
		Content:   p.Content,

		Repo:     f,
		UserID:   p.UserID,
		ThreadID: uid.UID(p.ThreadID),
		AuthorObj: entity.Author{
			AnonymousID: (*uid.UID)(p.AnonymousID),
			UserName:    p.UserName,
		},
	}
	if p.Blocked != nil && *p.Blocked {
		post.Blocked = true
		post.Content = blockedContent // TODO: move to service layer
	}

	return post
}

func (f *Forum) GetPost(ctx context.Context, search *entity.PostSearch) (*entity.Post, error) {
	var post Post
	if err := f.DB.Model(&post).Where("id = ?", search.ID).Select(); err != nil {
		return nil, err
	}
	return f.toEntityPost(&post), nil
}

func (f *Forum) searchPostsQuery(search *entity.PostsSearch, posts *[]Post) *orm.Query {
	q := f.DB.Model(posts)
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

func (f *Forum) GetPosts(ctx context.Context, search *entity.PostsSearch) ([]*entity.Post, error) {
	var posts []Post
	q := f.searchPostsQuery(search, &posts)
	if err := q.Select(); err != nil {
		return nil, err
	}
	var ePosts []*entity.Post
	for i := range posts {
		ePosts = append(ePosts, f.toEntityPost(&posts[i]))
	}
	return ePosts, nil
}

func (f *Forum) GetPostSlice(
	ctx context.Context, search *entity.PostsSearch, query entity.SliceQuery,
) (*entity.PostSlice, error) {
	var posts []Post
	q := f.searchPostsQuery(search, &posts)
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

func (f *Forum) GetPostCount(ctx context.Context, search *entity.PostsSearch) (int, error) {
	var posts []Post
	q := f.searchPostsQuery(search, &posts)
	return q.Count()
}

func (f *Forum) GetPostQuotesPosts(ctx context.Context, id uid.UID) ([]*entity.Post, error) {
	var posts []Post
	q := f.DB.Model(&posts).Join("INNER JOIN posts_quotes ON post.id = posts_quotes.quoted_id").
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

func (f *Forum) GetPostQuotedCount(ctx context.Context, id uid.UID) (int, error) {
	var count int
	_, err := f.DB.Query(orm.Scan(&count), "SELECT count(*) FROM posts_quotes WHERE quoted_id=?", id)
	return count, err
}

func (f *Forum) InsertPost(ctx context.Context, post *entity.Post) error {
	newPost := &Post{
		ID:        int64(post.ID),
		ThreadID:  int64(post.ThreadID),
		Anonymous: post.Anonymous,
		UserID:    post.UserID,
		Content:   post.Content,
	}
	if post.Anonymous {
		newPost.AnonymousID = (*int64)(post.AuthorObj.AnonymousID)
	} else {
		newPost.UserName = post.AuthorObj.UserName
	}
	if _, err := f.DB.Model(newPost).Insert(); err != nil {
		return err
	}
	if _, err := f.DB.Model((*Thread)(nil)).Set("last_post_id=?", post.ID).
		Where("id = ?", post.ThreadID).Update(); err != nil {
		return err
	}
	quotes, err := post.Quotes(ctx)
	if err != nil {
		return err
	}
	for _, q := range quotes {
		pq := PostsQuote{
			QuoterID: int64(post.ID),
			QuotedID: int64(q.ID),
		}
		if _, err := f.DB.Model(&pq).Insert(); err != nil {
			return err
		}
	}
	return nil
}

func (f *Forum) UpdatePost(ctx context.Context, id uid.UID, update *entity.PostUpdate) error {
	post := Post{}
	q := f.DB.Model(&post).Where("id = ?", id)
	if update.Blocked != nil {
		q.Set("blocked = ?Blocked", update)
	}
	_, err := q.Update()
	if err != nil {
		return err
	}
	return nil
}

func (f *Forum) GetTags(ctx context.Context, search *entity.TagSearch) ([]*entity.Tag, error) {
	var tags []Tag
	q := f.DB.Model(&tags).Limit(search.Limit)
	if search.Text != nil {
		q.Where("name LIKE '%?%'", *search.Text)
	}
	if search.UserID != nil {
		q.Where("name IN (SELECT tag_name FROM users_tags WHERE user_id=?)", *search.UserID)
	}
	if search.MainOnly {
		q.Where("is_main=true")
	}
	if err := q.Select(); err != nil {
		return nil, err
	}
	var entities []*entity.Tag
	for _, t := range tags {
		entities = append(entities, &entity.Tag{
			Name:   t.Name,
			IsMain: t.IsMain,
		})
	}
	return entities, nil
}

func (f *Forum) UpdateUserTags(ctx context.Context, userID int, update *entity.UserTagUpdate) error {
	for _, t := range update.AddTags {
		if _, err := f.DB.Model(&UsersTag{UserID: userID, TagName: t}).Insert(); err != nil {
			return err
		}
	}
	for _, t := range update.DelTags {
		if _, err := f.DB.Model(&UsersTag{}).Where("user_id=?", userID).Where("tag_name=?", t).Delete(); err != nil {
			return err
		}
	}
	return nil
}
