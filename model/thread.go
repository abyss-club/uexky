package model

import (
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky/config"
	"gitlab.com/abyss.club/uexky/uexky"
	"gitlab.com/abyss.club/uexky/uuid64"
)

var postIDGenerator = uuid64.Generator{Sections: []uuid64.Section{
	&uuid64.TimestampSection{Length: 6, Unit: time.Second, NoPadding: true},
	&uuid64.CounterSection{Length: 2, Unit: time.Second},
	&uuid64.RandomSection{Length: 1},
}}

// Thread ...
type Thread struct {
	ObjectID   bson.ObjectId `bson:"_id"` // not display in front
	ID         string        `bson:"id"`
	Anonymous  bool          `bson:"anonymous"`
	Author     string        `bson:"author"`
	UserID     bson.ObjectId `bson:"user_id"` // not display in front
	CreateTime time.Time     `bson:"created_time"`
	UpdateTime time.Time     `bson:"update_time"`
	Blocked    bool          `bson:"blocked"`
	Banned     bool          `bson:"banned"`

	MainTag string   `bson:"main_tag"`
	SubTags []string `bson:"sub_tags"`
	Tags    []string `bson:"tags"`
	Title   string   `bson:"title"`
	Content string   `bson:"content"`

	Posts []struct {
		PostID string        `bson:"post_id"`
		UserID string        `bson:"user_id"`
		Author bson.ObjectId `bson:"author"`
	} `bson:"posts"`
}

//
// Model Operations
//

// InsertThread init new thread and insert to db
func InsertThread(u *uexky.Uexky, input *ThreadInput) (*Thread, error) {
	if err := u.Flow.CostMut(config.Config.RateLimit.Cost.PubThread); err != nil {
		return nil, err
	}
	user, err := GetSignedInUser(u)
	if err != nil {
		return nil, err
	}
	if user.Role.Type == Banned {
		return nil, errors.New("permitted error, you are banned")
	}

	thread, err := input.ParseThead(u, user)
	if err != nil {
		return nil, err
	}

	c := u.Mongo.C(colleThread)
	if err := c.Insert(thread); err != nil {
		return nil, err
	}

	// Set Tag info
	if err := UpsertTags(u, thread.MainTag, thread.SubTags); err != nil {
		return nil, errors.Wrap(err, "set tag info")
	}
	return thread, nil
}

// UpdateThread ...
func UpdateThread(u *uexky.Uexky, selector, query bson.M) error {
	return u.Mongo.C(colleThread).Update(selector, query)
}

// FindThread ...
func FindThread(u *uexky.Uexky, selector interface{}) (*Thread, error) {
	thread := &Thread{}
	if err := u.Mongo.C(colleThread).Find(selector).One(thread); err != nil {
		return nil, err
	}
	return thread, nil
}

// FindThreads ...
func FindThreads(u *uexky.Uexky, selector bson.M, sq *SliceQuery) ([]*Thread, *SliceInfo, error) {
	sliceKey := "update_time"
	qry, err := sq.GenQueryByTime(sliceKey)
	if err != nil {
		return nil, nil, err
	}
	for k, v := range selector {
		qry[k] = v
	}
	qry["blocked"] = false

	var threads []*Thread
	if err := sq.Find(u, colleThread, sliceKey, qry, &threads); err != nil {
		return nil, nil, err
	}
	if len(threads) == 0 {
		return threads, &SliceInfo{}, nil
	}
	if !sq.Desc {
		ReverseSlice(threads)
	}
	return threads, &SliceInfo{
		FirstCursor: threads[0].genCursor(),
		LastCursor:  threads[len(threads)-1].genCursor(),
	}, nil
}

// CountThreads ...
func CountThreads(u *uexky.Uexky, selector bson.M) (int, error) {
	return u.Mongo.C(colleThread).Find(selector).Count()
}

// ThreadInput ...
type ThreadInput struct {
	Anonymous bool
	Content   string
	MainTag   string
	SubTags   *[]string
	Title     *string
}

// ParseThead ...
func (ti *ThreadInput) ParseThead(u *uexky.Uexky, user *User) (*Thread, error) {
	if !isMainTag(ti.MainTag) {
		return nil, errors.Errorf("Can't set main tag '%s'", ti.MainTag)
	}
	subTags := []string{}
	tags := []string{ti.MainTag}
	if ti.SubTags != nil && len(*ti.SubTags) != 0 {
		subTags = *ti.SubTags
	}
	for _, tag := range subTags {
		if isMainTag(tag) {
			return nil, errors.Errorf("Can't set main tag to sub tags '%s'", tag)
		}
		tags = append(tags, tag)
	}
	if ti.Content == "" {
		return nil, errors.New("Can't post an empty thread")
	}

	now := time.Now()
	thread := &Thread{
		ObjectID:   bson.NewObjectId(),
		Anonymous:  ti.Anonymous,
		UserID:     user.ID,
		CreateTime: now,
		UpdateTime: now,

		MainTag: ti.MainTag,
		SubTags: subTags,
		Tags:    tags,
		Content: ti.Content,
	}

	threadID, err := postIDGenerator.New()
	if err != nil {
		return nil, err
	}
	thread.ID = threadID

	if ti.Anonymous {
		author, err := user.AnonymousID(u, thread.ID, true)
		if err != nil {
			return nil, err
		}
		thread.Author = author
	} else {
		if user.Name == "" {
			return nil, errors.Errorf("Can't find name for user")
		}
		thread.Author = user.Name
	}

	if ti.Title != nil && *ti.Title != "" {
		thread.Title = *ti.Title
	}
	return thread, nil
}

// GetThreadsByTags ...
func GetThreadsByTags(u *uexky.Uexky, tags []string, sq *SliceQuery) (
	[]*Thread, *SliceInfo, error,
) {
	return FindThreads(u, bson.M{"tags": bson.M{"$in": tags}}, sq)
}

// FindThreadByID by id
func FindThreadByID(u *uexky.Uexky, id string) (*Thread, error) {
	return FindThread(u, bson.M{"id": id})
}

// FindThreadByPostID ...
func FindThreadByPostID(u *uexky.Uexky, id string) (*Thread, error) {
	return FindThread(u, bson.M{"posts.post_id": id})
}

func isThreadExist(u *uexky.Uexky, threadID string) (bool, error) {
	count, err := CountThreads(u, bson.M{"id": threadID})
	return count != 0, err
}

// GetReplies ...
func (t *Thread) GetReplies(u *uexky.Uexky, sq *SliceQuery) ([]*Post, *SliceInfo, error) {
	c := u.Mongo.C(collePost)
	c.EnsureIndexKey("thread_id")

	queryObj, err := sq.GenQueryByObjectID()
	if err != nil {
		return nil, nil, err
	}
	queryObj["thread_id"] = t.ID

	var posts []*Post
	if err := sq.Find(u, collePost, "_id", queryObj, &posts); err != nil {
		return nil, nil, err
	}
	if len(posts) == 0 {
		return posts, &SliceInfo{}, nil
	}
	if sq.Desc {
		ReverseSlice(posts)
	}
	si := &SliceInfo{
		FirstCursor: posts[0].ObjectID.Hex(),
		LastCursor:  posts[len(posts)-1].ObjectID.Hex(),
	}
	return posts, si, nil
}

// ReplyCount ...
func (t *Thread) ReplyCount() int {
	return len(t.Posts)
}

// return unix time of update time in millisecond(ms)
func (t *Thread) genCursor() string {
	return genTimeCursor(t.UpdateTime)
}
