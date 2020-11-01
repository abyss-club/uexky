package entity

import (
	"context"
	"fmt"
	"time"

	"gitlab.com/abyss.club/uexky/lib/config"
	"gitlab.com/abyss.club/uexky/lib/errors"
	"gitlab.com/abyss.club/uexky/lib/uid"
)

type ThreadsSearch struct {
	UserID *uid.UID
	Tags   []string
}

type ThreadRepo interface {
	CheckIfDuplicated(ctx context.Context, title *string, content string) error
	GetByID(ctx context.Context, id uid.UID) (*Thread, error)
	FindSlice(ctx context.Context, params *ThreadsSearch, query SliceQuery) (*ThreadSlice, error)

	Insert(ctx context.Context, thread *Thread) (*Thread, error)
	Update(ctx context.Context, thread *Thread) (*Thread, error)

	Replies(ctx context.Context, thread *Thread, query SliceQuery) (*PostSlice, error)
	ReplyCount(ctx context.Context, thread *Thread) (int, error)
	Catalog(ctx context.Context, thread *Thread) ([]*ThreadCatalogItem, error)
	PostAID(ctx context.Context, thread *Thread, user *User) (string, error)
}

type Thread struct {
	ID         uid.UID   `json:"id"`
	LastPostID uid.UID   `json:"-"` // for sort
	CreatedAt  time.Time `json:"createdAt"`
	Author     *Author   `json:"author"`
	Title      *string   `json:"title"`
	Content    string    `json:"content"`
	MainTag    string    `json:"main_tag"`
	SubTags    []string  `json:"sub_tags"`
	Blocked    bool      `json:"blocked"`
	Locked     bool      `json:"locked"`
}

const (
	BlockedContent       = "[此内容已被管理员屏蔽]"
	DuplicatedCheckRange = 3 * time.Minute
)

type Author struct {
	UserID    uid.UID `json:"-"`
	Guest     bool    `json:"-"`
	Anonymous bool    `json:"anonymous"`
	Author    string  `json:"author"`
}

func NewThread(user *User, input ThreadInput) (*Thread, error) {
	thread := &Thread{
		ID:        uid.NewUID(),
		CreatedAt: time.Now(),
		Author: &Author{
			UserID:    user.ID,
			Guest:     user.Role == RoleGuest,
			Anonymous: input.Anonymous,
		},
		Title:   input.Title,
		Content: input.Content,
		MainTag: input.MainTag,
		SubTags: input.SubTags,
	}
	if input.Anonymous {
		thread.Author.Author = uid.NewUID().ToBase64String()
	} else {
		if user.Name == nil {
			return nil, errors.BadParams.New("user name must be set")
		}
		thread.Author.Author = *user.Name
	}
	subTags, err := validateThreadTags(input.MainTag, input.SubTags)
	if err != nil {
		return nil, err
	}
	thread.SubTags = subTags
	return thread, nil
}

func validateThreadTags(mainTag string, subTags []string) ([]string, error) {
	mains, subs := config.SplitTags(append(subTags, mainTag)...)
	if len(mains) != 1 {
		return nil, errors.BadParams.New("must specify only one main tag")
	}
	if mains[0] != mainTag {
		return nil, errors.BadParams.Errorf("%s is not a main tag", mainTag)
	}
	return subs, nil
}

func (t *Thread) String() string {
	return fmt.Sprintf("<Thread:%v:%s>", t.ID, t.ID.ToBase64String())
}

func (t *Thread) EditTags(mainTag string, subTags []string) error {
	subTags, err := validateThreadTags(mainTag, subTags)
	if err != nil {
		return err
	}
	t.MainTag = mainTag
	t.SubTags = subTags
	return nil
}

func (t *Thread) Lock() {
	t.Locked = true
}

func (t *Thread) Block() {
	t.Blocked = true
}
