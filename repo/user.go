package repo

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gitlab.com/abyss.club/uexky/service"
)

type user struct {
	ID        int     `json:"id"`
	Name      *string `json:"name,omitempty"`
	Level     int     `json:"level"`
	FriendIDs []int   `json:"friend_ids"`
}

var usersData = func() []*user {
	var data []*user
	err := json.Unmarshal([]byte(`[
		{"id": 1, "level": 2, "friend_ids": [2, 3]},
		{"id": 2, "level": 0, "name": "arch", "friend_ids": [1, 3]},
		{"id": 3, "level": 5, "name": "gentoo", "friend_ids": [1, 2]},
		{"id": 4, "level": 9, "friend_ids": []}
	]`), &data)
	if err != nil {
		panic(err)
	}
	return data
}()

type UserRepository struct{}

func (u *UserRepository) NewServiceUser(user *user) *service.User {
	return &service.User{
		ID:        user.ID,
		Name:      user.Name,
		Level:     user.Level,
		FriendIDs: user.FriendIDs,
		Repo:      u,
	}
}

func (u *UserRepository) FindUsersByIDs(ctx context.Context, ids []int) ([]*service.User, error) {
	log.Infof("users data= %#v", usersData)
	for _, ud := range usersData {
		log.Infof("\t<%v>:\t%+v", ud.ID, *ud)
	}
	var sUsers []*service.User
	var notFound []int
	for _, id := range ids {
		found := false
		for _, user := range usersData {
			if user.ID == id {
				sUsers = append(sUsers, u.NewServiceUser(user))
				found = true
				break
			}
		}
		if !found {
			notFound = append(notFound, id)
		}
	}
	if len(notFound) > 0 {
		return nil, errors.Errorf("users not found: %v", notFound)
	}
	return sUsers, nil
}

func (u *UserRepository) UpdateUser(ctx context.Context, sUser *service.User) error {
	log.Infof("update user: %+v", sUser)
	if sUser.ID == 0 {
		return errors.New("user id is missing")
	}
	var user *user
	for _, u := range usersData {
		if u.ID == sUser.ID {
			user = u
		}
	}
	if u == nil {
		return errors.New("not found")
	}
	if sUser.Name != nil {
		if (*sUser.Name) == "" {
			user.Name = nil
		} else {
			user.Name = sUser.Name
		}
	}
	if sUser.Level != 0 {
		user.Level = sUser.Level
	}
	return nil
}
