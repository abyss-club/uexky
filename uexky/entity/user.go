// user aggragate: user

package entity

import (
	"context"
	"errors"
	"fmt"
)

type UserRepo interface{}

type MailService interface{}

type UserService struct {
	Repo UserRepo
	Mail MailService
}

func (s *UserService) GetUserRequirePermission(ctx context.Context, action Action) (*User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (s *UserService) SignInByEmail(ctx context.Context, email string) (bool, error) {
	// user, err := s.GetCurrentUser(ctx)
	// if err != nil {
	// 	return false, err
	// }
	// if user != nil {
	// 	return false, errors.New("you already signed in")
	// }
	panic(fmt.Errorf("not implemented"))
}

type Role string

const (
	RoleNormal = Role("")
	RoleGuest  = Role("guest")
	RoleMod    = Role("mod")
	RoleBanned = Role("banned")
	RoleAdmin  = Role("admin")
)

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
	ActionBlockPost   = Action("BLOCK_POST")
	ActionLockThread  = Action("LOCK_THREAD")
	ActionBlockThread = Action("BLOCK_THREAD")
	ActionEditTag     = Action("EDIT_TAG")
	ActionEditSetting = Action("EDIT_SETTING")
	ActionPubPost     = Action("PUB_POST")
	ActionPubThread   = Action("PUB_THREAD")
)

var ActionRole = map[Action]Role{
	ActionProfile:     RoleNormal,
	ActionBanUser:     RoleMod,
	ActionBlockPost:   RoleMod,
	ActionLockThread:  RoleMod,
	ActionBlockThread: RoleMod,
	ActionEditTag:     RoleMod,
	ActionEditSetting: RoleAdmin,
}

type User struct {
	Email string  `json:"email"`
	Name  *string `json:"name"`
	Role  *string `json:"role"`
}

func (u *User) EnsurePermission(ctx context.Context, action Action) error {
	needRole := ActionRole[action]
	userRole := Role("")
	if u.Role != nil {
		userRole = Role(*u.Role)
	}
	if userRole.Value() < needRole.Value() {
		return errors.New("permission denied")
	}
	return nil
}

func (u *User) SetName(ctx context.Context, name string) error {
	panic(fmt.Errorf("not implemented"))
}

func (u *User) BanUser(ctx context.Context, name string, anonymous bool) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}
