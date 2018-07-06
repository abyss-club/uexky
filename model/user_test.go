package model

import (
	"context"
	"log"
	"reflect"
	"testing"

	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
)

var mockUsers []*User

func addMockUser() {
	log.Print("addMockUser!")
	users := []*User{
		&User{bson.NewObjectId(), "0@mail.com", "test0", []string{"动画"}},
		&User{bson.NewObjectId(), "1@mail.com", "", []string{}},
		&User{bson.NewObjectId(), "2@mail.com", "", []string{}},
	}
	c, cs := Colle("users")
	defer cs()
	for _, user := range users {
		if err := c.Insert(user); err != nil {
			log.Fatal(errors.Wrap(err, "gen mock users"))
		}
	}
	mockUsers = users
}

func ctxWithUser(a *User) context.Context {
	ctx := context.Background()
	return context.WithValue(ctx, ContextLoggedInUser, a.ID)
}

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
		{"test invalid token", args{ctxWithUser(&User{ID: "?"})}, nil, true},
		{"test unexist token", args{ctxWithUser(&User{ID: bson.NewObjectId()})}, nil, true},
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
		{"new user", args{
			context.Background(), "3@mail.com",
		}, "3@mail.com", false},
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

func TestUser_SyncTags(t *testing.T) {
	t.Skip("skip due to inconsistency")
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
			ctx := context.Background()
			if err := tt.args.user.SyncTags(ctx, tt.args.tags); (err != nil) != tt.wantErr {
				t.Errorf("User.SyncTags() error = %v, wantErr %v", err, tt.wantErr)
			}
			ctx = ctxWithUser(tt.args.user)
			a, err := GetUser(ctx)
			if err != nil {
				t.Error(errors.Wrap(err, "User.AddName() get user error"))
			}
			if !reflect.DeepEqual(a.Tags, tt.want) || !reflect.DeepEqual(tt.args.user.Tags, tt.want) {
				t.Errorf("User.AddName() want = %v, in memory = %v, in db = %v",
					tt.want, tt.args.user.Tags, a.Tags)
			}
		})
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
			got, err := tt.args.user.AnonymousID(tt.args.threadID, tt.args.new)
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
