// forum notification: systemNoti repliedNoti quotedNoti

package entity

import (
	"context"
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"gitlab.com/abyss.club/uexky/lib/uerr"
	"gitlab.com/abyss.club/uexky/lib/uid"
)

type NotiService struct {
	Repo NotiRepo
}

func (n *NotiService) GetUnreadNotiCount(ctx context.Context, user *User) (int, error) {
	return n.Repo.GetUserUnreadCount(ctx, user)
}

func (n *NotiService) GetNotification(ctx context.Context, user *User, query SliceQuery) (*NotiSlice, error) {
	slice, err := n.Repo.GetNotiSlice(ctx, user, query)
	if err != nil {
		return nil, err
	}
	if len(slice.Notifications) > 0 {
		lastRead := slice.Notifications[0].SortKey
		if err := n.Repo.UpdateReadID(ctx, user.ID, lastRead); err != nil {
			return nil, err
		}
		user.LastReadNoti = lastRead
	}
	return slice, err
}

func (n *NotiService) NewSystemNoti(ctx context.Context, title, content string, receivers ...Receiver) error {
	if len(receivers) == 0 {
		return uerr.New(uerr.PermissionError, "no receivers")
	}
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

func (n *NotiService) NewRepliedNoti(ctx context.Context, user *User, thread *Thread, reply *Post) error {
	key := fmt.Sprintf("replied:%s", thread.ID.ToBase64String())
	oldNoti, err := n.Repo.GetNotiByKey(ctx, user, key)
	if err != nil {
		return err
	}
	content := RepliedNoti{
		Thread: &ThreadOutline{
			ID:      thread.ID,
			Title:   thread.Title,
			Content: thread.Content,
			MainTag: thread.MainTag,
			SubTags: thread.SubTags,
		},
		NewReplyID:    reply.ID,
		NewReplyCount: 1,
	}
	noti := &Notification{
		Type:      NotiTypeReplied,
		Key:       key,
		SortKey:   reply.ID,
		Receivers: []Receiver{SendToUser(thread.AuthorObj.UserID)},
	}
	if oldNoti != nil {
		if !oldNoti.HasRead {
			oldContent := oldNoti.Content.(RepliedNoti)
			content.NewReplyCount = oldContent.NewReplyCount + 1
			content.NewReplyID = oldContent.NewReplyID
		}
		noti.Content = content
		log.Infof("UpdateNotiContent(%#v), key=%s", noti, key)
		return n.Repo.UpdateNotiContent(ctx, noti)
	}
	noti.Content = content
	log.Infof("InsertNoti(%#v), key=%s", noti, key)
	return n.Repo.InsertNoti(ctx, noti)
}

func (n *NotiService) NewQuotedNoti(ctx context.Context, thread *Thread, post *Post, quotedPost *Post) error {
	noti := &Notification{
		Type:      NotiTypeQuoted,
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
	log.Infof("NewQuotedNoti, post %v quote post %v, key=%s", post, quotedPost, noti.Key)
	return n.Repo.InsertNoti(ctx, noti)
}

type Notification struct {
	Type      NotiType    `json:"type"`
	EventTime time.Time   `json:"eventTime"`
	HasRead   bool        `json:"hasRead"`
	Content   NotiContent `json:"content"`

	Key       string     `json:"-"` // use to merge notification, must be unique
	SortKey   uid.UID    `json:"-"` // use to sort and mark read for notification
	Receivers []Receiver `json:"-"`
}

func (n *Notification) DecodeContent(m map[string]interface{}) error {
	var err error
	switch n.Type {
	case NotiTypeSystem:
		var content SystemNoti
		err = mapstructure.Decode(m, &content)
		n.Content = content
	case NotiTypeReplied:
		var content RepliedNoti
		err = mapstructure.Decode(m, &content)
		n.Content = content
	case NotiTypeQuoted:
		var content QuotedNoti
		err = mapstructure.Decode(m, &content)
		n.Content = content
	default:
		err = fmt.Errorf("can't marshal noti content, invalid type '%s'", n.Type)
	}
	return err
}

func (n *Notification) EncodeContent() (map[string]interface{}, error) {
	m := make(map[string]interface{})
	var err error
	switch n.Type {
	case NotiTypeSystem:
		err = mapstructure.Decode(n.Content.(SystemNoti), &m)
	case NotiTypeReplied:
		err = mapstructure.Decode(n.Content.(RepliedNoti), &m)
	case NotiTypeQuoted:
		err = mapstructure.Decode(n.Content.(QuotedNoti), &m)
	default:
		err = uerr.Errorf(uerr.ParamsError, "invalid noti type '%s'", n.Type)
	}
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (n *Notification) String() string {
	return fmt.Sprintf("<Notification:%s:%s:%v>", n.Type, n.Key, n.SortKey)
}
