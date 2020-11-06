package main

import (
	"time"

	"gitlab.com/abyss.club/uexky/lib/uid"
	"gitlab.com/abyss.club/uexky/uexky/entity"
)

type Notification struct {
	//nolint: structcheck, unused
	tableName struct{} `pg:"notification,,discard_unknown_columns"`

	ID          int                    `pg:"id,pk"`
	Key         *string                `pg:"key"`
	CreatedAt   time.Time              `pg:"created_at,use_zero"`
	UpdatedAt   time.Time              `pg:"updated_at,use_zero"`
	Type        entity.NotiType        `pg:"type,use_zero"`
	SendTo      *int                   `pg:"send_to"`
	SendToGroup *string                `pg:"send_to_group"`
	Content     map[string]interface{} `pg:"content"`
}

type Post struct {
	//nolint: structcheck, unused
	tableName struct{} `pg:"post,,discard_unknown_columns"`

	ID          uid.UID   `pg:"id,pk"`
	CreatedAt   time.Time `pg:"created_at,use_zero"`
	UpdatedAt   time.Time `pg:"updated_at,use_zero"`
	ThreadID    uid.UID   `pg:"thread_id,use_zero"`
	Anonymous   bool      `pg:"anonymous,use_zero"`
	UserID      int       `pg:"user_id,use_zero"`
	UserName    *string   `pg:"user_name"`
	AnonymousID *uid.UID  `pg:"anonymous_id"`
	Blocked     *bool     `pg:"blocked"`
	Content     string    `pg:"content,use_zero"`
}

type PostsQuote struct {
	//nolint: structcheck, unused
	tableName struct{} `pg:"posts_quotes,,discard_unknown_columns"`

	ID       int   `pg:"id,pk"`
	QuoterID int64 `pg:"quoter_id,use_zero"`
	QuotedID int64 `pg:"quoted_id,use_zero"`
}

type Tag struct {
	//nolint: structcheck, unused
	tableName struct{} `pg:"tag,,discard_unknown_columns"`

	Name      string    `pg:"name,pk"`
	IsMain    bool      `pg:"is_main,use_zero"`
	CreatedAt time.Time `pg:"created_at,use_zero"`
	UpdatedAt time.Time `pg:"updated_at,use_zero"`
}

type Thread struct {
	//nolint: structcheck, unused
	tableName struct{} `pg:"thread,,discard_unknown_columns"`

	ID          uid.UID   `pg:"id,pk"`
	CreatedAt   time.Time `pg:"created_at,use_zero"`
	UpdatedAt   time.Time `pg:"updated_at,use_zero"`
	Anonymous   bool      `pg:"anonymous,use_zero"`
	UserID      int       `pg:"user_id,use_zero"`
	UserName    *string   `pg:"user_name"`
	AnonymousID *uid.UID  `pg:"anonymous_id"`
	Title       *string   `pg:"title"`
	Content     string    `pg:"content,use_zero"`
	Locked      bool      `pg:"locked,use_zero"`
	Blocked     bool      `pg:"blocked,use_zero"`
	LastPostID  uid.UID   `pg:"last_post_id,use_zero"`
}

type ThreadsTag struct {
	//nolint: structcheck, unused
	tableName struct{} `pg:"threads_tags,,discard_unknown_columns"`

	ID        int       `pg:"id,pk"`
	CreatedAt time.Time `pg:"created_at,use_zero"`
	UpdatedAt time.Time `pg:"updated_at,use_zero"`
	ThreadID  int64     `pg:"thread_id,use_zero"`
	TagName   string    `pg:"tag_name,use_zero"`
}

type User struct {
	//nolint: structcheck, unused
	tableName struct{} `pg:"user,,discard_unknown_columns"`

	ID                  int         `pg:"id,pk"`
	CreatedAt           time.Time   `pg:"created_at,use_zero"`
	UpdatedAt           time.Time   `pg:"updated_at,use_zero"`
	Email               string      `pg:"email,use_zero"`
	Name                *string     `pg:"name"`
	Role                entity.Role `pg:"role,use_zero"`
	LastReadSystemNoti  int         `pg:"last_read_system_noti,use_zero"`
	LastReadRepliedNoti int         `pg:"last_read_replied_noti,use_zero"`
	LastReadQuotedNoti  int         `pg:"last_read_quoted_noti,use_zero"`
}

type UsersTag struct {
	//nolint: structcheck, unused
	tableName struct{} `pg:"users_tags,,discard_unknown_columns"`

	ID      int    `pg:"id,pk"`
	UserID  int    `pg:"user_id,use_zero"`
	TagName string `pg:"tag_name,use_zero"`
}
