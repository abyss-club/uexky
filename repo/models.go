package repo

import (
	"time"

	"gitlab.com/abyss.club/uexky/lib/uid"
	"gitlab.com/abyss.club/uexky/uexky/entity"
)

type User struct {
	//nolint: structcheck, unused
	tableName struct{} `pg:"user,,discard_unknown_columns"`

	ID           uid.UID     `pg:"id,pk" json:"id"`
	CreatedAt    time.Time   `pg:"created_at" json:"-"`
	UpdatedAt    time.Time   `pg:"updated_at" json:"-"`
	Email        *string     `pg:"email,use_zero" json:"-"`
	Name         *string     `pg:"name" json:"-"`
	Role         entity.Role `pg:"role,use_zero" json:"role"`
	LastReadNoti uid.UID     `pg:"last_read_noti,use_zero" json:"-"`
	Tags         []string    `pg:"tags,array" json:"tags"`
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
	Title      *string   `pg:"title"`
	Content    string    `pg:"content,use_zero"`
	Locked     bool      `pg:"locked,use_zero"`
	Blocked    bool      `pg:"blocked,use_zero"`
	Tags       []string  `pg:"tags,array"`
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
	Blocked   *bool     `pg:"blocked"`
	Content   string    `pg:"content,use_zero"`
	QuotedIDs []uid.UID `pg:"quoted_ids,array"`
}

type Tag struct {
	//nolint: structcheck, unused
	tableName struct{} `pg:"tag,,discard_unknown_columns"`

	Name      string    `pg:"name,pk"`
	CreatedAt time.Time `pg:"created_at"`
	UpdatedAt time.Time `pg:"updated_at"`
	TagType   *string   `pg:"type"`
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
