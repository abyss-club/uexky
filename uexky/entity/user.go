package entity

import (
	"context"
	"time"

	"gitlab.com/abyss.club/uexky/lib/algo"
	"gitlab.com/abyss.club/uexky/lib/uerr"
	"gitlab.com/abyss.club/uexky/lib/uid"
)

type UserRepo interface {
	// Read
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id uid.UID) (*User, error)
	GetGuestByID(ctx context.Context, id uid.UID) (*User, error)

	// Write
	Insert(ctx context.Context, user *User) (*User, error)
	Update(ctx context.Context, user *User) (*User, error)

	// Related
	ThreadSlice(ctx context.Context, user *User, query SliceQuery) (*ThreadSlice, error)
	PostSlice(ctx context.Context, user *User, query SliceQuery) (*PostSlice, error)
}

type User struct {
	ID           uid.UID  `json:"-"`
	Email        *string  `json:"email"`
	Name         *string  `json:"name"`
	Role         Role     `json:"role"`
	Tags         []string `json:"tags"`
	LastReadNoti uid.UID  `json:"-"`
}

const GuestExpireTime = 30 * time.Hour * 24

func NewSignedInUser(email string) *User {
	return &User{
		ID:    uid.NewUID(),
		Email: &email,
		Role:  RoleNormal,
	}
}

func NewGuestUser(id uid.UID) *User {
	return &User{
		ID:   id,
		Role: RoleGuest,
	}
}

func (u *User) Ban() {
	u.Role = RoleBanned
}

func (u *User) SetName(name string) error {
	if u.Name != nil {
		return uerr.New(uerr.ParamsError, "already have a name")
	}
	u.Name = algo.NullString(name)
	return nil
}

func (u *User) SetTags(tags []string) {
	repeat := map[string]bool{}
	var ts []string
	for _, tag := range tags {
		if !repeat[tag] {
			ts = append(ts, tag)
		}
		repeat[tag] = true
	}
	u.Tags = ts
}

func (u *User) AddTag(tag string) {
	for _, t := range u.Tags {
		if t == tag {
			return
		}
	}
	u.Tags = append(u.Tags, tag)
}

func (u *User) DelTag(tag string) {
	var tags []string
	for _, t := range u.Tags {
		if t != tag {
			tags = append(tags, t)
		}
	}
	u.Tags = tags
}

func (u *User) NotiReceivers() []Receiver {
	return []Receiver{SendToUser(u.ID), SendToGroup(AllUser)}
}

func (u *User) UpdateReadID(readID uid.UID) {
	u.LastReadNoti = readID
}

// ---- save to and get from context ----

type contextKey int

const userKey contextKey = 1

// contextUser used for keep user referance to sync change
type contextUser struct {
	User *User
}

// GetCurrentUser maybe nil
func GetCurrentUser(ctx context.Context) *User {
	user, ok := ctx.Value(userKey).(contextUser)
	if !ok {
		return nil
	}
	return user.User
}

func (u *User) AttachContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, userKey, contextUser{User: u})
}

// ---- roles and permissions ----

func ParseRole(s string) Role {
	if s == "" {
		return RoleNormal
	}
	r := Role(s)
	if !r.IsValid() {
		return RoleBanned
	}
	return r
}

func (r Role) Value() int {
	switch r {
	case RoleAdmin:
		return 100
	case RoleMod:
		return 10
	case RoleNormal:
		return 1
	case RoleGuest:
		return 0
	case RoleBanned:
		return -1
	default:
		return -10
	}
}

type Action string

const (
	ActionProfile     = Action("PROFILE")
	ActionBanUser     = Action("BAN_USER")
	ActionPromoteUser = Action("PROMOTE_USER")
	ActionBlockPost   = Action("BLOCK_POST")
	ActionLockThread  = Action("LOCK_THREAD")
	ActionBlockThread = Action("BLOCK_THREAD")
	ActionEditTag     = Action("EDIT_TAG")
	ActionEditSetting = Action("EDIT_SETTING")
	ActionPubPost     = Action("PUB_POST")
	ActionPubThread   = Action("PUB_THREAD")
)

var ActionRole = map[Action]Role{
	ActionProfile:     RoleBanned, // Because a user can only read the profile own by himself.
	ActionBanUser:     RoleMod,
	ActionPromoteUser: RoleAdmin,
	ActionBlockPost:   RoleMod,
	ActionLockThread:  RoleMod,
	ActionBlockThread: RoleMod,
	ActionEditTag:     RoleMod,
	ActionEditSetting: RoleAdmin,
	ActionPubPost:     RoleGuest,
	ActionPubThread:   RoleGuest,
}

func (u *User) RequirePermission(action Action) error {
	if u == nil {
		return uerr.New(uerr.AuthError, "permission denied, no user found")
	}
	needRole := ActionRole[action]
	if u.Role.Value() < needRole.Value() {
		return uerr.New(uerr.PermissionError, "permission denied")
	}
	return nil
}
