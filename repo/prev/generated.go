//nolint
//lint:file-ignore U1000 ignore unused code, it's generated
package repo

import (
	"time"
)

var Columns = struct {
	AnonymouId struct {
		ID, CreatedAt, UpdatedAt, ThreadID, UserID, AnonymousID string
	}
	Config struct {
		ID, RateLimit, RateCost string
	}
	Counter struct {
		Name, Count string
	}
	Notification struct {
		ID, Key, CreatedAt, UpdatedAt, Type, SendTo, SendToGroup, Content string

		SendToRel string
	}
	Pgmigration struct {
		ID, Name, RunOn string
	}
	Post struct {
		ID, CreatedAt, UpdatedAt, ThreadID, Anonymous, UserID, UserName, AnonymousID, Blocked, Content string

		Thread, User, UserNameRel string
	}
	PostsQuote struct {
		ID, QuoterID, QuotedID string

		Quoted, Quoter string
	}
	Tag struct {
		Name, IsMain, CreatedAt, UpdatedAt string
	}
	TagsMainTag struct {
		ID, CreatedAt, UpdatedAt, Name, BelongsTo string

		BelongsToRel, NameRel string
	}
	Thread struct {
		ID, CreatedAt, UpdatedAt, Anonymous, UserID, UserName, AnonymousID, Title, Content, Locked, Blocked, LastPostID string

		User, UserNameRel string
	}
	ThreadsTag struct {
		ID, CreatedAt, UpdatedAt, ThreadID, TagName string

		TagNameRel, Thread string
	}
	User struct {
		ID, CreatedAt, UpdatedAt, Email, Name, Role, LastReadSystemNoti, LastReadRepliedNoti, LastReadQuotedNoti string
	}
	UsersTag struct {
		ID, UserID, TagName string

		TagNameRel, User string
	}
}{
	AnonymouId: struct {
		ID, CreatedAt, UpdatedAt, ThreadID, UserID, AnonymousID string
	}{
		ID:          "id",
		CreatedAt:   "created_at",
		UpdatedAt:   "updated_at",
		ThreadID:    "thread_id",
		UserID:      "user_id",
		AnonymousID: "anonymous_id",
	},
	Config: struct {
		ID, RateLimit, RateCost string
	}{
		ID:        "id",
		RateLimit: "rate_limit",
		RateCost:  "rate_cost",
	},
	Counter: struct {
		Name, Count string
	}{
		Name:  "name",
		Count: "count",
	},
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
	Pgmigration: struct {
		ID, Name, RunOn string
	}{
		ID:    "id",
		Name:  "name",
		RunOn: "run_on",
	},
	Post: struct {
		ID, CreatedAt, UpdatedAt, ThreadID, Anonymous, UserID, UserName, AnonymousID, Blocked, Content string

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

		Thread:      "Thread",
		User:        "User",
		UserNameRel: "UserNameRel",
	},
	PostsQuote: struct {
		ID, QuoterID, QuotedID string

		Quoted, Quoter string
	}{
		ID:       "id",
		QuoterID: "quoter_id",
		QuotedID: "quoted_id",

		Quoted: "Quoted",
		Quoter: "Quoter",
	},
	Tag: struct {
		Name, IsMain, CreatedAt, UpdatedAt string
	}{
		Name:      "name",
		IsMain:    "is_main",
		CreatedAt: "created_at",
		UpdatedAt: "updated_at",
	},
	TagsMainTag: struct {
		ID, CreatedAt, UpdatedAt, Name, BelongsTo string

		BelongsToRel, NameRel string
	}{
		ID:        "id",
		CreatedAt: "created_at",
		UpdatedAt: "updated_at",
		Name:      "name",
		BelongsTo: "belongs_to",

		BelongsToRel: "BelongsToRel",
		NameRel:      "NameRel",
	},
	Thread: struct {
		ID, CreatedAt, UpdatedAt, Anonymous, UserID, UserName, AnonymousID, Title, Content, Locked, Blocked, LastPostID string

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

		User:        "User",
		UserNameRel: "UserNameRel",
	},
	ThreadsTag: struct {
		ID, CreatedAt, UpdatedAt, ThreadID, TagName string

		TagNameRel, Thread string
	}{
		ID:        "id",
		CreatedAt: "created_at",
		UpdatedAt: "updated_at",
		ThreadID:  "thread_id",
		TagName:   "tag_name",

		TagNameRel: "TagNameRel",
		Thread:     "Thread",
	},
	User: struct {
		ID, CreatedAt, UpdatedAt, Email, Name, Role, LastReadSystemNoti, LastReadRepliedNoti, LastReadQuotedNoti string
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
	},
	UsersTag: struct {
		ID, UserID, TagName string

		TagNameRel, User string
	}{
		ID:      "id",
		UserID:  "user_id",
		TagName: "tag_name",

		TagNameRel: "TagNameRel",
		User:       "User",
	},
}

var Tables = struct {
	AnonymouId struct {
		Name string
	}
	Config struct {
		Name string
	}
	Counter struct {
		Name string
	}
	Notification struct {
		Name string
	}
	Pgmigration struct {
		Name string
	}
	Post struct {
		Name string
	}
	PostsQuote struct {
		Name string
	}
	Tag struct {
		Name string
	}
	TagsMainTag struct {
		Name string
	}
	Thread struct {
		Name string
	}
	ThreadsTag struct {
		Name string
	}
	User struct {
		Name string
	}
	UsersTag struct {
		Name string
	}
}{
	AnonymouId: struct {
		Name string
	}{
		Name: "anonymous_id",
	},
	Config: struct {
		Name string
	}{
		Name: "config",
	},
	Counter: struct {
		Name string
	}{
		Name: "counter",
	},
	Notification: struct {
		Name string
	}{
		Name: "notification",
	},
	Pgmigration: struct {
		Name string
	}{
		Name: "pgmigrations",
	},
	Post: struct {
		Name string
	}{
		Name: "post",
	},
	PostsQuote: struct {
		Name string
	}{
		Name: "posts_quotes",
	},
	Tag: struct {
		Name string
	}{
		Name: "tag",
	},
	TagsMainTag: struct {
		Name string
	}{
		Name: "tags_main_tags",
	},
	Thread: struct {
		Name string
	}{
		Name: "thread",
	},
	ThreadsTag: struct {
		Name string
	}{
		Name: "threads_tags",
	},
	User: struct {
		Name string
	}{
		Name: "user",
	},
	UsersTag: struct {
		Name string
	}{
		Name: "users_tags",
	},
}

type AnonymouId struct {
	tableName struct{} `pg:"anonymous_id,,discard_unknown_columns"`

	ID          int       `pg:"id,pk"`
	CreatedAt   time.Time `pg:"created_at,use_zero"`
	UpdatedAt   time.Time `pg:"updated_at,use_zero"`
	ThreadID    int64     `pg:"thread_id,use_zero"`
	UserID      int       `pg:"user_id,use_zero"`
	AnonymousID int64     `pg:"anonymous_id,use_zero"`
}

type Config struct {
	tableName struct{} `pg:"config,,discard_unknown_columns"`

	ID        int                    `pg:"id,pk"`
	RateLimit map[string]interface{} `pg:"rate_limit,use_zero"`
	RateCost  map[string]interface{} `pg:"rate_cost,use_zero"`
}

type Counter struct {
	tableName struct{} `pg:"counter,,discard_unknown_columns"`

	Name  string `pg:"name,pk"`
	Count *int   `pg:"count"`
}

type Notification struct {
	tableName struct{} `pg:"notification,,discard_unknown_columns"`

	ID          int                    `pg:"id,pk"`
	Key         *string                `pg:"key"`
	CreatedAt   time.Time              `pg:"created_at,use_zero"`
	UpdatedAt   time.Time              `pg:"updated_at,use_zero"`
	Type        string                 `pg:"type,use_zero"`
	SendTo      *int                   `pg:"send_to"`
	SendToGroup *string                `pg:"send_to_group"`
	Content     map[string]interface{} `pg:"content"`

	SendToRel *User `pg:"fk:send_to"`
}

type Pgmigration struct {
	tableName struct{} `pg:"pgmigrations,,discard_unknown_columns"`

	ID    int       `pg:"id,pk"`
	Name  string    `pg:"name,use_zero"`
	RunOn time.Time `pg:"run_on,use_zero"`
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

	Thread      *Thread `pg:"fk:thread_id"`
	User        *User   `pg:"fk:user_id"`
	UserNameRel *User   `pg:"fk:user_name"`
}

type PostsQuote struct {
	tableName struct{} `pg:"posts_quotes,,discard_unknown_columns"`

	ID       int   `pg:"id,pk"`
	QuoterID int64 `pg:"quoter_id,use_zero"`
	QuotedID int64 `pg:"quoted_id,use_zero"`

	Quoted *Post `pg:"fk:quoted_id"`
	Quoter *Post `pg:"fk:quoter_id"`
}

type Tag struct {
	tableName struct{} `pg:"tag,,discard_unknown_columns"`

	Name      string    `pg:"name,pk"`
	IsMain    bool      `pg:"is_main,use_zero"`
	CreatedAt time.Time `pg:"created_at,use_zero"`
	UpdatedAt time.Time `pg:"updated_at,use_zero"`
}

type TagsMainTag struct {
	tableName struct{} `pg:"tags_main_tags,,discard_unknown_columns"`

	ID        int       `pg:"id,pk"`
	CreatedAt time.Time `pg:"created_at,use_zero"`
	UpdatedAt time.Time `pg:"updated_at,use_zero"`
	Name      string    `pg:"name,use_zero"`
	BelongsTo string    `pg:"belongs_to,use_zero"`

	BelongsToRel *Tag `pg:"fk:belongs_to"`
	NameRel      *Tag `pg:"fk:name"`
}

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

	User        *User `pg:"fk:user_id"`
	UserNameRel *User `pg:"fk:user_name"`
}

type ThreadsTag struct {
	tableName struct{} `pg:"threads_tags,,discard_unknown_columns"`

	ID        int       `pg:"id,pk"`
	CreatedAt time.Time `pg:"created_at,use_zero"`
	UpdatedAt time.Time `pg:"updated_at,use_zero"`
	ThreadID  int64     `pg:"thread_id,use_zero"`
	TagName   string    `pg:"tag_name,use_zero"`

	TagNameRel *Tag    `pg:"fk:tag_name"`
	Thread     *Thread `pg:"fk:thread_id"`
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
}

type UsersTag struct {
	tableName struct{} `pg:"users_tags,,discard_unknown_columns"`

	ID      int    `pg:"id,pk"`
	UserID  int    `pg:"user_id,use_zero"`
	TagName string `pg:"tag_name,use_zero"`

	TagNameRel *Tag  `pg:"fk:tag_name"`
	User       *User `pg:"fk:user_id"`
}
