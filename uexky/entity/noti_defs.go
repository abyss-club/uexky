package entity

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky/lib/uid"
)

//-- NotiRepo

type NotiRepo interface {
	GetUserUnreadCount(ctx context.Context, user *User) (int, error)
	GetNotiByKey(ctx context.Context, userID int64, key string) (*Notification, error)
	GetNotiSlice(ctx context.Context, user *User, query SliceQuery) (*NotiSlice, error)
	InsertNoti(ctx context.Context, notification *Notification) error
	UpdateNotiContent(ctx context.Context, noti *Notification) error
	UpdateReadID(ctx context.Context, userID int64, id uid.UID) error
}

// -- NotiType

func ParseNotiType(s string) (NotiType, error) {
	t := NotiType(s)
	if !t.IsValid() {
		return NotiType(""), errors.Errorf("invalid noti type: %s", s)
	}
	return t, nil
}

// -- Receiver

type SendGroup string

const AllUser SendGroup = "all_user"

type Receiver string

func SendToUser(userID int64) Receiver {
	return Receiver(fmt.Sprintf("u:%v", userID))
}

func SendToGroup(group SendGroup) Receiver {
	return Receiver(fmt.Sprintf("g:%v", group))
}

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
