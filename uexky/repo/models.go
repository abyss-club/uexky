package repo

import (
	"time"

	"gitlab.com/abyss.club/uexky/lib/config"
	"gitlab.com/abyss.club/uexky/lib/uerr"
	"gitlab.com/abyss.club/uexky/lib/uid"
	"gitlab.com/abyss.club/uexky/uexky/entity"
)

type User struct {
	//nolint: structcheck, unused
	tableName struct{} `pg:"user,,discard_unknown_columns"`

	ID           uid.UID     `pg:"id,pk" json:"id"`
	CreatedAt    time.Time   `pg:"created_at" json:"created_at"`
	UpdatedAt    time.Time   `pg:"updated_at" json:"updated_at"`
	Email        *string     `pg:"email,use_zero" json:"-"`
	Name         *string     `pg:"name,user_zero" json:"-"`
	Role         entity.Role `pg:"role,use_zero" json:"role"`
	LastReadNoti uid.UID     `pg:"last_read_noti,use_zero" json:"-"`
	Tags         []string    `pg:"tags,array" json:"tags"`
}

func NewUserFromEntity(user *entity.User) *User {
	// unmapped: CreatedAt, UpdatedAt
	return &User{
		ID:           user.ID,
		Email:        user.Email,
		Name:         user.Name,
		Role:         user.Role,
		LastReadNoti: user.LastReadNoti,
		Tags:         user.Tags,
	}
}

func (u *User) ToEntity() *entity.User {
	user := &entity.User{
		ID:           u.ID,
		Email:        u.Email,
		Name:         u.Name,
		Role:         u.Role,
		Tags:         u.Tags,
		LastReadNoti: u.LastReadNoti,
	}
	// TODO: should in service level?
	if len(user.Tags) == 0 {
		user.Tags = config.GetMainTags()
	}
	if user.Role == "" {
		user.Role = entity.RoleNormal
	}
	return user
}

type Thread struct {
	//nolint: structcheck, unused
	tableName struct{} `pg:"thread,,discard_unknown_columns"`

	ID         uid.UID   `pg:"id,pk"`
	LastPostID uid.UID   `pg:"last_post_id,use_zero"`
	CreatedAt  time.Time `pg:"created_at"`
	UpdatedAt  time.Time `pg:"updated_at"`
	UserID     uid.UID   `pg:"user_id,use_zero"`
	Anonymous  bool      `pg:"anonymous,use_zero"`
	Guest      bool      `pg:"guest,use_zero"`
	Author     string    `pg:"author"`
	Title      *string   `pg:"title,use_zero"`
	Content    string    `pg:"content,use_zero"`
	Locked     bool      `pg:"locked,use_zero"`
	Blocked    bool      `pg:"blocked,use_zero"`
	Tags       []string  `pg:"tags,array"`
}

func NewThreadFromEntity(thread *entity.Thread) *Thread {
	// unmapped: LastPostID, UpdatedAt, Content
	t := &Thread{
		ID:        thread.ID,
		CreatedAt: thread.CreatedAt,
		UserID:    thread.Author.UserID,
		Anonymous: thread.Author.Anonymous,
		Guest:     thread.Author.Guest,
		Author:    thread.Author.Author,
		Title:     thread.Title,
		Locked:    thread.Locked,
		Blocked:   thread.Blocked,
		Tags:      []string{thread.MainTag},
	}
	t.Tags = append(t.Tags, thread.SubTags...)
	if !thread.Blocked {
		t.Content = thread.Content
	}
	return t
}

func (t *Thread) ToEntity() *entity.Thread {
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
	}
	if thread.Blocked {
		thread.Content = entity.BlockedContent
	}
	return thread
}

type Post struct {
	//nolint: structcheck, unused
	tableName struct{} `pg:"post,,discard_unknown_columns"`

	ID        uid.UID   `pg:"id,pk"`
	CreatedAt time.Time `pg:"created_at"`
	UpdatedAt time.Time `pg:"updated_at"`
	ThreadID  uid.UID   `pg:"thread_id,use_zero"`
	UserID    uid.UID   `pg:"user_id,use_zero"`
	Anonymous bool      `pg:"anonymous,use_zero"`
	Guest     bool      `pg:"guest,use_zero"`
	Author    string    `pg:"author"`
	Blocked   bool      `pg:"blocked"`
	Content   string    `pg:"content,use_zero"`
	QuotedIDs []uid.UID `pg:"quoted_ids,array"`
}

func NewPostFromEntity(post *entity.Post) *Post {
	// unmapped: UpdatedAt, Content
	p := &Post{
		ID:        post.ID,
		CreatedAt: post.CreatedAt,
		ThreadID:  post.ThreadID,
		UserID:    post.Author.UserID,
		Anonymous: post.Author.Anonymous,
		Guest:     post.Author.Guest,
		Author:    post.Author.Author,
		Blocked:   post.Blocked,
		QuotedIDs: post.QuoteIDs,
	}
	if !p.Blocked {
		p.Content = post.Content
	}
	return p
}

func (p *Post) ToEntity() *entity.Post {
	post := &entity.Post{
		ID:        p.ID,
		ThreadID:  p.ThreadID,
		CreatedAt: p.CreatedAt,
		Author: &entity.Author{
			UserID:    p.UserID,
			Guest:     p.Guest,
			Anonymous: p.Anonymous,
			Author:    p.Author,
		},
		QuoteIDs: p.QuotedIDs,
		Content:  p.Content,
		Blocked:  p.Blocked,
	}
	if post.Blocked {
		post.Content = entity.BlockedContent
	}
	return post
}

type Tag struct {
	//nolint: structcheck, unused
	tableName struct{} `pg:"tag,,discard_unknown_columns"`

	Name      string    `pg:"name,pk"`
	CreatedAt time.Time `pg:"created_at"`
	UpdatedAt time.Time `pg:"updated_at"`
	TagType   *string   `pg:"type,use_zero"`
}

type Notification struct {
	//nolint: structcheck, unused
	tableName struct{} `pg:"notification,,discard_unknown_columns"`

	Key       string                 `pg:"key,pk"`
	SortKey   uid.UID                `pg:"sort_key"`
	CreatedAt time.Time              `pg:"created_at"`
	UpdatedAt time.Time              `pg:"updated_at"`
	Type      entity.NotiType        `pg:"type,use_zero"`
	Receivers []entity.Receiver      `pg:"receivers,array"`
	Content   map[string]interface{} `pg:"content,json_use_number"`
}

type NotificationQuery struct {
	Notification `pg:",inherit"`

	HasRead bool `pg:"has_read"`
}

func NewNotificaionFromEntity(notification *entity.Notification) (*Notification, error) {
	// unmapped: CreatedAt
	n := &Notification{
		Key:       notification.Key,
		SortKey:   notification.SortKey,
		UpdatedAt: notification.EventTime,
		Type:      notification.Type,
		Receivers: notification.Receivers,
	}
	content, err := notification.EncodeContent()
	if err != nil {
		return nil, err
	}
	n.Content = content
	return n, nil
}

func (n *NotificationQuery) ToEntity() *entity.Notification {
	notification := &entity.Notification{
		Type:      n.Type,
		EventTime: n.CreatedAt,
		HasRead:   n.HasRead,
		Key:       n.Key,
		SortKey:   n.SortKey,
		Receivers: n.Receivers,
	}
	if err := notification.DecodeContent(n.Content); err != nil {
		panic(uerr.Errorf(uerr.InternalError, "read notification error: %w", err))
	}
	return notification
}
