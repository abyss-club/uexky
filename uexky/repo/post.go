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
	"gitlab.com/abyss.club/uexky/lib/errors"
	"gitlab.com/abyss.club/uexky/lib/postgres"
	librd "gitlab.com/abyss.club/uexky/lib/redis"
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

func (r *PostRepo) GetByID(ctx context.Context, id uid.UID) (*entity.Post, error) {
	var post Post
	if err := db(ctx).Model(&post).Where("id = ?", id).Select(); err != nil {
		return nil, postgres.ErrHandlef(err, "GetPost(id=%v)", id)
	}
	return post.ToEntity(), nil
}

func (r *PostRepo) Insert(ctx context.Context, post *entity.Post) (*entity.Post, error) {
	log.Infof("InsertPost(%+v)", post)
	p := NewPostFromEntity(post)
	if _, err := db(ctx).Model(p).Returning("*").Insert(); err != nil {
		return nil, postgres.ErrHandlef(err, "InsertPost.Insert(post=%+v)", post)
	}
	if _, err := db(ctx).Model(&Thread{}).Set("last_post_id=?", post.ID).
		Where("id = ?", post.ThreadID).Update(); err != nil {
		return nil, postgres.ErrHandlef(err, "InsertPost.UpdateThread(post=%+v)", post)
	}
	return p.ToEntity(), nil
}

func (r *PostRepo) Update(ctx context.Context, post *entity.Post) (*entity.Post, error) {
	p := Post{}
	q := db(ctx).Model(&p).Where("id = ?", post.ID).
		Set("blocked = ?", post.Blocked)
	_, err := q.Returning("*").Update()
	return p.ToEntity(), postgres.ErrHandlef(err, "UpdatePost(post=%+v)", p)
}

func (r *PostRepo) QuotedPosts(ctx context.Context, post *entity.Post) ([]*entity.Post, error) {
	var posts []Post
	q := db(ctx).Model(&posts).Where("id = ANY(?)", pg.Array(post.QuoteIDs))
	if err := q.Select(); err != nil {
		return nil, postgres.ErrHandlef(err, "GetPostsQuotedPosts(post=%+v)", post)
	}
	ePostMap := map[uid.UID]*entity.Post{}
	for i := range posts {
		p := (&posts[i]).ToEntity()
		ePostMap[p.ID] = p
	}
	var rst []*entity.Post
	for _, id := range post.QuoteIDs {
		rst = append(rst, ePostMap[id])
	}
	return rst, nil
}

func (r *PostRepo) QuotedCount(ctx context.Context, post *entity.Post) (int, error) {
	var count int
	_, err := db(ctx).Query(orm.Scan(&count), "SELECT count(*) FROM post WHERE ? = ANY(quoted_ids)", post.ID)
	return count, postgres.ErrHandlef(err, "GetPostQuotedCount(id=%v)", post.ID)
}

func getPostSlice(ctx context.Context, qf queryFunc, sq *entity.SliceQuery, desc bool) (*entity.PostSlice, error) {
	var posts []Post
	var entities []*entity.Post
	h := sliceHelper{
		Column: "id",
		Desc:   desc,
		TransCursor: func(s string) (interface{}, error) {
			return uid.ParseUID(s)
		},
		SQ: sq,
	}
	if err := h.Select(qf(db(ctx).Model(&posts))); err != nil {
		return nil, postgres.ErrHandlef(err, "GetPostSlice")
	}
	h.DealResults(len(posts), func(i int) {
		entities = append(entities, (&posts[i]).ToEntity())
	})
	sliceInfo := &entity.SliceInfo{HasNext: len(posts) > sq.Limit}
	if len(entities) > 0 {
		sliceInfo.FirstCursor = entities[0].ID.ToBase64String()
		sliceInfo.LastCursor = entities[len(entities)-1].ID.ToBase64String()
	}
	return &entity.PostSlice{
		Posts:     entities,
		SliceInfo: sliceInfo,
	}, nil
}
