package model

import (
	"reflect"
	"testing"

	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky/mgmt"
)

func TestPost(t *testing.T) {
	user := mockUsers[0]
	ctx := ctxWithUser(user)
	thread, err := NewThread(ctx, &ThreadInput{
		Content: "thread!", MainTag: mgmt.Config.MainTags[0], Anonymous: true,
	})
	if err != nil {
		t.Fatal(errors.Wrap(err, "create thread"))
	}

	t.Log("Post1, normal post, signed name")
	input1 := &PostInput{
		ThreadID:  thread.ID,
		Anonymous: false,
		Content:   "post1",
	}
	post1, err := NewPost(ctx, input1)
	if err != nil {
		t.Fatal(errors.Wrap(err, "create post1"))
	}
	if post1.ObjectID == "" || post1.ID == "" || post1.Anonymous != false ||
		post1.Author != user.Name || post1.UserID != user.ID ||
		post1.ThreadID != thread.ID || post1.Content != input1.Content ||
		len(post1.Quotes) != 0 {
		t.Fatal(errors.Errorf("Post1 wrong! get: %+v, input = %+v, user = %+v", post1, input1, user))
	}

	t.Log("Post2, Anonymous Post")
	input2 := &PostInput{
		ThreadID:  thread.ID,
		Anonymous: true,
		Content:   "post2",
	}
	post2, err := NewPost(ctx, input2)
	if err != nil {
		t.Fatal(errors.Wrap(err, "create post2"))
	}
	if post2.Anonymous == false || post2.Author == "" {
		t.Fatal(errors.Errorf("Post2 wrong! get: %v", post2))
	}
	if post2.Author != thread.Author {
		t.Fatal(errors.Errorf(
			"In one thread, AnonymousID of one user must be same, want %s, get %s",
			thread.Author, post2.Author,
		))
	}

	t.Log("Post3, has quotes")
	input3 := &PostInput{
		ThreadID:  thread.ID,
		Anonymous: true,
		Content:   "post3",
		Quotes:    &[]string{post1.ID, post2.ID},
	}
	post3, err := NewPost(ctx, input3)
	if err != nil {
		t.Fatal(errors.Wrap(err, "create post2"))
	}
	if !reflect.DeepEqual(post3.Quotes, *(input3.Quotes)) {
		t.Fatalf(
			"Post 3 quotes error: %v, want %v", post3.Quotes, input3.Quotes,
		)
	}
	quotes, err := post3.QuotePosts(ctx)
	if err != nil {
		t.Fatal(errors.Wrap(err, "find quote posts"))
	}
	if !equal(post1, quotes[0]) {
		t.Fatalf("quotes[0] = %v, want = %v", quotes[0], post1)
	}
	if !equal(post2, quotes[1]) {
		t.Fatalf("quotes[0] = %v, want %v", quotes[0], post2)
	}

	t.Log("Check quoted count of Post1")
	c, err := post1.QuoteCount(ctx)
	if err != nil {
		t.Errorf("Post.CountOfQuoted() error = %v", err)
	}
	if c != 1 {
		t.Errorf("Post.CountOfQuoted() = %v, want %v", c, 1)
	}

	t.Log("Find post")
	post4, err := FindPost(ctx, post3.ID)
	if err != nil {
		t.Fatal(errors.Wrap(err, "find post"))
	}
	if !equal(post3, post4) {
		t.Fatalf("FindPost() = %v, want %v", post4, post3)
	}

	t.Log("checkout thread update time")
	nThread, err := FindThread(ctx, thread.ID)
	if err != nil {
		t.Fatal(errors.Wrap(err, "find thread"))
	}
	if nThread.UpdateTime != post4.CreateTime {
		t.Fatalf("Checkout Thread UpdateTime = %v, want %v",
			nThread.UpdateTime, post4.CreateTime)
	}
}
