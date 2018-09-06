package model

import (
	"context"
	"reflect"
	"testing"

	"github.com/pkg/errors"
)

var mockUsers []*User

func TestGetUser(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *User
		wantErr bool
	}{
		{"normal", args{ctxWithUser(mockUsers[0])}, mockUsers[0], false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUser(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				return
			}
			if !reflect.DeepEqual(got.ID, tt.want.ID) {
				t.Errorf("GetUser() ID = %+v, want %+v", got, tt.want)
			}
			if !reflect.DeepEqual(got.Email, tt.want.Email) {
				t.Errorf("GetUser() Token = %+v, want %+v", got, tt.want)
			}
			if !reflect.DeepEqual(got.Name, tt.want.Name) {
				t.Errorf("GetUser() Names = %+v, want %+v", got, tt.want)
			}
			if !reflect.DeepEqual(got.Tags, tt.want.Tags) {
				t.Errorf("GetUser() Tags = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestGetUserByEmail(t *testing.T) {
	type args struct {
		ctx   context.Context
		email string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"exist user", args{
			ctxWithUser(mockUsers[0]), mockUsers[0].Email,
		}, mockUsers[0].Email, false},
		{"new user", args{testCtx, "3@mail.com"}, "3@mail.com", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUserByEmail(tt.args.ctx, tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Email, tt.want) {
				t.Errorf("GetUserByEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_SetName(t *testing.T) {
	mockUsers[0].Name = "test0"
	mockUsers[1].Name = ""
	mockUsers[2].Name = ""
	type args struct {
		user *User
		name string
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		wantName string
	}{
		{"has name", args{mockUsers[0], "test1"}, true, "test0"},
		{"no name", args{mockUsers[1], "testX"}, false, "testX"},
		{"same name", args{mockUsers[2], "testX"}, true, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := ctxWithUser(tt.args.user)
			if err := tt.args.user.SetName(ctx, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("User.AddName() error = %v, wantErr %v", err, tt.wantErr)
			}
			a, err := GetUser(ctx)
			if err != nil {
				t.Error(errors.Wrap(err, "User.AddName() get user error"))
			}
			if a.Name != tt.wantName {
				t.Errorf("User.AddName() want = %s, in memory = %s, in db = %s",
					tt.wantName, tt.args.user.Name, a.Name)
			}
		})
	}
}

func parseToStringSet(sList []string) map[string]bool {
	set := map[string]bool{}
	for _, s := range sList {
		set[s] = true
	}
	return set
}

func isSet(s []string) bool {
	set := parseToStringSet(s)
	return len(set) == len(s)
}

func cmpTags(lTags []string, rTags []string) bool {
	lt := parseToStringSet(lTags)
	rt := parseToStringSet(rTags)
	if len(lt) != len(rt) {
		return false
	}
	for s := range lt {
		_, exists := rt[s]
		if !exists {
			return false
		}
	}
	return true
}

func TestUser_SyncTags(t *testing.T) {
	type args struct {
		user *User
		tags []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    []string
	}{
		{"add tag", args{mockUsers[1], []string{"tag1"}}, false, []string{"tag1"}},
		{"delete tag", args{mockUsers[1], []string{}}, false, []string{}},
		{"add tag with repeated", args{mockUsers[1], []string{"tag1", "tag2", "tag1"}}, false, []string{
			"tag1", "tag2"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := ctxWithUser(tt.args.user)
			if err := tt.args.user.SyncTags(ctx, tt.args.tags); (err != nil) != tt.wantErr {
				t.Errorf("User.SyncTags() error = %v, wantErr %v", err, tt.wantErr)
			}
			ctx = ctxWithUser(tt.args.user)
			a, err := GetUser(ctx)
			if err != nil {
				t.Error(errors.Wrap(err, "User.AddName() get user error"))
			}
			if !isSet(a.Tags) {
				t.Errorf("User.SyncTags() error, repeated tag, %q", a.Tags)
			}
			if !cmpTags(a.Tags, tt.want) || !cmpTags(tt.args.user.Tags, tt.want) {
				t.Errorf("User.AddName() want = %v, in memory = %v, in db = %v",
					tt.want, tt.args.user.Tags, a.Tags)
			}
		})
	}
}

func TestUser_AddSubbedTags(t *testing.T) {
	user := mockUsers[2]
	ctx := ctxWithUser(user)

	t.Log("reset tags subscribed")
	{
		if err := user.SyncTags(ctx, []string{"A", "B", "C"}); err != nil {
			t.Fatalf("reset tags error: %v", err)
		}
	}
	// Tags: A, B, C
	t.Log("test add tags")
	{
		want := []string{"A", "B", "C", "D", "E"}
		if err := user.AddSubbedTags(ctx, []string{"B", "B", "D", "E"}); err != nil {
			t.Fatalf("AddSubbedTags() error: %v", err)
		}
		u, err := GetUser(ctx)
		if err != nil {
			t.Fatalf("GetUser() error: %v", err)
		}
		if !cmpTags(u.Tags, user.Tags) || !cmpTags(u.Tags, want) {
			t.Fatalf("AddSubbedTags() want %q, in memory = %q, in db = %q",
				want, user.Tags, u.Tags)
		}
	}
	// Tags: A, B, C, D, E
	t.Log("test add tags")
	{
		want := []string{"A", "C"}
		if err := user.DelSubbedTags(ctx, []string{"B", "B", "D", "E"}); err != nil {
			t.Fatalf("AddSubbedTags() error: %v", err)
		}
		u, err := GetUser(ctx)
		if err != nil {
			t.Fatalf("GetUser() error: %v", err)
		}
		if !cmpTags(u.Tags, user.Tags) || !cmpTags(u.Tags, want) {
			t.Fatalf("AddSubbedTags() want %q, in memory = %q, in db = %q",
				want, user.Tags, u.Tags)
		}
	}
}

func TestUser_AnonymousID(t *testing.T) {
	type args struct {
		user     *User
		threadID string
		new      bool
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		equalLast bool
	}{
		{"new thread", args{mockUsers[0], "Thread1", false}, false, false},
		{"same thread", args{mockUsers[0], "Thread1", false}, false, true},
		{"same thread, renew id", args{mockUsers[0], "Thread1", true}, false, false},
	}
	last := ""
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.user.AnonymousID(testCtx, tt.args.threadID, tt.args.new)
			if (err != nil) != tt.wantErr {
				t.Errorf("User.AnonymousID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("get AnonymousID '%s'", got)
			if (last == got) != tt.equalLast {
				t.Errorf("User.AnonymousID() = %v, last %v, want equal %v", got, last, tt.equalLast)
			}
			last = got
		})
	}
}
