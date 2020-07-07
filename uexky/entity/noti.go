// forum notification: systemNoti repliedNoti quotedNoti

package entity

import (
	"context"
	"fmt"
	"time"

	"gitlab.com/abyss.club/uexky/lib/uid"
)

type NotiService struct {
	Repo NotiRepo
}

func (n *NotiService) GetUnreadNotiCount(ctx context.Context, user *User) (int, error) {
	return n.Repo.GetUserUnreadCount(ctx, user)
}

func (n *NotiService) GetNotification(ctx context.Context, user *User, query SliceQuery) (*NotiSlice, error) {
	return n.Repo.GetNotiSlice(ctx, &NotiSearch{UserID: user.ID}, query)
}

func (n *NotiService) NewSystemNoti(ctx context.Context, title, content string, receivers []Receiver) error {
	noti := &Notification{
		Type:      NotiTypeSystem,
		EventTime: time.Now(),
		Content: SystemNoti{
			Title:   title,
			Content: content,
		},

		SortKey:   uid.NewUID(),
		Receivers: receivers,
	}
	noti.Key = noti.SortKey.ToBase64String()
	return n.Repo.InsertNoti(ctx, noti)
}

func (n *NotiService) NewRepliedNoti(ctx context.Context, thread *Thread, reply *Post) error {
	noti := &Notification{
		Type:      NotiTypeSystem,
		EventTime: time.Now(),
		Content: RepliedNoti{
			Thread: &ThreadOutline{
				ID:      thread.ID,
				Title:   thread.Title,
				Content: thread.Content,
				MainTag: thread.MainTag,
				SubTags: thread.SubTags,
			},
			NewReplyID: reply.ID, // repo should handle merge logic
		},

		Key:       fmt.Sprintf("reply:%s", thread.ID.ToBase64String()),
		SortKey:   uid.NewUID(),
		Receivers: []Receiver{SendToUser(thread.AuthorObj.UserID)},
	}
	noti.Key = noti.SortKey.ToBase64String()
	return n.Repo.InsertNoti(ctx, noti)
}

func (n *NotiService) NewQuotedNoti(ctx context.Context, thread *Thread, post *Post, quotedPost *Post) error {
	noti := &Notification{
		Type:      NotiTypeSystem,
		EventTime: time.Now(),
		Content: QuotedNoti{
			ThreadID: thread.ID,
			QuotedPost: &PostOutline{
				Author:  quotedPost.Author(),
				Content: quotedPost.Content,
			},
			Post: &PostOutline{
				Author:  post.Author(),
				Content: post.Content,
			},
		},

		Key:       fmt.Sprintf("quoted:%s:%s", quotedPost.ID.ToBase64String(), post.ID.ToBase64String()),
		SortKey:   uid.NewUID(),
		Receivers: []Receiver{SendToUser(quotedPost.Data.Author.UserID)},
	}
	noti.Key = noti.SortKey.ToBase64String()
	return n.Repo.InsertNoti(ctx, noti)
}

type Notification struct {
	Type      NotiType    `json:"type"`
	EventTime time.Time   `json:"eventTime"`
	Content   NotiContent `json:"content"`

	Key       string     `json:"-"` // use to merge notification, must be unique
	SortKey   uid.UID    `json:"-"` // use to sort and mark read for notification
	Receivers []Receiver `json:"-"`
}

func (n *Notification) HasRead(user *User) bool {
	panic(fmt.Errorf("not implemented"))
}
