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

type PostRepo struct {
	Redis *redis.Client
}

func (r *PostRepo) CheckIfDuplicated(ctx context.Context, userID uid.UID, content string) error {
	msg := fmt.Sprintf("%s:%s", userID.ToBase64String(), content)
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

func (r *PostRepo) GetByID(ctx context.Context, id uid.UID) (*entity.Post, error) {
	var post Post
	if err := db(ctx).Model(&post).Where("id = ?", id).Select(); err != nil {
		return nil, dbErrWrapf(err, "GetPost(id=%v)", id)
	}
	return post.ToEntity(), nil
}

func (r *PostRepo) Insert(ctx context.Context, post *entity.Post) (*entity.Post, error) {
	log.Infof("InsertPost(%+v)", post)
	p := NewPostFromEntity(post)
	if _, err := db(ctx).Model(p).Returning("*").Insert(); err != nil {
		return nil, dbErrWrapf(err, "InsertPost.Insert(post=%+v)", post)
	}
	if _, err := db(ctx).Model(&Thread{}).Set("last_post_id=?", post.ID).
		Where("id = ?", post.ThreadID).Update(); err != nil {
		return nil, dbErrWrapf(err, "InsertPost.UpdateThread(post=%+v)", post)
	}
	return p.ToEntity(), nil
}

func (r *PostRepo) Update(ctx context.Context, post *entity.Post) (*entity.Post, error) {
	p := Post{}
	q := db(ctx).Model(&p).Where("id = ?", post.ID).
		Set("blocked = ?", post.Blocked)
	_, err := q.Returning("*").Update()
	return p.ToEntity(), dbErrWrapf(err, "UpdatePost(post=%+v)", p)
}

func (r *PostRepo) QuotedPosts(ctx context.Context, post *entity.Post) ([]*entity.Post, error) {
	var posts []Post
	q := db(ctx).Model(&posts).Where("id = ANY(?)", pg.Array(post.QuoteIDs))
	if err := q.Select(); err != nil {
		return nil, dbErrWrapf(err, "GetPostsQuotedPosts(post=%+v)", post)
	}
	var ePosts []*entity.Post
	for i := range posts {
		ePosts = append(ePosts, (&posts[i]).ToEntity())
	}
	return ePosts, nil
}

func (r *PostRepo) QuotedCount(ctx context.Context, post *entity.Post) (int, error) {
	var count int
	_, err := db(ctx).Query(orm.Scan(&count), "SELECT count(*) FROM post WHERE ? = ANY(quoted_ids)", post.ID)
	return count, dbErrWrapf(err, "GetPostQuotedCount(id=%v)", post.ID)
}

func getPostSlice(ctx context.Context, qf queryFunc, sq *entity.SliceQuery, desc bool) (*entity.PostSlice, error) {
	var posts []Post
	q := qf(db(ctx).Model(posts))
	applySlice := func(q *orm.Query, isAfter bool, cursor string) (*orm.Query, error) {
		if cursor != "" {
			c, err := uid.ParseUID(cursor)
			if err != nil {
				return nil, errors.Wrapf(err, "ParseUID(%s)", cursor)
			}
			if isAfter != desc {
				q = q.Where("id > ?", c)
			} else {
				q.Where("id < ?", c)
			}
		}
		if isAfter != desc {
			return q.Order("id"), nil
		}
		return q.Order("id DESC"), nil
	}
	var err error
	q, err = applySliceQuery(applySlice, q, sq)
	if err != nil {
		return nil, err
	}
	if err := q.Select(); err != nil {
		return nil, dbErrWrapf(err, "GetPostSlice")
	}

	sliceInfo := &entity.SliceInfo{HasNext: len(posts) > sq.Limit}
	var entities []*entity.Post
	dealSlice := func(i int, isFirst bool, isLast bool) {
		entities = append(entities, (&posts[i]).ToEntity())
		if isFirst {
			sliceInfo.FirstCursor = posts[i].ID.ToBase64String()
		}
		if isLast {
			sliceInfo.LastCursor = posts[i].ID.ToBase64String()
		}
	}
	dealSliceResult(dealSlice, sq, len(posts), sq.Before != nil)
	return &entity.PostSlice{
		Posts:     entities,
		SliceInfo: sliceInfo,
	}, nil
}
