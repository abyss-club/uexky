package model

import (
	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky/uexky"
)

// RoleType ...
type RoleType string

// Roles
const (
	Normal     RoleType = ""
	Banned     RoleType = "banned"
	TagAdmin   RoleType = "tag_admin"
	SuperAdmin RoleType = "super_admin"
)

var allActions = []string{
	"ban_user",
	"block_post",
	"lock_thread",
	"block_thread",
	"edit_tags",
}

type rolePriority struct {
	Role    RoleType
	Actions []string
}

var rolePriorities = []rolePriority{
	rolePriority{
		Role:    TagAdmin,
		Actions: allActions,
	},
	rolePriority{
		Role:    SuperAdmin,
		Actions: allActions,
	},
}

// CheckPriority ...
func (a *User) CheckPriority(
	u *uexky.Uexky, action, target string,
) (bool, error) {
	if a.Role.Type != TagAdmin && a.Role.Type != SuperAdmin {
		return false, nil
	} else if a.Role.Type == SuperAdmin {
		return true, nil
	}

	// TagAdmin:
	if len(a.Role.Range) == 0 {
		return false, nil
	}
	var findThread func(*uexky.Uexky, string) (*Thread, error)
	if action == "block_post" || action == "ban_user" {
		findThread = FindThreadByPostID
	} else {
		findThread = FindThreadByID
	}
	thread, err := findThread(u, target)
	if err != nil {
		return false, err
	}
	for _, tag := range a.Role.Range {
		if tag == thread.MainTag {
			return true, nil
		}
	}
	return false, nil
}

// BanUser ...
func BanUser(u *uexky.Uexky, postID string) error {
	if ok, err := u.Auth.CheckPriority("ban_user", postID); err != nil {
		return err
	} else if !ok {
		return errors.New("Permitted Error")
	}
	post, err := FindPost(u, postID)
	if err != nil {
		return err
	}
	var user *User
	c := u.Mongo.C(colleUser)
	if err := c.FindId(post.UserID).One(user); err != nil {
		return err
	}
	if user.Role.Type == TagAdmin || user.Role.Type == SuperAdmin {
		return errors.New("Permitted Error")
	}
	return c.Update(bson.M{"_id": user.ID}, bson.M{
		"$set": bson.M{"role.type": Banned},
	})
}

// BlockPost ...
func BlockPost(u *uexky.Uexky, postID string) error {
	if ok, err := u.Auth.CheckPriority("block_post", postID); err != nil {
		return err
	} else if !ok {
		return errors.New("Permitted Error")
	}
	return u.Mongo.C(collePost).Update(bson.M{"id": postID}, bson.M{
		"$set": bson.M{"blocked": true},
	})
}

// LockThread ...
func LockThread(u *uexky.Uexky, threadID string) error {
	if ok, err := u.Auth.CheckPriority("lock_thread", threadID); err != nil {
		return err
	} else if !ok {
		return errors.New("Permitted Error")
	}
	return u.Mongo.C(colleThread).Update(bson.M{"id": threadID}, bson.M{
		"$set": bson.M{"lock": true},
	})
}

// BlockThread ...
func BlockThread(u *uexky.Uexky, threadID string) error {
	if ok, err := u.Auth.CheckPriority("block_thread", threadID); err != nil {
		return err
	} else if !ok {
		return errors.New("Permitted Error")
	}
	return u.Mongo.C(colleThread).Update(bson.M{"id": threadID}, bson.M{
		"$set": bson.M{"block": true},
	})
}

// EditTags ...
func EditTags(u *uexky.Uexky, threadID, mainTag string, subTags []string) error {
	if ok, err := u.Auth.CheckPriority("edit_tags", threadID); err != nil {
		return err
	} else if !ok {
		return errors.New("Permitted Error")
	}
	tags := []string{mainTag}
	for _, tag := range subTags {
		tags = append(tags, tag)
	}
	return u.Mongo.C(colleThread).Update(bson.M{"id": threadID}, bson.M{
		"$set": bson.M{"main_tag": mainTag, "sub_tags": subTags, "tags": tags},
	})
}
