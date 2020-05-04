// noti aggragate: SystemNoti, RepliedNoti, QuotedNoti

package entity

import (
	"context"
	"fmt"
	"time"
)

type NotiRepo interface{}

type NotiService struct {
	Repo NotiRepo
}

func (n *NotiService) GetUnreadNotiCount(ctx context.Context, user *User) (*UnreadNotiCount, error) {
	panic(fmt.Errorf("not implemented"))
}

func (n *NotiService) GetNotification(
	ctx context.Context, user *User, typeArg string, query SliceQuery,
) (*NotiSlice, error) {
	panic(fmt.Errorf("not implemented"))
}

type SystemNotiContent struct {
	Title       string `json:"title"`
	ContentText string `json:"content"`
}

type SystemNoti struct {
	ID        string            `json:"id"`
	Type      string            `json:"type"`
	EventTime time.Time         `json:"eventTime"`
	Content   SystemNotiContent `json:"content"`
}

func (n *NotiService) NewSystemNoti(ctx context.Context, title, context string) error {
	panic(fmt.Errorf("not implemented"))
}

func (n *SystemNoti) HasRead(ctx context.Context) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (n *SystemNoti) Title(ctx context.Context) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (n *SystemNoti) ContentText(ctx context.Context) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

type RepliedNotiContent struct {
	ThreadID string `json:"threadId"` // int64.toString()
}

type RepliedNoti struct {
	ID        string             `json:"id"`
	Type      string             `json:"type"`
	EventTime time.Time          `json:"eventTime"`
	Content   RepliedNotiContent `json:"content"`
}

func (n *NotiService) NewRepliedNoti(ctx context.Context, thread *Thread, reply *Post) error {
	panic(fmt.Errorf("not implemented"))
}

func (n *RepliedNoti) HasRead(ctx context.Context) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (n *RepliedNoti) Thread(ctx context.Context) (*Thread, error) {
	panic(fmt.Errorf("not implemented"))
}

func (n *RepliedNoti) Repliers(ctx context.Context) ([]string, error) {
	panic(fmt.Errorf("not implemented"))
}

type QuotedNotiContent struct {
	ThreadID string `json:"threadId"` // int64.toString()
	QuotedID string `json:"quotedId"` // int64.toString()
	PostID   string `json:"postId"`   // int64.toString()
}

type QuotedNoti struct {
	ID        string            `json:"id"`
	Type      string            `json:"type"`
	EventTime time.Time         `json:"eventTime"`
	Content   QuotedNotiContent `json:"content"`
}

func (n *NotiService) NewQuotedNoti(
	ctx context.Context, thread *Thread, post *Post, quotedPost *Post,
) error {
	panic(fmt.Errorf("not implemented"))
}

func (n *QuotedNoti) HasRead(ctx context.Context) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (n *QuotedNoti) Thread(ctx context.Context) (*Thread, error) {
	panic(fmt.Errorf("not implemented"))
}

func (n *QuotedNoti) QuotedPost(ctx context.Context) (*Post, error) {
	panic(fmt.Errorf("not implemented"))
}

func (n *QuotedNoti) Post(ctx context.Context) (*Post, error) {
	panic(fmt.Errorf("not implemented"))
}
