package service

import "context"

type UserRepository interface {
	FindUsersByIDs(ctx context.Context, id []int) ([]*User, error)
	UpdateUser(ctx context.Context, user *User) error
}

type User struct {
	ID    int     `json:"id"`
	Name  *string `json:"name"`
	Level int     `json:"level"`

	FriendIDs []int
	Repo      UserRepository
}

func (u *User) GetFriends(ctx context.Context) ([]*User, error) {
	return u.Repo.FindUsersByIDs(ctx, u.FriendIDs)
}
