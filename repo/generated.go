//nolint
//lint:file-ignore U1000 ignore unused code, it's generated
package repo

import (
	"encoding/json"
	"time"
)

var Columns = struct {
	Notification struct {
		ID, Key, CreatedAt, UpdatedAt, Type, SendTo, SendToGroup, Content string

		SendToRel string
	}
	Post struct {
		ID, CreatedAt, UpdatedAt, ThreadID, Anonymous, UserID, UserName, AnonymousID, Blocked, Content, QuotedIDs string

		Thread, User, UserNameRel string
	}
	SchemaMigration struct {
		Version, Dirty string
	}
	Tag struct {
		Name, CreatedAt, UpdatedAt, TagType string
	}
	Thread struct {
		ID, CreatedAt, UpdatedAt, Anonymous, UserID, UserName, AnonymousID, Title, Content, Locked, Blocked, LastPostID, Tags string

		User, UserNameRel string
	}
	User struct {
		ID, CreatedAt, UpdatedAt, Email, Name, Role, LastReadSystemNoti, LastReadRepliedNoti, LastReadQuotedNoti, Tags string
	}
}{
	Notification: struct {
		ID, Key, CreatedAt, UpdatedAt, Type, SendTo, SendToGroup, Content string

		SendToRel string
	}{
		ID:          "id",
		Key:         "key",
		CreatedAt:   "created_at",
		UpdatedAt:   "updated_at",
		Type:        "type",
		SendTo:      "send_to",
		SendToGroup: "send_to_group",
		Content:     "content",

		SendToRel: "SendToRel",
	},
	Post: struct {
		ID, CreatedAt, UpdatedAt, ThreadID, Anonymous, UserID, UserName, AnonymousID, Blocked, Content, QuotedIDs string

		Thread, User, UserNameRel string
	}{
		ID:          "id",
		CreatedAt:   "created_at",
		UpdatedAt:   "updated_at",
		ThreadID:    "thread_id",
		Anonymous:   "anonymous",
		UserID:      "user_id",
		UserName:    "user_name",
		AnonymousID: "anonymous_id",
		Blocked:     "blocked",
		Content:     "content",
		QuotedIDs:   "quoted_ids",

		Thread:      "Thread",
		User:        "User",
		UserNameRel: "UserNameRel",
	},
	SchemaMigration: struct {
		Version, Dirty string
	}{
		Version: "version",
		Dirty:   "dirty",
	},
	Tag: struct {
		Name, CreatedAt, UpdatedAt, TagType string
	}{
		Name:      "name",
		CreatedAt: "created_at",
		UpdatedAt: "updated_at",
		TagType:   "tag_type",
	},
	Thread: struct {
		ID, CreatedAt, UpdatedAt, Anonymous, UserID, UserName, AnonymousID, Title, Content, Locked, Blocked, LastPostID, Tags string

		User, UserNameRel string
	}{
		ID:          "id",
		CreatedAt:   "created_at",
		UpdatedAt:   "updated_at",
		Anonymous:   "anonymous",
		UserID:      "user_id",
		UserName:    "user_name",
		AnonymousID: "anonymous_id",
		Title:       "title",
		Content:     "content",
		Locked:      "locked",
		Blocked:     "blocked",
		LastPostID:  "last_post_id",
		Tags:        "tags",

		User:        "User",
		UserNameRel: "UserNameRel",
	},
	User: struct {
		ID, CreatedAt, UpdatedAt, Email, Name, Role, LastReadSystemNoti, LastReadRepliedNoti, LastReadQuotedNoti, Tags string
	}{
		ID:                  "id",
		CreatedAt:           "created_at",
		UpdatedAt:           "updated_at",
		Email:               "email",
		Name:                "name",
		Role:                "role",
		LastReadSystemNoti:  "last_read_system_noti",
		LastReadRepliedNoti: "last_read_replied_noti",
		LastReadQuotedNoti:  "last_read_quoted_noti",
		Tags:                "tags",
	},
}

var Tables = struct {
	Notification struct {
		Name string
	}
	Post struct {
		Name string
	}
	SchemaMigration struct {
		Name string
	}
	Tag struct {
		Name string
	}
	Thread struct {
		Name string
	}
	User struct {
		Name string
	}
}{
	Notification: struct {
		Name string
	}{
		Name: "notification",
	},
	Post: struct {
		Name string
	}{
		Name: "post",
	},
	SchemaMigration: struct {
		Name string
	}{
		Name: "schema_migrations",
	},
	Tag: struct {
		Name string
	}{
		Name: "tag",
	},
	Thread: struct {
		Name string
	}{
		Name: "thread",
	},
	User: struct {
		Name string
	}{
		Name: "user",
	},
}

type Notification struct {
	tableName struct{} `pg:"notification,,discard_unknown_columns"`

	ID          int             `pg:"id,pk"`
	Key         *string         `pg:"key"`
	CreatedAt   time.Time       `pg:"created_at,use_zero"`
	UpdatedAt   time.Time       `pg:"updated_at,use_zero"`
	Type        string          `pg:"type,use_zero"`
	SendTo      *int            `pg:"send_to"`
	SendToGroup *string         `pg:"send_to_group"`
	Content     json.RawMessage `pg:"content"`

	SendToRel *User `pg:"fk:send_to"`
}

type Post struct {
	tableName struct{} `pg:"post,,discard_unknown_columns"`

	ID          int64     `pg:"id,pk"`
	CreatedAt   time.Time `pg:"created_at,use_zero"`
	UpdatedAt   time.Time `pg:"updated_at,use_zero"`
	ThreadID    int64     `pg:"thread_id,use_zero"`
	Anonymous   bool      `pg:"anonymous,use_zero"`
	UserID      int       `pg:"user_id,use_zero"`
	UserName    *string   `pg:"user_name"`
	AnonymousID *int64    `pg:"anonymous_id"`
	Blocked     *bool     `pg:"blocked"`
	Content     string    `pg:"content,use_zero"`
	QuotedIDs   []int64   `pg:"quoted_ids,array"`

	Thread      *Thread `pg:"fk:thread_id"`
	User        *User   `pg:"fk:user_id"`
	UserNameRel *User   `pg:"fk:user_name"`
}

type SchemaMigration struct {
	tableName struct{} `pg:"schema_migrations,,discard_unknown_columns"`

	Version int64 `pg:"version,pk"`
	Dirty   bool  `pg:"dirty,use_zero"`
}

type Tag struct {
	tableName struct{} `pg:"tag,,discard_unknown_columns"`

	Name      string    `pg:"name,pk"`
	CreatedAt time.Time `pg:"created_at,use_zero"`
	UpdatedAt time.Time `pg:"updated_at,use_zero"`
	TagType   *string   `pg:"tag_type"`
}

// TODO: tag_type index

type Thread struct {
	tableName struct{} `pg:"thread,,discard_unknown_columns"`

	ID          int64     `pg:"id,pk"`
	CreatedAt   time.Time `pg:"created_at,use_zero"`
	UpdatedAt   time.Time `pg:"updated_at,use_zero"`
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

	User        *User `pg:"fk:user_id"`
	UserNameRel *User `pg:"fk:user_name"`
}

type User struct {
	tableName struct{} `pg:"user,,discard_unknown_columns"`

	ID                  int       `pg:"id,pk"`
	CreatedAt           time.Time `pg:"created_at,use_zero"`
	UpdatedAt           time.Time `pg:"updated_at,use_zero"`
	Email               string    `pg:"email,use_zero"`
	Name                *string   `pg:"name"`
	Role                string    `pg:"role,use_zero"`
	LastReadSystemNoti  int       `pg:"last_read_system_noti,use_zero"`
	LastReadRepliedNoti int       `pg:"last_read_replied_noti,use_zero"`
	LastReadQuotedNoti  int       `pg:"last_read_quoted_noti,use_zero"`
	Tags                []string  `pg:"tags,array"`
}
