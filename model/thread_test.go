package model

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/globalsign/mgo/bson"
)

func TestThread(t *testing.T) {
	input1 := &ThreadInput{
		Author:  &(mockAccounts[1].Names[0]),
		Content: "Thread1",
		MainTag: mainTags[0],
		SubTags: *[]string{},
	}
	type args struct {
		ctx   context.Context
		input *ThreadInput
	}
	tests := []struct {
		name    string
		args    args
		want    *Thread
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewThread(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewThread() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewThread() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetThreadsByTags(t *testing.T) {
	type args struct {
		ctx  context.Context
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GetThreadsByTags(tt.args.ctx, tt.args.tags, tt.args.sq)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetThreadsByTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetThreadsByTags() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GetThreadsByTags() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestFindThread(t *testing.T) {
	type args struct {
		ctx context.Context
		ID  string
	}
	tests := []struct {
		name    string
		args    args
		want    *Thread
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindThread(tt.args.ctx, tt.args.ID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindThread() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindThread() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isThreadExist(t *testing.T) {
	type args struct {
		threadID string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := isThreadExist(tt.args.threadID)
			if (err != nil) != tt.wantErr {
				t.Errorf("isThreadExist() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("isThreadExist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestThread_GetReplies(t *testing.T) {
	type fields struct {
		ObjectID   bson.ObjectId
		ID         string
		Anonymous  bool
		Author     string
		Account    string
		CreateTime time.Time
		MainTag    string
		SubTags    []string
		Title      string
		Content    string
	}
	type args struct {
		ctx context.Context
		sq  *SliceQuery
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*Post
		want1   *SliceInfo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t := &Thread{
				ObjectID:   tt.fields.ObjectID,
				ID:         tt.fields.ID,
				Anonymous:  tt.fields.Anonymous,
				Author:     tt.fields.Author,
				Account:    tt.fields.Account,
				CreateTime: tt.fields.CreateTime,
				MainTag:    tt.fields.MainTag,
				SubTags:    tt.fields.SubTags,
				Title:      tt.fields.Title,
				Content:    tt.fields.Content,
			}
			got, got1, err := t.GetReplies(tt.args.ctx, tt.args.sq)
			if (err != nil) != tt.wantErr {
				t.Errorf("Thread.GetReplies() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Thread.GetReplies() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Thread.GetReplies() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
