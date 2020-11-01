package repo

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math/rand"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/go-redis/redis/v7"
	log "github.com/sirupsen/logrus"
	"gitlab.com/abyss.club/uexky/lib/algo"
	"gitlab.com/abyss.club/uexky/lib/errors"
	"gitlab.com/abyss.club/uexky/lib/postgres"
	librd "gitlab.com/abyss.club/uexky/lib/redis"
	"gitlab.com/abyss.club/uexky/lib/uid"
	"gitlab.com/abyss.club/uexky/uexky/entity"
)

type ThreadRepo struct {
	Redis *redis.Client
}

func (r *ThreadRepo) CheckIfDuplicated(ctx context.Context, title *string, content string) error {
	msg := fmt.Sprintf("%s:%s", algo.NullToString(title), content)
	key := fmt.Sprintf("%x", sha256.Sum256([]byte(msg)))
	value := fmt.Sprintf("%v", rand.Int63())
	if _, err := r.Redis.SetNX(key, value, entity.DuplicatedCheckRange).Result(); err != nil {
		return librd.ErrHandlef(err, "CheckDuplicate, SetNX(%s, %s)", key, value)
	}
	got, err := r.Redis.Get(key).Result()
	if err != nil {
		return librd.ErrHandlef(err, "CheckDuplicate, Get(%s)", key)
	}
	if got != value { // value already exist
		return errors.Duplicated.New("content is duplicated in 5 minutes")
	}
	return nil
}

func (r *ThreadRepo) GetByID(ctx context.Context, id uid.UID) (*entity.Thread, error) {
	thread := Thread{}
	q := db(ctx).Model(&thread).Where("id = ?", id)
	if err := q.Select(); err != nil {
		return nil, postgres.ErrHandlef(err, "GetByID(id=%+v)", id)
	}
	return thread.ToEntity(), nil
}

func (r *ThreadRepo) FindSlice(
	ctx context.Context, params *entity.ThreadsSearch, query entity.SliceQuery,
) (*entity.ThreadSlice, error) {
	qf := func(prev *orm.Query) *orm.Query {
		return prev.Where("id IN (SELECT id FROM thread WHERE ? && thread.tags)", pg.Array(params.Tags))
	}
	return getThreadSlice(ctx, qf, &query)
}

func (r *ThreadRepo) Insert(ctx context.Context, thread *entity.Thread) (*entity.Thread, error) {
	log.Infof("InsertThread(%v)", thread)
	t := NewThreadFromEntity(thread)
	t.LastPostID = t.ID
	if _, err := db(ctx).Model(t).Returning("*").Insert(); err != nil {
		return nil, postgres.ErrHandlef(err, "InsertThread(thread=%+v)", thread)
	}
	return t.ToEntity(), nil
}

func (r *ThreadRepo) Update(ctx context.Context, thread *entity.Thread) (*entity.Thread, error) {
	t := NewThreadFromEntity(thread)
	q := db(ctx).Model(t).Where("id = ?", t.ID).
		Set("tags = ?", pg.Array(t.Tags)).
		Set("blocked = ?", t.Blocked).
		Set("locked = ?", t.Locked)
	_, err := q.Returning("*").Update()
	return t.ToEntity(), postgres.ErrHandlef(err, "UpdateThread(thread=%+v)", t)
}

func (r *ThreadRepo) Replies(ctx context.Context, thread *entity.Thread, sq entity.SliceQuery) (*entity.PostSlice, error) {
	qf := func(prev *orm.Query) *orm.Query {
		return prev.Where("thread_id = ?", thread.ID)
	}
	return getPostSlice(ctx, qf, &sq, false)
}

func (r *ThreadRepo) ReplyCount(ctx context.Context, thread *entity.Thread) (int, error) {
	var posts []Post
	q := db(ctx).Model(&posts).Where("thread_id = ?", thread.ID)
	count, err := q.Count()
	return count, postgres.ErrHandle(err, "GetThreadReplyCount")
}

func (r *ThreadRepo) Catalog(ctx context.Context, thread *entity.Thread) ([]*entity.ThreadCatalogItem, error) {
	var posts []Post
	q := db(ctx).Model(&posts).Column("id", "created_at").Where("thread_id=?", thread.ID).Order("id")
	if err := q.Select(); err != nil {
		return nil, postgres.ErrHandlef(err, "GetThreadCatalog(id=%v)", thread.ID)
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

func (r *ThreadRepo) PostAID(ctx context.Context, thread *entity.Thread, user *entity.User) (string, error) {
	var posts []Post
	q := db(ctx).Model(&posts).Column("author").
		Where("thread_id = ?", thread.ID).Where("anonymous = true").Order("id DESC").Limit(1)
	if err := q.Select(); err != nil {
		return "", postgres.ErrHandlef(err, "GetAnonyID(userID=%v, threadID=%v", user.ID, thread.ID)
	}
	if len(posts) > 0 {
		return posts[0].Author, nil
	}
	return uid.NewUID().ToBase64String(), nil
}

func getThreadSlice(ctx context.Context, qf queryFunc, sq *entity.SliceQuery) (*entity.ThreadSlice, error) {
	var threads []Thread
	var entities []*entity.Thread
	h := sliceHelper{
		Column: "last_post_id",
		Desc:   true,
		TransCursor: func(s string) (interface{}, error) {
			return uid.ParseUID(s)
		},
		SQ: sq,
	}
	if err := h.Select(qf(db(ctx).Model(&threads))); err != nil {
		return nil, postgres.ErrHandle(err, "GetThreadSlice")
	}
	h.DealResults(len(threads), func(i int) {
		entities = append(entities, (&threads[i]).ToEntity())
	})
	sliceInfo := &entity.SliceInfo{HasNext: len(threads) > sq.Limit}
	if len(entities) > 0 {
		sliceInfo.FirstCursor = entities[0].LastPostID.ToBase64String()
		sliceInfo.LastCursor = entities[len(entities)-1].LastPostID.ToBase64String()
	}
	return &entity.ThreadSlice{
		Threads:   entities,
		SliceInfo: sliceInfo,
	}, nil
}
