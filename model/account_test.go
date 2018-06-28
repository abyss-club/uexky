package model

import (
	"context"
	"log"
	"reflect"
	"testing"

	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
)

var mockAccounts []*Account

func addMockUser() {
	accounts := []*Account{
		&Account{bson.NewObjectId(), "0@mail.com", []string{"test0"},
			[]string{"动画"}},
		&Account{bson.NewObjectId(), "1@mail.com", []string{"test1"},
			[]string{}},
		&Account{bson.NewObjectId(), "2@mail.com", []string{"test2"},
			[]string{}},
	}
	c, cs := Colle("accounts")
	defer cs()
	for _, account := range accounts {
		if err := c.Insert(account); err != nil {
			log.Fatal(errors.Wrap(err, "gen mock accounts"))
		}
	}
	mockAccounts = accounts
}

func ctxWithAccount(a *Account) context.Context {
	ctx := context.Background()
	return context.WithValue(ctx, ContextLoggedInAccount, a.ID)
}

func TestGetAccount(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *Account
		wantErr bool
	}{
		{"normal", args{ctxWithAccount(mockAccounts[0])}, mockAccounts[0], false},
		{"test invalid token", args{ctxWithAccount(&Account{ID: "?"})}, nil, true},
		{"test unexist token", args{ctxWithAccount(&Account{ID: bson.NewObjectId()})}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAccount(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				return
			}
			if !reflect.DeepEqual(got.ID, tt.want.ID) {
				t.Errorf("GetAccount() ID = %+v, want %+v", got, tt.want)
			}
			if !reflect.DeepEqual(got.Email, tt.want.Email) {
				t.Errorf("GetAccount() Token = %+v, want %+v", got, tt.want)
			}
			if !reflect.DeepEqual(got.Names, tt.want.Names) {
				t.Errorf("GetAccount() Names = %+v, want %+v", got, tt.want)
			}
			if !reflect.DeepEqual(got.Tags, tt.want.Tags) {
				t.Errorf("GetAccount() Tags = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestGetAccountByEmail(t *testing.T) {
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
		{"exist account", args{
			ctxWithAccount(mockAccounts[0]), mockAccounts[0].Email,
		}, mockAccounts[0].Email, false},
		{"new account", args{
			context.Background(), "3@mail.com",
		}, "3@mail.com", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAccountByEmail(tt.args.ctx, tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAccountByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Email, tt.want) {
				t.Errorf("GetAccountByEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccount_AddName(t *testing.T) {
	type args struct {
		account *Account
		name    string
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantNames []string
	}{
		{"new name", args{mockAccounts[0], "testX"}, false, []string{"test0", "testX"}},
		{"same name", args{mockAccounts[0], "test1"}, true, []string{"test0", "testX"}},
		{"same name", args{mockAccounts[0], "testX1"}, false, []string{
			"test0", "testX", "testX1"}},
		{"same name", args{mockAccounts[0], "testX2"}, false, []string{
			"test0", "testX", "testX1", "testX2"}},
		{"same name", args{mockAccounts[0], "testX3"}, false, []string{
			"test0", "testX", "testX1", "testX2", "testX3"}},
		{"same name", args{mockAccounts[0], "testX4"}, true, []string{
			"test0", "testX", "testX1", "testX2", "testX3"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := ctxWithAccount(tt.args.account)
			if err := tt.args.account.AddName(ctx, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("Account.AddName() error = %v, wantErr %v", err, tt.wantErr)
			}
			a, err := GetAccount(ctx)
			if err != nil {
				t.Error(errors.Wrap(err, "Account.AddName() get account error"))
			}
			if a == nil {
				return
			}
			if !reflect.DeepEqual(a.Names, tt.wantNames) || !reflect.DeepEqual(tt.args.account.Names, tt.wantNames) {
				t.Errorf("Account.AddName() want = %v, in memory = %v, in db = %v",
					tt.wantNames, tt.args.account.Names, a.Names)
			}
		})
	}
}

func TestAccount_SyncTags(t *testing.T) {
	type args struct {
		account *Account
		tags    []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    []string
	}{
		{"add tag", args{mockAccounts[1], []string{"tag1"}}, false, []string{"tag1"}},
		{"delete tag", args{mockAccounts[1], []string{}}, false, []string{}},
		{"add tag with repeated", args{mockAccounts[1], []string{"tag1", "tag2", "tag1"}}, false, []string{
			"tag1", "tag2"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if err := tt.args.account.SyncTags(ctx, tt.args.tags); (err != nil) != tt.wantErr {
				t.Errorf("Account.SyncTags() error = %v, wantErr %v", err, tt.wantErr)
			}
			ctx = ctxWithAccount(tt.args.account)
			a, err := GetAccount(ctx)
			if err != nil {
				t.Error(errors.Wrap(err, "Account.AddName() get account error"))
			}
			if !reflect.DeepEqual(a.Tags, tt.want) || !reflect.DeepEqual(tt.args.account.Tags, tt.want) {
				t.Errorf("Account.AddName() want = %v, in memory = %v, in db = %v",
					tt.want, tt.args.account.Tags, a.Tags)
			}
		})
	}
}

func TestAccount_HaveName(t *testing.T) {
	type args struct {
		account *Account
		name    string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"had name", args{mockAccounts[2], "test2"}, true},
		{"had name", args{mockAccounts[2], "test2X"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.account.HaveName(tt.args.name); got != tt.want {
				t.Errorf("Account.HaveName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccount_AnonymousID(t *testing.T) {
	type args struct {
		account  *Account
		threadID string
		new      bool
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		equalLast bool
	}{
		{"new thread", args{mockAccounts[0], "Thread1", false}, false, false},
		{"same thread", args{mockAccounts[0], "Thread1", false}, false, true},
		{"same thread, renew id", args{mockAccounts[0], "Thread1", true}, false, false},
	}
	last := ""
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.account.AnonymousID(tt.args.threadID, tt.args.new)
			if (err != nil) != tt.wantErr {
				t.Errorf("Account.AnonymousID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("get AnonymousID '%s'", got)
			if (last == got) != tt.equalLast {
				t.Errorf("Account.AnonymousID() = %v, last %v, want equal %v", got, last, tt.equalLast)
			}
			last = got
		})
	}
}
