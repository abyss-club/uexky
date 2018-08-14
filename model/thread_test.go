package model

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky/mgmt"
)

func TestNewThread(t *testing.T) {
	user := mockUsers[0]
	t.Log("user:", user)
	ctx := ctxWithUser(user)
	titles := []string{"thread1"}
	tests := []struct {
		name    string
		input   *ThreadInput
		check   bool
		wantErr bool
	}{
		{"normal, signed, titled", &ThreadInput{
			false, "thread1", mgmt.Config.MainTags[0],
			&[]string{"tag1", "tag2"}, &titles[0],
		}, true, false},
		{"normal, anonymous, non-title", &ThreadInput{
			true, "thread2", mgmt.Config.MainTags[0],
			&[]string{"tag1", "tag2"}, nil,
		}, true, false},
		{"error, no-main-tag", &ThreadInput{
			true, "thread3", "em..", nil, nil,
		}, false, true},
		{"error, multi-main-tag", &ThreadInput{
			true, "thread3", mgmt.Config.MainTags[0], &[]string{mgmt.Config.MainTags[1]}, nil,
		}, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewThread(ctx, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewThread() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.check {
				return
			}
			if got.ObjectID == "" || got.ID == "" ||
				got.UserID != user.ID ||
				got.MainTag != tt.input.MainTag ||
				!cmp.Equal(got.SubTags, *(tt.input.SubTags)) ||
				got.Content != tt.input.Content {
				t.Errorf("NewThread() = %v, input = %v", got, tt.input)
			}
			if tt.input.Title != nil && *(tt.input.Title) != "" {
				if got.Title == "" {
					t.Errorf("NewThread() = %v, should have title", got)
				}
			} else {
				if got.Title != "" {
					t.Errorf("NewThread() = %v, shouldn't have title", got)
				}
			}
			if !tt.input.Anonymous {
				if got.Anonymous == true || got.Author != user.Name {
					t.Errorf("NewThread() = %v, input = %v", got, tt.input)
				}
			} else {
				if got.Anonymous == false || got.Author == user.Name {
					t.Errorf("NewThread() = %v, input = %v", got, tt.input)
				}
			}
		})
	}
}

func TestGetThreadsByTags(t *testing.T) {
	threads := []*Thread{}
	user := mockUsers[1]
	ctx := ctxWithUser(user)
	for i := 0; i != 20; i++ {
		subTags := []string{}
		if i%2 == 0 {
			subTags = append(subTags, "2")
		}
		if i%3 == 0 {
			subTags = append(subTags, "3")
		}
		input := &ThreadInput{
			Anonymous: true,
			Content:   "content",
			MainTag:   mgmt.Config.MainTags[0],
			SubTags:   &subTags,
		}
		thread, err := NewThread(ctx, input)
		if err != nil {
			t.Fatal(errors.Wrap(err, "create thread"))
		}
		threads = append(threads, thread)
	}
	type args struct {
		tags []string
		sq   *SliceQuery
	}
	tests := []struct {
		name    string
		args    args
		want    []*Thread
		want1   *SliceInfo
		wantErr bool
	}{
		{
			"find tag 1", args{[]string{"1"}, &SliceQuery{Limit: 3}},
			[]*Thread{}, &SliceInfo{"", ""}, false,
		},
		{
			"find tag 2", args{[]string{"2"}, &SliceQuery{Limit: 3, Desc: true}},
			[]*Thread{threads[18], threads[16], threads[14]},
			&SliceInfo{threads[18].genCursor(), threads[14].genCursor()}, false,
		},
		{
			"find tag 3", args{[]string{"3"}, &SliceQuery{Limit: 3, Desc: true}},
			[]*Thread{threads[18], threads[15], threads[12]},
			&SliceInfo{threads[18].genCursor(), threads[12].genCursor()}, false,
		},
		{
			"find tag 3 desc", args{[]string{"3"},
				&SliceQuery{Limit: 3, Desc: true, Cursor: threads[12].genCursor()}},
			[]*Thread{threads[9], threads[6], threads[3]},
			&SliceInfo{threads[9].genCursor(), threads[3].genCursor()}, false,
		},
		{
			"find tag 3 asc", args{[]string{"3"},
				&SliceQuery{Limit: 3, Cursor: threads[9].genCursor()}},
			[]*Thread{threads[18], threads[15], threads[12]},
			&SliceInfo{threads[18].genCursor(), threads[12].genCursor()}, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GetThreadsByTags(ctx, tt.args.tags, tt.args.sq)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetThreadsByTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("GetThreadsByTags() got = %v, want %v", got, tt.want)
				return
			}
			for i := 0; i < len(got); i++ {
				if got[i].ID != tt.want[i].ID {
					t.Errorf("GetThreadsByTags() got = %v, want %v", got[i].ID, tt.want[i].ID)
					return
				}
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GetThreadsByTags() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	{
		thread, err := FindThread(ctx, threads[0].ID)
		if err != nil {
			t.Errorf("FindThread(%v) error = %v", thread, err)
		}
		if thread.ID != threads[0].ID {
			t.Errorf("FindThread(%v) got %v", threads[0].ID, thread.ID)
		}
	}
	{
		thread, err := FindThread(ctx, "AA")
		if err == nil {
			t.Errorf("FindThread(%v) should be error, found %v", err, thread)
		}
	}
}

func TestThread_GetReplies(t *testing.T) {
	user := mockUsers[1]
	ctx := ctxWithUser(user)
	input := &ThreadInput{
		Content:   "content",
		Anonymous: true,
		MainTag:   mgmt.Config.MainTags[0],
	}
	thread, err := NewThread(ctx, input)
	if err != nil {
		t.Errorf("FindThread(%v) should be error, found %v", err, thread)
	}
	posts := []*Post{}
	for i := 0; i < 6; i++ {
		pInput := &PostInput{
			ThreadID:  thread.ID,
			Anonymous: true,
			Content:   "post",
		}
		post, err := NewPost(ctx, pInput)
		if err != nil {
			t.Fatalf("new post error: %v", err)
		}
		posts = append(posts, post)
	}

	tests := []struct {
		name    string
		sq      *SliceQuery
		want    []*Post
		want1   *SliceInfo
		wantErr bool
	}{
		{"first 3", &SliceQuery{Limit: 3}, []*Post{posts[0], posts[1], posts[2]},
			&SliceInfo{posts[0].ObjectID.Hex(), posts[2].ObjectID.Hex()}, false},
		{"3 after 3", &SliceQuery{Limit: 3, Cursor: posts[2].ObjectID.Hex()},
			[]*Post{posts[3], posts[4], posts[5]},
			&SliceInfo{posts[3].ObjectID.Hex(), posts[5].ObjectID.Hex()}, false},
		{"3 before 3", &SliceQuery{Limit: 3, Desc: true, Cursor: posts[3].ObjectID.Hex()},
			[]*Post{posts[0], posts[1], posts[2]},
			&SliceInfo{posts[0].ObjectID.Hex(), posts[2].ObjectID.Hex()}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := thread.GetReplies(ctx, tt.sq)
			if (err != nil) != tt.wantErr {
				t.Errorf("Thread.GetReplies() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("Thread.GetReplies() got = %v, want %v", got, tt.want)
			}
			for i := 0; i < len(got); i++ {
				if got[i].ID != tt.want[i].ID {
					t.Errorf("Thread.GetReplies () got = %v, want %v", got[i].ID, tt.want[i].ID)
				}
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Thread.GetReplies() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
