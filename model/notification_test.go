package model

import (
	"testing"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/google/go-cmp/cmp"
)

func TestTriggerNotifForPost(t *testing.T) {
	// prepare
	receiver := mockUsers[1]
	author := mockUsers[2]
	ctx := ctxWithUser(author)
	thread := &Thread{
		ID:         "NotiTestThread",
		UserID:     receiver.ID,
		CreateTime: time.Now(),
	}
	post := &Post{
		ID:         "NotiTestPost",
		UserID:     author.ID,
		Author:     "TestAuthor",
		CreateTime: time.Now(),
	}
	refers := []*Post{
		&Post{
			ID:     "NotiTestReferredPost1",
			UserID: receiver.ID,
		},
		&Post{
			ID:     "NotiTestReferredPost2",
			UserID: author.ID,
		},
	}

	// publish test post
	if err := TriggerNotifForPost(ctx, thread, post, refers); err != nil {
		t.Fatalf("TriggerNotifForPost() error = %v", err)
	}

	// check unread count
	ctx = ctxWithUser(receiver)
	if c, err := GetUnreadNotificationCount(ctx, NotiTypeSystem); err != nil {
		t.Fatalf("GetUnreadNotificationCount(System) error = %v", err)
	} else if c != 0 {
		t.Fatalf("GetUnreadNotificationCount(System) = %v, want = %v", c, 0)
	}
	if c, err := GetUnreadNotificationCount(ctx, NotiTypeReplied); err != nil {
		t.Fatalf("GetUnreadNotificationCount(Replied) error = %v", err)
	} else if c != 0 {
		t.Fatalf("GetUnreadNotificationCount(Replied) = %v, want = %v", c, 1)
	}
	if c, err := GetUnreadNotificationCount(ctx, NotiTypeReferred); err != nil {
		t.Fatalf("GetUnreadNotificationCount(Referred) error = %v", err)
	} else if c != 0 {
		t.Fatalf("GetUnreadNotificationCount(Referred) = %v, want = %v", c, 1)
	}

	// check notification
	sq := &SliceQuery{Limit: 10, Desc: true, Cursor: genTimeCursor(time.Now())}
	{
		noti, slice, err := GetNotification(ctx, NotiTypeSystem, sq)
		if err != nil {
			t.Fatalf("GetNotification(System) error = %v", err)
		}
		if len(noti) != 0 {
			t.Fatalf("GetNotification(System) != [], len = %v", len(noti))
		}
		if slice.FirstCursor != "" || slice.LastCursor != "" {
			t.Fatalf("GetNotification(System).slice != {}, slice = %v", *slice)
		}
	}
	{
		noti, slice, err := GetNotification(ctx, NotiTypeReplied, sq)
		if err != nil {
			t.Fatalf("GetNotification(Replied) error = %v", err)
		}
		if len(noti) != 1 {
			t.Fatalf("GetNotification(Replied).len != 1, len = %v", len(noti))
		}
		want := &NotiStore{
			ID:        "replied:NotiTestThread",
			Type:      NotiTypeReplied,
			SendTo:    receiver.ID,
			EventTime: post.CreateTime,
			HasRead:   false,
			Replied: &RepliedNotiContent{
				ThreadID:   thread.ID,
				Repliers:   []string{post.Author},
				ReplierIDs: []bson.ObjectId{post.UserID},
			},
		}
		if diff := cmp.Diff(want, noti[0], timeCmp); diff != "" {
			t.Fatalf("GetNotification(Replied) error, diff = %v", diff)
		}
		if slice.FirstCursor != noti[0].genCursor() ||
			slice.LastCursor != noti[0].genCursor() {
			t.Fatalf("GetNotification(Replied).slice != {}, slice = %v", *slice)
		}
	}
	{
		noti, slice, err := GetNotification(ctx, NotiTypeReferred, sq)
		if err != nil {
			t.Fatalf("GetNotification(Referred) error = %v", err)
		}
		if len(noti) != 1 {
			t.Fatalf("GetNotification(Referred) != [], len = %v", len(noti))
		}
		want := &NotiStore{
			ID:        "referred:NotiTestReferredPost1",
			Type:      NotiTypeReferred,
			SendTo:    receiver.ID,
			EventTime: post.CreateTime,
			HasRead:   false,
			Referred: &ReferredNotiContent{
				ThreadID:   thread.ID,
				PostID:     refers[0].ID,
				Referers:   []string{post.Author},
				RefererIDs: []bson.ObjectId{post.UserID},
			},
		}
		if diff := cmp.Diff(want, noti[0], timeCmp); diff != "" {
			t.Fatalf("GetNotification(Referred) error, diff = %+v", diff)
		}
		if slice.FirstCursor != noti[0].genCursor() ||
			slice.LastCursor != noti[0].genCursor() {
			t.Fatalf("GetNotification(Referred).slice != {}, slice = %v", *slice)
		}
	}
}
