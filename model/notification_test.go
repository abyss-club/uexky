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
	quotes := []*Post{
		&Post{
			ID:     "NotiTestQuotedPost1",
			UserID: receiver.ID,
		},
		&Post{
			ID:     "NotiTestQuotedPost2",
			UserID: author.ID,
		},
	}

	// publish test post
	if err := TriggerNotifForPost(ctx, thread, post, quotes); err != nil {
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
	} else if c != 1 {
		t.Fatalf("GetUnreadNotificationCount(Replied) = %v, want = %v", c, 1)
	}
	if c, err := GetUnreadNotificationCount(ctx, NotiTypeQuoted); err != nil {
		t.Fatalf("GetUnreadNotificationCount(Quoted) error = %v", err)
	} else if c != 1 {
		t.Fatalf("GetUnreadNotificationCount(Quoted) = %v, want = %v", c, 1)
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
		want := &Notification{
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
		noti, slice, err := GetNotification(ctx, NotiTypeQuoted, sq)
		if err != nil {
			t.Fatalf("GetNotification(Quoted) error = %v", err)
		}
		if len(noti) != 1 {
			t.Fatalf("GetNotification(Quoted) != [], len = %v", len(noti))
		}
		want := &Notification{
			ID:        "quoted:NotiTestQuotedPost1:NotiTestPost",
			Type:      NotiTypeQuoted,
			SendTo:    receiver.ID,
			EventTime: post.CreateTime,
			HasRead:   false,
			Quoted: &QuotedNotiContent{
				ThreadID:     thread.ID,
				PostID:       post.ID,
				QuotedPostID: quotes[0].ID,
				Quoter:       post.Author,
				QuoterID:     post.UserID,
			},
		}
		if diff := cmp.Diff(want, noti[0], timeCmp); diff != "" {
			t.Fatalf("GetNotification(Quoted) error, diff = %+v", diff)
		}
		if slice.FirstCursor != noti[0].genCursor() ||
			slice.LastCursor != noti[0].genCursor() {
			t.Fatalf("GetNotification(Quoted).slice != {}, slice = %v", *slice)
		}
	}
}

func TestSendSystemNotification(t *testing.T) {
	ctx := ctxWithUser(mockUsers[2])
	title := "Test!"
	content := `This is a *markdown* Notification`
	if err := SendSystemNotification(ctx, title, content); err != nil {
		t.Fatalf("SendSystemNotification() error = %v", err)
	}
	t.Log("check unread count")
	{
		if c, err := GetUnreadNotificationCount(ctx, NotiTypeSystem); err != nil {
			t.Fatalf("GetUnreadNotificationCount(System) error = %v", err)
		} else if c != 1 {
			t.Fatalf("GetUnreadNotificationCount(System) = %v, want = 1", c)
		}
	}
	t.Log("check noti content")
	{
		sq := &SliceQuery{Limit: 10, Desc: true, Cursor: genTimeCursor(time.Now())}
		noti, slice, err := GetNotification(ctx, NotiTypeSystem, sq)
		if err != nil {
			t.Fatalf("GetNotification(System) error = %v", err)
		}
		if len(noti) != 1 {
			t.Fatalf("GetNotification(System) != [], len = %v", len(noti))
		}
		got := noti[0]
		want := &Notification{
			Type:        NotiTypeSystem,
			SendToGroup: AllUser,
			System: &SystemNotiContent{
				Title:   title,
				Content: content,
			},
		}
		if slice.FirstCursor != got.genCursor() ||
			slice.LastCursor != got.genCursor() {
			t.Fatalf("GetNotification(System).slice != {}, slice = %v", *slice)
		}
		if got.Type != want.Type {
			t.Fatalf("GetNotification(System).Type = %v, want %v",
				got.Type, want.Type)
		}
		if got.SendToGroup != want.SendToGroup {
			t.Fatalf("GetNotification(System).SendToGroup = %v, want %v",
				got.SendToGroup, want.SendToGroup)
		}
		if diff := cmp.Diff(want.System, got.System, timeCmp); diff != "" {
			t.Fatalf("GetNotification(System).System = %v, diff = %s",
				got.SendToGroup, diff)
		}
	}
}
