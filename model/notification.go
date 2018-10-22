package model

import (
	"context"
	"fmt"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky/mw"
)

// NotiType ...
type NotiType string

// NotifTypes
const (
	NotiTypeSystem  NotiType = "system"
	NotiTypeReplied NotiType = "replied"
	NotiTypeRefered NotiType = "refered"
)

var allNotiTypes = map[NotiType]bool{
	NotiTypeSystem:  true,
	NotiTypeReplied: true,
	NotiTypeRefered: true,
}

// UserGroup ...
type UserGroup string

// UserGroups
const (
	AllUser UserGroup = "all"
)

// BaseNoti ...
type BaseNoti struct {
	ID          string        `bson:"id"`
	Type        NotiType      `bson:"type"`
	SendTo      bson.ObjectId `bson:"send_to"`
	SendToGroup UserGroup     `bson:"send_to_group"`
	EventTime   time.Time     `bson:"event_time"`
	HasRead     bool          `bson:"-"`
}

// SystemNotiContent ...
type SystemNotiContent struct {
	Title   string `bson:"title"`
	Content string `bson:"content"`
}

// SystemNoti ...
type SystemNoti struct {
	*BaseNoti
	*SystemNotiContent
}

// RepliedNotiContent ...
type RepliedNotiContent struct {
	ThreadID   string          `bson:"thread_id"`
	Repliers   []string        `bson:"repliers"`
	ReplierIDs []bson.ObjectId `bson:"replier_ids"`
}

// RepliedNoti ...
type RepliedNoti struct {
	*BaseNoti
	*RepliedNotiContent
}

// ReferedNotiContent ...
type ReferedNotiContent struct {
	ThreadID   string          `bson:"thread_id"`
	PostID     string          `bson:"post_id"`
	Referers   []string        `bson:"repliers"`
	RefererIDs []bson.ObjectId `bson:"replier_ids"`
}

// ReferedNoti ...
type ReferedNoti struct {
	*BaseNoti
	*ReferedNotiContent
}

// NotiStore for save notification in DB
type NotiStore struct {
	// BaseNoti...
	ID          string        `bson:"id"`
	Type        NotiType      `bson:"type"`
	SendTo      bson.ObjectId `bson:"send_to"`
	SendToGroup UserGroup     `bson:"send_to_group"`
	EventTime   time.Time     `bson:"event_time"`
	HasRead     bool          `bson:"-"`

	System  *SystemNotiContent  `bson:"system"`
	Replied *RepliedNotiContent `bson:"replied"`
	Refered *ReferedNotiContent `bson:"refered"`
}

func (ns *NotiStore) genCursor() string {
	return genTimeCursor(ns.EventTime)
}

func (ns *NotiStore) checkIfRead(user *User, t NotiType) {
	switch t {
	case NotiTypeSystem:
		ns.HasRead = user.ReadNotiTime.System.After(ns.EventTime)
	case NotiTypeReplied:
		ns.HasRead = user.ReadNotiTime.Replied.After(ns.EventTime)
	case NotiTypeRefered:
		ns.HasRead = user.ReadNotiTime.Refered.After(ns.EventTime)
	}
}

func (ns *NotiStore) baseNoti() *BaseNoti {
	return &BaseNoti{ns.ID, ns.Type, ns.SendTo, ns.SendToGroup,
		ns.EventTime, ns.HasRead}
}

// GetSystemNoti ...
func (ns *NotiStore) GetSystemNoti() *SystemNoti {
	if ns.System == nil {
		return nil
	}
	return &SystemNoti{ns.baseNoti(), ns.System}
}

// GetRepliedNoti ...
func (ns *NotiStore) GetRepliedNoti() *RepliedNoti {
	if ns.Replied == nil {
		return nil
	}
	return &RepliedNoti{ns.baseNoti(), ns.Replied}
}

// GetReferedNoti ...
func (ns *NotiStore) GetReferedNoti() *ReferedNoti {
	if ns.Refered == nil {
		return nil
	}
	return &ReferedNoti{ns.baseNoti(), ns.Refered}
}

// GetUnreadNotificationCount ...
func GetUnreadNotificationCount(ctx context.Context, t NotiType) (int, error) {
	if !allNotiTypes[t] {
		return 0, errors.Errorf("Invalidate notification type: %v", t)
	}
	user, err := requireSignIn(ctx)
	if err != nil {
		return 0, err
	}
	c := mw.GetMongo(ctx).C(colleNotification)
	c.EnsureIndexKey("send_to", "type", "release_time")
	c.EnsureIndexKey("send_to_group", "type", "release_time")
	query := bson.M{
		"$or": []bson.M{
			bson.M{"send_to": user.ID},
			bson.M{"send_to_group": AllUser},
		},
		"type":       t,
		"event_time": bson.M{"$lt": user.getReadNotiTime(t)},
	}
	count, err := c.Find(query).Count()
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetNotification ...
func GetNotification(
	ctx context.Context, t NotiType, sq *SliceQuery,
) ([]*NotiStore, *SliceInfo, error) {
	if !allNotiTypes[t] {
		return nil, nil, errors.Errorf("Invalidate notification type: %v", t)
	}
	user, err := requireSignIn(ctx)
	if err != nil {
		return nil, nil, err
	}
	query, err := sq.GenQueryByTime("event_time")
	if err != nil {
		return nil, nil, err
	}
	query["$or"] = []bson.M{
		bson.M{"send_to": user.ID},
		bson.M{"send_to_group": AllUser},
	}
	query["type"] = t

	var noti []*NotiStore
	now := time.Now()
	err = sq.Find(ctx, colleNotification, "event_time", query, &noti)
	if err != nil {
		return nil, nil, err
	}
	for _, n := range noti {
		n.checkIfRead(user, t)
	}
	user.setReadNotiTime(ctx, t, now)

	if len(noti) == 0 {
		return noti, &SliceInfo{}, nil
	}

	ReverseSlice(noti)
	return noti, &SliceInfo{
		FirstCursor: noti[0].genCursor(),
		LastCursor:  noti[len(noti)-1].genCursor(),
	}, nil
}

// trigger notifications by event:

// TriggerNotifForPost ...
func TriggerNotifForPost(
	ctx context.Context, thread *Thread, post *Post, refers []*Post,
) error {
	c := mw.GetMongo(ctx).C(colleNotification)
	c.EnsureIndexKey("id")
	id := fmt.Sprintf("replied:%v", thread.ID)
	if post.UserID != thread.UserID {
		if _, err := c.Upsert(bson.M{"id": id}, bson.M{
			"$set": bson.M{
				"id":                id,
				"type":              NotiTypeReplied,
				"send_to":           thread.UserID,
				"event_time":        post.CreateTime,
				"replied.thread_id": thread.ID,
			},
			"$addToSet": bson.M{
				"replied.repliers":    post.Author,
				"replied.replier_ids": post.UserID,
			},
		}); err != nil {
			return err
		}
	}

	for _, rp := range refers {
		if post.UserID == rp.UserID {
			continue
		}
		id := fmt.Sprintf("refered:%v", rp.ID)
		if _, err := c.Upsert(bson.M{"id": id}, bson.M{
			"$set": bson.M{
				"id":                id,
				"type":              NotiTypeRefered,
				"send_to":           rp.UserID,
				"event_time":        post.CreateTime,
				"refered.thread_id": thread.ID,
				"refered.post_id":   rp.ID,
			},
			"$addToSet": bson.M{
				"refered.repliers":    post.Author,
				"refered.replier_ids": post.UserID,
			},
		}); err != nil {
			return err
		}
	}
	return nil
}
