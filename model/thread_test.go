package model

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
)

func TestNewThread(t *testing.T) {
	user := mockUsers[1]
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
			&user.Name, "thread1", pkg.mainTags[0],
			&[]string{"tag1", "tag2"}, &titles[0],
		}, true, false},
		{"normal, anonymous, non-title", &ThreadInput{
			nil, "thread2", pkg.mainTags[0],
			&[]string{"tag1", "tag2"}, nil,
		}, true, false},
		{"error, no-main-tag", &ThreadInput{
			nil, "thread3", "em..", nil, nil,
		}, false, true},
		{"error, multi-main-tag", &ThreadInput{
			nil, "thread3", pkg.mainTags[0], &[]string{pkg.mainTags[1]}, nil,
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
			if tt.input.Author != nil && *(tt.input.Author) != "" {
				if got.Anonymous == true || got.Author != *(tt.input.Author) {
					t.Errorf("NewThread() = %v, input = %v", got, tt.input)
				}
			} else {
				if got.Anonymous == false || got.Author == "" {
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
			Content: "content",
			MainTag: pkg.mainTags[0],
			SubTags: &subTags,
		}
		thread, err := NewThread(ctx, input)
		if err != nil {
			t.Fatal(errors.Wrap(err, "create thread"))
		}
		threads = append(threads, thread)
	}
	t.Log("test for count")
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
			"find tag 1", args{[]string{pkg.mainTags[0], "1"}, &SliceQuery{Limit: 3}},
			[]*Thread{}, &SliceInfo{"", ""}, false,
		},
		{
			"find tag 2", args{[]string{pkg.mainTags[0], "2"}, &SliceQuery{Limit: 3}},
			[]*Thread{threads[18], threads[16], threads[14]},
			&SliceInfo{threads[18].ID, threads[14].ID}, false,
		},
		{
			"find tag 3", args{[]string{pkg.mainTags[0], "3"}, &SliceQuery{Limit: 3}},
			[]*Thread{threads[18], threads[15], threads[12]},
			&SliceInfo{threads[18].ID, threads[12].ID}, false,
		},
		{
			"find tag 3 before", args{[]string{pkg.mainTags[0], "3"},
				&SliceQuery{Limit: 3, Before: threads[12].ID}},
			[]*Thread{threads[9], threads[6], threads[3]},
			&SliceInfo{threads[9].ID, threads[3].ID}, false,
		},
		{
			"find tag 3 after", args{[]string{pkg.mainTags[0], "3"},
				&SliceQuery{Limit: 3, After: threads[12].ID}},
			[]*Thread{threads[18], threads[15]},
			&SliceInfo{threads[18].ID, threads[15].ID}, false,
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
			}
			for i := 0; i < len(got); i++ {
				if got[i].ID != tt.want[i].ID {
					t.Errorf("GetThreadsByTags() got = %v, want %v", got[i].ID, tt.want[i].ID)
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
		Content: "content",
		MainTag: pkg.mainTags[0],
	}
	thread, err := NewThread(ctx, input)
	if err != nil {
		t.Errorf("FindThread(%v) should be error, found %v", err, thread)
	}
	posts := []*Post{}
	for i := 0; i < 6; i++ {
		pInput := &PostInput{
			ThreadID: thread.ID,
			Content:  "post",
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
			&SliceInfo{posts[0].ID, posts[2].ID}, false},
		{"3 after 3", &SliceQuery{Limit: 3, After: posts[2].ID},
			[]*Post{posts[3], posts[4], posts[5]},
			&SliceInfo{posts[3].ID, posts[5].ID}, false},
		{"3 before 3", &SliceQuery{Limit: 3, Before: posts[3].ID},
			[]*Post{posts[0], posts[1], posts[2]},
			&SliceInfo{posts[0].ID, posts[2].ID}, false},
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
