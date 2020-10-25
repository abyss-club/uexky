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
	"gitlab.com/abyss.club/uexky/lib/uerr"
	"gitlab.com/abyss.club/uexky/lib/uid"
	"gitlab.com/abyss.club/uexky/uexky/entity"
)

type ThreadRepo struct {
	Redis *redis.Client
}

func (r *ThreadRepo) CheckIfDuplicated(ctx context.Context, title *string, content string) error {
	msg := fmt.Sprintf("%s:%s", title, content)
	key := fmt.Sprintf("%x", sha256.Sum256([]byte(msg)))
	value := fmt.Sprintf("%v", rand.Int63())
	if _, err := r.Redis.SetNX(key, value, entity.DuplicatedCheckRange).Result(); err != nil {
		return redisErrWrapf(err, "CheckDuplicate, SetNX(%s, %s)", key, value)
	}
	got, err := r.Redis.Get(key).Result()
	if err != nil {
		return redisErrWrapf(err, "CheckDuplicate, Get(%s)", key)
	}
	if got != value { // value already exist
		return uerr.New(uerr.DuplicatedError, "content is duplicated in 5 minutes")
	}
	return nil
}

func (r *ThreadRepo) GetByID(ctx context.Context, id uid.UID) (*entity.Thread, error) {
	thread := Thread{}
	q := db(ctx).Model(&thread).Where("id = ?", id)
	if err := q.Select(); err != nil {
		return nil, dbErrWrapf(err, "GetByID(id=%+v)", id)
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
	if _, err := db(ctx).Model(&t).Returning("*").Insert(); err != nil {
		return nil, dbErrWrapf(err, "InsertThread(thread=%+v)", thread)
	}
	return t.ToEntity(), nil
}

func (r *ThreadRepo) Update(ctx context.Context, thread *entity.Thread) (*entity.Thread, error) {
	t := NewThreadFromEntity(thread)
	q := db(ctx).Model(&t).Where("id = ?", t.ID).
		Set("tags = ?", pg.Array(t.Tags)).
		Set("blocked = ?", t.Blocked).
		Set("locked = ?", t.Locked)
	_, err := q.Returning("*").Update()
	return t.ToEntity(), dbErrWrapf(err, "UpdateThread(thread=%+v)", t)
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
	return count, dbErrWrap(err, "GetThreadReplyCount")
}

func (r *ThreadRepo) Catalog(ctx context.Context, thread *entity.Thread) ([]*entity.ThreadCatalogItem, error) {
	var posts []Post
	q := db(ctx).Model(&posts).Column("id", "created_at").Where("thread_id=?", thread.ID).Order("id")
	if err := q.Select(); err != nil {
		return nil, dbErrWrapf(err, "GetThreadCatalog(id=%v)", thread.ID)
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
		return "", dbErrWrapf(err, "GetAnonyID(userID=%v, threadID=%v", user.ID, thread.ID)
	}
	if len(posts) > 0 {
		return posts[0].Author, nil
	}
	return uid.NewUID().ToBase64String(), nil
}

func getThreadSlice(ctx context.Context, qf queryFunc, sq *entity.SliceQuery) (*entity.ThreadSlice, error) {
	var threads []Thread
	q := db(ctx).Model(&threads)
	q = qf(q)
	applySlice := func(q *orm.Query, isAfter bool, cursor string) (*orm.Query, error) {
		if cursor != "" {
			c, err := uid.ParseUID(cursor)
			if err != nil {
				return nil, errors.Wrapf(err, "ParseUID(%s)", cursor)
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
	q, err = applySliceQuery(applySlice, q, sq)
	if err != nil {
		return nil, err
	}
	if err := q.Select(); err != nil {
		return nil, dbErrWrap(err, "GetThreadSlice")
	}

	sliceInfo := &entity.SliceInfo{HasNext: len(threads) > sq.Limit}
	var entities []*entity.Thread
	dealSlice := func(i int, isFirst bool, isLast bool) {
		entities = append(entities, (&threads[i]).ToEntity())
		if isFirst {
			sliceInfo.FirstCursor = threads[i].LastPostID.ToBase64String()
		}
		if isLast {
			sliceInfo.LastCursor = threads[i].LastPostID.ToBase64String()
		}
	}
	dealSliceResult(dealSlice, sq, len(threads), sq.Before != nil)
	return &entity.ThreadSlice{
		Threads:   entities,
		SliceInfo: sliceInfo,
	}, nil
}
