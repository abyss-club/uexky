// forum notificatoin: systemNoti repliedNoti quotedNoti

package entity

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky/lib/uid"
)

type NotiService struct {
	Repo NotiRepo
}

func (n *NotiService) GetUnreadNotiCount(ctx context.Context, user *User) (*UnreadNotiCount, error) {
	return n.Repo.GetUserUnreadCount(ctx, user)
}

func (n *NotiService) GetNotification(
	ctx context.Context, user *User, typeArg NotiType, query SliceQuery,
) (*NotiSlice, error) {
	slice, err := n.Repo.GetNotiSlice(ctx, &NotiSearch{UserID: user.ID, Type: typeArg}, query)
	if err != nil {
		return nil, err
	}
	read := user.LastReadNoti.Get(typeArg)
	max := slice.GetMaxID(typeArg)
	if max >= read {
		if err := n.Repo.UpdateReadID(ctx, user.ID, typeArg, max); err != nil {
			return nil, err
		}
		user.LastReadNoti.Set(typeArg, max)
	}
	return slice, nil
}

func ParseNotiType(s string) (NotiType, error) {
	t := NotiType(s)
	if !t.IsValid() {
		return NotiType(""), errors.Errorf("invalid noti type: %s", s)
	}
	return t, nil
}

type LastReadNoti struct {
	SystemNoti  int
	RepliedNoti int
	QuotedNoti  int
}

func (lrn *LastReadNoti) Get(t NotiType) int {
	switch t {
	case NotiTypeQuoted:
		return lrn.QuotedNoti
	case NotiTypeReplied:
		return lrn.RepliedNoti
	case NotiTypeSystem:
		return lrn.SystemNoti
	default:
		return 0
	}
}

func (lrn *LastReadNoti) Set(t NotiType, id int) {
	switch t {
	case NotiTypeQuoted:
		lrn.QuotedNoti = id
	case NotiTypeReplied:
		lrn.RepliedNoti = id
	case NotiTypeSystem:
		lrn.SystemNoti = id
	}
}

func (ns *NotiSlice) GetMaxID(t NotiType) int {
	max := 0
	switch t {
	case NotiTypeQuoted:
		for _, n := range ns.Quoted {
			if n.ID > max {
				max = n.ID
			}
		}
	case NotiTypeReplied:
		for _, n := range ns.Replied {
			if n.ID > max {
				max = n.ID
			}
		}
	case NotiTypeSystem:
		for _, n := range ns.System {
			if n.ID > max {
				max = n.ID
			}
		}
	}
	return max
}

type SendGroup string

const AllUser SendGroup = "all_user"

type SendTo struct {
	UserID    *int       `json:"send_to"`
	SendGroup *SendGroup `json:"send_to_group"`
}

func SendToUser(userID int) SendTo {
	return SendTo{UserID: &userID}
}

func SendToGourp(group SendGroup) SendTo {
	return SendTo{SendGroup: &group}
}

type SystemNotiContent struct {
	Title       string `json:"title"`
	ContentText string `json:"content"`
}

type SystemNoti struct {
	ID        int               `json:"id"`
	Type      NotiType          `json:"type"`
	EventTime time.Time         `json:"eventTime"`
	Content   SystemNotiContent `json:"content"`

	SendTo SendTo `json:"-"`
}

func (n *NotiService) NewSystemNoti(ctx context.Context, title, content string, sendTo SendTo) error {
	noti := &SystemNoti{
		Type:      NotiTypeSystem,
		EventTime: time.Now(),
		Content: SystemNotiContent{
			Title:       title,
			ContentText: content,
		},
		SendTo: sendTo,
	}
	return n.Repo.InsertNoti(ctx, NotiInsert{System: noti})
}

func (n *SystemNoti) HasRead(user *User) bool {
	return user.LastReadNoti.SystemNoti >= n.ID
}

func (n *SystemNoti) Title() string {
	return n.Content.Title
}

func (n *SystemNoti) ContentText() string {
	return n.Content.ContentText
}

type RepliedNotiContent struct {
	ThreadID string `json:"threadId"` // int64.toString()

	Thread *Thread `json:"-"`
}

type RepliedNoti struct {
	ID        int                `json:"id"`
	Type      NotiType           `json:"type"`
	EventTime time.Time          `json:"eventTime"`
	Content   RepliedNotiContent `json:"content"`

	SendTo SendTo `json:"-"`
}

func (n *NotiService) NewRepliedNoti(ctx context.Context, thread *Thread, reply *Post) error {
	noti := &RepliedNoti{
		Type:      NotiTypeReplied,
		EventTime: time.Now(),
		Content: RepliedNotiContent{
			ThreadID: thread.ID.ToBase64String(),
		},
		SendTo: SendToUser(thread.AuthorObj.UserID),
	}
	return n.Repo.InsertNoti(ctx, NotiInsert{Replied: noti})
}

func (n *RepliedNoti) HasRead(user *User) bool {
	return user.LastReadNoti.RepliedNoti >= n.ID
}

func (n *RepliedNoti) ThreadID() (uid.UID, error) {
	return uid.ParseUID(n.Content.ThreadID)
}

type QuotedNotiContent struct {
	ThreadID string `json:"threadId"` // int64.toString()
	QuotedID string `json:"quotedId"` // int64.toString()
	PostID   string `json:"postId"`   // int64.toString()
}

type QuotedNoti struct {
	ID        int               `json:"id"`
	Type      NotiType          `json:"type"`
	EventTime time.Time         `json:"eventTime"`
	Content   QuotedNotiContent `json:"content"`

	SendTo SendTo `json:"-"`
}

func (n *NotiService) NewQuotedNoti(
	ctx context.Context, thread *Thread, post *Post, quotedPost *Post,
) error {
	noti := &QuotedNoti{
		Type:      NotiTypeQuoted,
		EventTime: time.Now(),
		Content: QuotedNotiContent{
			ThreadID: thread.ID.ToBase64String(),
			QuotedID: quotedPost.ID.ToBase64String(),
			PostID:   post.ID.ToBase64String(),
		},
		SendTo: SendToUser(quotedPost.Data.Author.UserID),
	}
	return n.Repo.InsertNoti(ctx, NotiInsert{Quoted: noti})
}

func (n *QuotedNoti) HasRead(user *User) bool {
	return user.LastReadNoti.QuotedNoti >= n.ID
}

func (n *QuotedNoti) ThreadID() (uid.UID, error) {
	return uid.ParseUID(n.Content.ThreadID)
}

func (n *QuotedNoti) QuotedPostID() (uid.UID, error) {
	return uid.ParseUID(n.Content.QuotedID)
}

func (n *QuotedNoti) PostID() (uid.UID, error) {
	return uid.ParseUID(n.Content.PostID)
}
