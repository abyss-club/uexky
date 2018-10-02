package model

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/globalsign/mgo/bson"
	"gitlab.com/abyss.club/uexky/mw"
)

// NotifType ...
type NotifType string

// NotifTypes
const (
	NotifTypeSystem  NotifType = "system"
	NotifTypeReplied NotifType = "replied"
	NotifTypeRefered NotifType = "refered"
)

// UserGroup ...
type UserGroup string

// UserGroups
const (
	AllUser UserGroup = "all"
)

// Notification ...
type Notification struct {
	SendTo       bson.ObjectId `bson:"send_to"`
	SendToGroup  UserGroup     `bson:"send_to_group"`
	ReleaseTime  time.Time     `bson:"release_time"`
	NotifType    NotifType     `bson:"notif_type"`
	SystemNotif  *systemNotif  `bson:"system_notif"`
	RepliedNotif *repliedNotif `bson:"replied_notif"`
	ReferedNotif *referedNotif `bson:"refered_notif"`
}

type systemNotif struct {
	Content string `json:"content" bson:"content"`
}

type repliedNotif struct {
	ThreadID    string   `json:"thread_id" bson:"thread_id"`
	ThreadTitle string   `json:"thread_title" bson:"thread_title"`
	Tags        []string `json:"tags" bson:"tags"`
	Repliers    []string `json:"repliers" bson:"repliers"`
}

type referedNotif struct {
	ThreadID    string   `json:"thread_id" bson:"thread_id"`
	ThreadTitle string   `json:"thread_title" bson:"thread_title"`
	Tags        []string `json:"tags" bson:"tags"`
	Content     string   `json:"content" bson:"content"`
	Referers    []string `json:"referers" bson:"referers"`
}

// GetContent ...
func (n *Notification) GetContent() string {
	var content []byte
	switch n.NotifType {
	case "system":
		content, _ = json.Marshal(n.SystemNotif)
	case "replied":
		content, _ = json.Marshal(n.RepliedNotif)
	case "refered":
		content, _ = json.Marshal(n.ReferedNotif)
	default:
		log.Fatal("Unknown Notification Type")
	}
	return string(content)
}

func (n *Notification) genCursor() string {
	return genTimeCursor(n.ReleaseTime)
}

func reverseNotification(notif []*Notification) {
	l := len(notif)
	for i := 0; i != l/2; i++ {
		notif[i], notif[l-i-1] = notif[l-i-1], notif[i]
	}
}

// GetUnreadNotificationCount ...
func GetUnreadNotificationCount(ctx context.Context) (int, error) {
	user, err := requireSignIn(ctx)
	if err != nil {
		return 0, err
	}
	c := mw.GetMongo(ctx).C(colleNotification)
	c.EnsureIndexKey("release_time", "send_to")
	c.EnsureIndexKey("release_time", "send_to_group")
	query := bson.M{
		"release_time": bson.M{"$lt": user.ReadNotifTime},
		"$or": []bson.M{
			bson.M{"send_to": user.ID},
			bson.M{"send_to_group": AllUser},
		},
	}
	count, err := c.Find(query).Count()
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetNotificationByUser ...
func GetNotificationByUser(
	ctx context.Context, sq *SliceQuery,
) ([]*Notification, *SliceInfo, error) {
	user, err := requireSignIn(ctx)
	if err != nil {
		return nil, nil, err
	}

	query, err := sq.GenQueryByTime("release_time")
	if err != nil {
		return nil, nil, err
	}
	query["$or"] = []bson.M{
		bson.M{"send_to": user.ID},
		bson.M{"send_to_group": AllUser},
	}

	var notifs []*Notification
	err = sq.Find(ctx, colleNotification, "release_time", query, notifs)
	if err != nil {
		return nil, nil, err
	}
	now := time.Now()
	c := mw.GetMongo(ctx).C(colleUser)
	if err := c.Update(bson.M{"_id": user.ID}, bson.M{"read_notif_time": now}); err != nil {
		return nil, nil, err
	}
	user.ReadNotifTime = now

	if len(notifs) == 0 {
		return notifs, &SliceInfo{}, nil
	}

	reverseNotification(notifs)
	return notifs, &SliceInfo{
		FirstCursor: notifs[0].genCursor(),
		LastCursor:  notifs[len(notifs)-1].genCursor(),
	}, nil
}
