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
	NotiTypeQuoted  NotiType = "quoted"
)

var allNotiTypes = map[NotiType]bool{
	NotiTypeSystem:  true,
	NotiTypeReplied: true,
	NotiTypeQuoted:  true,
}

// UserGroup ...
type UserGroup string

// UserGroups
const (
	AllUser UserGroup = "all"
)

// SystemNotiContent ...
type SystemNotiContent struct {
	Title   string `bson:"title"`
	Content string `bson:"content"`
}

// RepliedNotiContent ...
type RepliedNotiContent struct {
	ThreadID   string          `bson:"thread_id"`
	Repliers   []string        `bson:"repliers"`
	ReplierIDs []bson.ObjectId `bson:"replier_ids"`
}

// QuotedNotiContent ...
type QuotedNotiContent struct {
	ThreadID string        `bson:"thread_id"`
	PostID   string        `bson:"post_id"`
	Quoter   string        `bson:"quoter"`
	QuoterID bson.ObjectId `bson:"quoter_id"`
}

// Notification for save notification in DB
type Notification struct {
	// base info
	ID          string        `bson:"id"`
	Type        NotiType      `bson:"type"`
	SendTo      bson.ObjectId `bson:"send_to"`
	SendToGroup UserGroup     `bson:"send_to_group"`
	EventTime   time.Time     `bson:"event_time"`
	HasRead     bool          `bson:"-"`

	System  *SystemNotiContent  `bson:"system"`
	Replied *RepliedNotiContent `bson:"replied"`
	Quoted  *QuotedNotiContent  `bson:"quoted"`
}

func (ns *Notification) genCursor() string {
	return genTimeCursor(ns.EventTime)
}

func (ns *Notification) checkIfRead(user *User, t NotiType) {
	switch t {
	case NotiTypeSystem:
		ns.HasRead = user.ReadNotiTime.System.After(ns.EventTime)
	case NotiTypeReplied:
		ns.HasRead = user.ReadNotiTime.Replied.After(ns.EventTime)
	case NotiTypeQuoted:
		ns.HasRead = user.ReadNotiTime.Quoted.After(ns.EventTime)
	}
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
) ([]*Notification, *SliceInfo, error) {
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

	var noti []*Notification
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
	ctx context.Context, thread *Thread, post *Post, quotes []*Post,
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

	for _, q := range quotes {
		if post.UserID == q.UserID {
			continue
		}
		qn := &Notification{
			ID:        fmt.Sprintf("quoted:%v", q.ID),
			Type:      NotiTypeQuoted,
			SendTo:    q.UserID,
			EventTime: post.CreateTime,
			Quoted: &QuotedNotiContent{
				ThreadID: thread.ID,
				PostID:   q.ID,
				Quoter:   post.Author,
				QuoterID: post.UserID,
			},
		}
		if err := c.Insert(qn); err != nil {
			return err
		}
	}
	return nil
}
