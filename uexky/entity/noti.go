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

type NotiRepo interface {
	GetUnreadCount(ctx context.Context, user *User) (int, error)
	GetByKey(ctx context.Context, user *User, key string) (*Notification, error)
	GetSlice(ctx context.Context, user *User, query SliceQuery) (*NotiSlice, error)
	Insert(ctx context.Context, notification *Notification) error

	UpdateContent(ctx context.Context, noti *Notification) error
	UpdateReadID(ctx context.Context, user *User, id uid.UID) error
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

// -- Receiver

type SendGroup string

const AllUser SendGroup = "all_user"

type Receiver string

func SendToUser(userID uid.UID) Receiver {
	return Receiver(fmt.Sprintf("u:%v", userID))
}

func SendToGroup(group SendGroup) Receiver {
	return Receiver(fmt.Sprintf("g:%v", group))
}

func NewSystemNoti(title, content string, receivers ...Receiver) (*Notification, error) {
	if len(receivers) == 0 {
		return nil, uerr.New(uerr.PermissionError, "no receivers")
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
	return noti, nil
}

func RepliedNotiKey(thread *Thread) string {
	return fmt.Sprintf("replied:%s", thread.ID.ToBase64String())
}

func NewRepliedNoti(prev *Notification, user *User, thread *Thread, reply *Post) *Notification {
	if thread.Author.Guest {
		return nil
	}
	log.Infof("NewRepliedNoti, User = %#v GetNotiByKey = %#v", user, prev)
	content := RepliedNoti{
		Thread: &ThreadOutline{
			ID:      thread.ID,
			Title:   thread.Title,
			Content: thread.Content,
			MainTag: thread.MainTag,
			SubTags: thread.SubTags,
		},
		FirstReplyID:    reply.ID,
		NewRepliesCount: 1,
	}
	noti := &Notification{
		Type:      NotiTypeReplied,
		Key:       RepliedNotiKey(thread),
		SortKey:   reply.ID,
		EventTime: time.Now(),
		Receivers: []Receiver{SendToUser(thread.Author.UserID)},
	}
	if prev != nil {
		if !prev.HasRead {
			oldContent := prev.Content.(RepliedNoti)
			content.NewRepliesCount = oldContent.NewRepliesCount + 1
			content.FirstReplyID = oldContent.FirstReplyID
		}
		noti.Content = content
		log.Infof("UpdateNotiContent(%#v), key=%s", noti, noti.Key)
		return noti
	}
	noti.Content = content
	log.Infof("InsertNoti(%#v), key=%s", noti, noti.Key)
	return noti
}

func QuotedNotiKey(quotedPost, post *Post) string {
	return fmt.Sprintf("quoted:%s:%s", quotedPost.ID.ToBase64String(), post.ID.ToBase64String())
}

func NewQuotedNoti(thread *Thread, post *Post, quotedPost *Post) *Notification {
	if quotedPost.Author.Guest {
		return nil
	}
	noti := &Notification{
		Type:      NotiTypeQuoted,
		EventTime: time.Now(),
		Content: QuotedNoti{
			ThreadID: thread.ID,
			QuotedPost: &PostOutline{
				ID:      quotedPost.ID,
				Author:  quotedPost.Author,
				Content: quotedPost.Content,
			},
			Post: &PostOutline{
				ID:      post.ID,
				Author:  post.Author,
				Content: post.Content,
			},
		},

		Key:       QuotedNotiKey(quotedPost, post),
		SortKey:   uid.NewUID(),
		Receivers: []Receiver{SendToUser(quotedPost.Author.UserID)},
	}
	log.Infof("NewQuotedNoti, post %v quote post %v, key=%s", post, quotedPost, noti.Key)
	return noti
}

func (n *Notification) String() string {
	return fmt.Sprintf("<Notification:%s:%s:%v>", n.Type, n.Key, n.SortKey)
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
	return uerr.Wrap(uerr.InternalError, err, "DecodeContent")
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
		err = fmt.Errorf("invalid noti type '%s'", n.Type)
	}
	if err != nil {
		return nil, uerr.Wrap(uerr.ParamsError, err, "EncodeContent")
	}
	return m, nil
}

// ---- special notifications ----

const (
	WelcomeTitle   = "欢迎来到 abyss!"
	WelcomeContent = `这是一个可匿名、标签化的讨论版。目标以现代观念打造美观便利的体验，
从聊天工具和社交网络的信息洪流中，回归到明晰的讨论。通过动画、游戏或者更多自定义的标签寻找和参与感兴趣的话题吧。

发言之前请阅读社区规则和隐私声明。这是 abyss 的第一个公开测试版本，尚有诸多不足，欢迎使用 abyss 标签给我们建议和反馈。我们的代码全部以 AGPL [开源](https://gitlab.com/abyss.club/abyss)，欢迎提出 issue 和 PR。

---

### 隐私声明

你可以在此保持匿名，abyss 不会收集任何用户信息，也不保证任何用户的身份。注册之后我们将随帐号记录邮箱地址用于登录；在你使用时会暂存IP地址用于流量控制。以上信息均不会透露给其他用户。

---

### 社区规则

1. 符合任一主标签内容话题均可以畅所欲言。如需更多的主标签，可以使用 abyss 标签发帖申请
2. 禁止歧视，仇恨，反人类的言论；禁止对性、儿童色情直接描写的言论和图片。
3. 由于不便明言的原因，暂时禁止讨论和隐射有关中国的政治话题。
`
)

func NewWelcomeNoti(user *User) (*Notification, error) {
	return NewSystemNoti(WelcomeTitle, WelcomeContent, SendToUser(user.ID))
}
