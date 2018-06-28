package model

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
)

func TestPost(t *testing.T) {
	account := mockAccounts[2]
	ctx := ctxWithAccount(account)
	thread, err := NewThread(ctx, &ThreadInput{
		Content: "thread!", MainTag: pkg.mainTags[0],
	})
	if err != nil {
		t.Fatal(errors.Wrap(err, "create thread"))
	}
	if err := account.AddName(ctx, "testPost"); err != nil {
		t.Fatal(errors.Wrap(err, "add name"))
	}

	t.Log("Post1, normal post, signed name")
	input1 := &PostInput{
		ThreadID: thread.ID,
		Author:   &(account.Names[0]),
		Content:  "post1",
	}
	post1, err := NewPost(ctx, input1)
	if err != nil {
		t.Fatal(errors.Wrap(err, "create post1"))
	}
	if post1.ObjectID == "" || post1.ID == "" || post1.Anonymous == true ||
		post1.Author != account.Names[0] || post1.AccountID != account.ID ||
		post1.ThreadID != thread.ID || post1.Content != input1.Content ||
		len(post1.Refers) != 0 {
		t.Fatal(errors.Errorf("Post1 wrong! get: %v", post1))
	}

	t.Log("Post2, Anonymous Post")
	input2 := &PostInput{
		ThreadID: thread.ID,
		Author:   nil,
		Content:  "post2",
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
			"In one thread, AnonymousID of one account must be same, want %s, get %s",
			thread.Author, post2.Author,
		))
	}

	t.Log("Post3, has refers")
	input3 := &PostInput{
		ThreadID: thread.ID,
		Author:   nil,
		Content:  "post3",
		Refers:   &[]string{post1.ID, post2.ID},
	}
	post3, err := NewPost(ctx, input3)
	if err != nil {
		t.Fatal(errors.Wrap(err, "create post2"))
	}
	if !reflect.DeepEqual(post3.Refers, *(input3.Refers)) {
		t.Fatalf(
			"Post 3 refers error: %v, want %v", post3.Refers, input3.Refers,
		)
	}
	refers, err := post3.ReferPosts(ctx)
	if err != nil {
		t.Fatal(errors.Wrap(err, "find refer posts"))
	}
	if cmp.Equal(post1, refers[0], strSliceCmp) {
		t.Fatalf("refers[0] = %v, want %v", refers[0], post1)
	}
	if cmp.Equal(post2, refers[1], strSliceCmp) {
		t.Fatalf("refers[0] = %v, want %v", refers[0], post2)
	}

	t.Log("Find post")
	post4, err := FindPost(ctx, post1.ID)
	if err != nil {
		t.Fatal(errors.Wrap(err, "find post"))
	}
	if cmp.Equal(post1, post4, strSliceCmp) {
		t.Fatalf("FindPost() = %v, want %v", post4, post1)
	}
}
