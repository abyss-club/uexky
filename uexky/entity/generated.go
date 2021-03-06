// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package entity

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"gitlab.com/abyss.club/uexky/lib/uid"
)

//  Union type of Notificatiion contents
type NotiContent interface {
	IsNotiContent()
}

//  NotiSlice object is for selecting specific 'slice' of an object to return.
// Affects the returning SliceInfo.
type NotiSlice struct {
	Notifications []*Notification `json:"notifications"`
	SliceInfo     *SliceInfo      `json:"sliceInfo"`
}

//  Input object describing a Post to be published.
type PostInput struct {
	//  ID of the replying thread's.
	ThreadID uid.UID `json:"threadId"`
	//  Should be sent anonymously or not.
	Anonymous bool `json:"anonymous"`
	//  Markdown formatted content.
	Content string `json:"content"`
	//  Set quoting PostIDs.
	QuoteIds []uid.UID `json:"quoteIds"`
}

//  Stripped version of Post object.
type PostOutline struct {
	//  post.id
	ID uid.UID `json:"id"`
	//  post.author
	Author *Author `json:"author"`
	//  Markdown formatted content.
	Content string `json:"content"`
}

//  PostSlice object is for selecting specific 'slice' of Post objects to
// return. Affects the returning SliceInfo.
type PostSlice struct {
	Posts     []*Post    `json:"posts"`
	SliceInfo *SliceInfo `json:"sliceInfo"`
}

//  Object describing contents of a quoted notification.
type QuotedNoti struct {
	//  ID of the Thread quoted.
	ThreadID uid.UID `json:"threadId"`
	//  The Post object quoted.
	QuotedPost *PostOutline `json:"quotedPost"`
	//  The Post object that made this notification.
	Post *PostOutline `json:"post"`
}

func (QuotedNoti) IsNotiContent() {}

//  Object describing contents of a replied notification.
type RepliedNoti struct {
	//  The Thread object that is replied.
	Thread *ThreadOutline `json:"thread"`
	//  How much new replies since the last read of this notification.
	NewRepliesCount int `json:"newRepliesCount"`
	//  ID of the first reply since the last read of this notification.
	FirstReplyID uid.UID `json:"firstReplyId"`
}

func (RepliedNoti) IsNotiContent() {}

//  SliceInfo objects are generated by the server.
// Can be used in consecutive queries.
type SliceInfo struct {
	FirstCursor string `json:"firstCursor"`
	LastCursor  string `json:"lastCursor"`
	//  If more results exist after lastCursor.
	HasNext bool `json:"hasNext"`
}

//  SliceQuery object is for selecting specific 'slice' of an object to return.
// Affects the returning SliceInfo.
type SliceQuery struct {
	//  Either this field or 'after' is required
	//  An empty string means slice from the beginning.
	Before *string `json:"before"`
	//  Either this field or 'before' is required.
	//  An empty string means slice to the end.
	After *string `json:"after"`
	//  Set the amount of returned items.
	Limit int `json:"limit"`
}

//  Object describing contents of a system notification.
type SystemNoti struct {
	//  Notification title.
	Title string `json:"title"`
	//  Markdown formatted notification content.
	Content string `json:"content"`
}

func (SystemNoti) IsNotiContent() {}

type Tag struct {
	//  Name of tag.
	Name string `json:"name"`
	//  The tag is a MainTag if true, SubTag otherwise.
	IsMain bool `json:"isMain"`
}

//  The ID and timestamp of post replied in the thread.
type ThreadCatalogItem struct {
	//  The ID of post.
	PostID    uid.UID   `json:"postId"`
	CreatedAt time.Time `json:"createdAt"`
}

//  Construct a new Thread.
type ThreadInput struct {
	//  Toggle anonymousness. If true, a new ID will be generated in each thread.
	Anonymous bool `json:"anonymous"`
	//  Markdown formmatted text.
	Content string `json:"content"`
	//  Required. Only one mainTag is allowed.
	MainTag string `json:"mainTag"`
	//  Optional, maximum of 4.
	SubTags []string `json:"subTags"`
	//  Optional. If not set, the title will be '无题'.
	Title *string `json:"title"`
}

//  Stripped version of Thread object.
type ThreadOutline struct {
	//  thread.id
	ID uid.UID `json:"id"`
	//  thread.title
	Title *string `json:"title"`
	//  Markdown formatted content.
	Content string   `json:"content"`
	MainTag string   `json:"mainTag"`
	SubTags []string `json:"subTags"`
}

type ThreadSlice struct {
	Threads   []*Thread  `json:"threads"`
	SliceInfo *SliceInfo `json:"sliceInfo"`
}

type NotiType string

const (
	NotiTypeSystem  NotiType = "system"
	NotiTypeReplied NotiType = "replied"
	NotiTypeQuoted  NotiType = "quoted"
)

var AllNotiType = []NotiType{
	NotiTypeSystem,
	NotiTypeReplied,
	NotiTypeQuoted,
}

func (e NotiType) IsValid() bool {
	switch e {
	case NotiTypeSystem, NotiTypeReplied, NotiTypeQuoted:
		return true
	}
	return false
}

func (e NotiType) String() string {
	return string(e)
}

func (e *NotiType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = NotiType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid NotiType", str)
	}
	return nil
}

func (e NotiType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type Role string

const (
	//  Administrators
	RoleAdmin Role = "admin"
	//  Moderators
	RoleMod Role = "mod"
	//  Normal role, default value for signed up user
	RoleNormal Role = "normal"
	//  For not signed up user
	RoleGuest Role = "guest"
	//  Banned user
	RoleBanned Role = "banned"
)

var AllRole = []Role{
	RoleAdmin,
	RoleMod,
	RoleNormal,
	RoleGuest,
	RoleBanned,
}

func (e Role) IsValid() bool {
	switch e {
	case RoleAdmin, RoleMod, RoleNormal, RoleGuest, RoleBanned:
		return true
	}
	return false
}

func (e Role) String() string {
	return string(e)
}

func (e *Role) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Role(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Role", str)
	}
	return nil
}

func (e Role) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
