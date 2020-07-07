package repo

import (
	"encoding/json"
	"time"

	"gitlab.com/abyss.club/uexky/uexky/entity"
)

type Notification struct {
	tableName struct{} `pg:"notification,,discard_unknown_columns"`

	Key       string            `pg:"key,pk"`
	SortKey   int64             `pg:"sort_key"`
	CreatedAt time.Time         `pg:"created_at"`
	UpdatedAt time.Time         `pg:"updated_at"`
	Type      string            `pg:"type,use_zero"`
	Receivers []entity.Receiver `pg:"receivers,array"`
	Content   json.RawMessage   `pg:"content"`
}

type Post struct {
	tableName struct{} `pg:"post,,discard_unknown_columns"`

	ID          int64     `pg:"id,pk"`
	CreatedAt   time.Time `pg:"created_at"`
	UpdatedAt   time.Time `pg:"updated_at"`
	ThreadID    int64     `pg:"thread_id,use_zero"`
	Anonymous   bool      `pg:"anonymous,use_zero"`
	UserID      int       `pg:"user_id,use_zero"`
	UserName    *string   `pg:"user_name"`
	AnonymousID *int64    `pg:"anonymous_id"`
	Blocked     *bool     `pg:"blocked"`
	Content     string    `pg:"content,use_zero"`
	QuotedIDs   []int64   `pg:"quoted_ids,array"`

	Thread      *Thread `pg:"post_thread_id_fkey,fk"`
	User        *User   `pg:"post_user_id_fkey,fk"`
	UserNameRel *User   `pg:"post_user_name_fkey,fk"`
}

type SchemaMigration struct {
	tableName struct{} `pg:"schema_migrations,,discard_unknown_columns"`

	Version int64 `pg:"version,pk"`
	Dirty   bool  `pg:"dirty,use_zero"`
}

type Tag struct {
	tableName struct{} `pg:"tag,,discard_unknown_columns"`

	Name      string    `pg:"name,pk"`
	CreatedAt time.Time `pg:"created_at"`
	UpdatedAt time.Time `pg:"updated_at"`
	TagType   *string   `pg:"tag_type"`
}

// TODO: tag_type index

type Thread struct {
	tableName struct{} `pg:"thread,,discard_unknown_columns"`

	ID          int64     `pg:"id,pk"`
	CreatedAt   time.Time `pg:"created_at"`
	UpdatedAt   time.Time `pg:"updated_at"`
	Anonymous   bool      `pg:"anonymous,use_zero"`
	UserID      int       `pg:"user_id,use_zero"`
	UserName    *string   `pg:"user_name"`
	AnonymousID *int64    `pg:"anonymous_id"`
	Title       *string   `pg:"title"`
	Content     string    `pg:"content,use_zero"`
	Locked      bool      `pg:"locked,use_zero"`
	Blocked     bool      `pg:"blocked,use_zero"`
	LastPostID  int64     `pg:"last_post_id,use_zero"`
	Tags        []string  `pg:"tags,array"`

	User        *User `pg:"thread_user_id_fkey,fk"`
	UserNameRel *User `pg:"thread_user_name_fkey,fk"`
}

type User struct {
	tableName struct{} `pg:"user,,discard_unknown_columns"`

	ID           int       `pg:"id,pk"`
	CreatedAt    time.Time `pg:"created_at"`
	UpdatedAt    time.Time `pg:"updated_at"`
	Email        string    `pg:"email,use_zero"`
	Name         *string   `pg:"name"`
	Role         string    `pg:"role,use_zero"`
	LastReadNoti int64     `pg:"last_read_noti,use_zero"`
	Tags         []string  `pg:"tags,array"`
}
