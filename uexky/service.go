package uexky

import (
	"context"
	"fmt"

	"gitlab.com/abyss.club/uexky/lib/uid"
	"gitlab.com/abyss.club/uexky/uexky/entity"
)

type Repo interface{}

type Service struct {
	Repo  Repo `wire:"-"` // TODO
	User  *entity.UserService
	Forum *entity.ForumService
	Noti  *entity.NotiService
}

func (s *Service) SignInByEmail(ctx context.Context, email string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (s *Service) SetUserName(ctx context.Context, name string) (*entity.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (s *Service) SyncUserTags(ctx context.Context, tags []*string) (*entity.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (s *Service) AddUserSubbedTag(ctx context.Context, tag string) (*entity.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (s *Service) DelUserSubbedTag(ctx context.Context, tag string) (*entity.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (s *Service) BanUser(ctx context.Context, postID *uid.UID, threadID *uid.UID) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (s *Service) BlockPost(ctx context.Context, postID uid.UID) (*entity.Post, error) {
	panic(fmt.Errorf("not implemented"))
}

func (s *Service) LockThread(ctx context.Context, threadID uid.UID) (*entity.Thread, error) {
	panic(fmt.Errorf("not implemented"))
}

func (s *Service) BlockThread(ctx context.Context, threadID uid.UID) (*entity.Thread, error) {
	panic(fmt.Errorf("not implemented"))
}

func (s *Service) EditTags(
	ctx context.Context, threadID uid.UID, mainTag string, subTags []string,
) (*entity.Thread, error) {
	panic(fmt.Errorf("not implemented"))
}

func (s *Service) Profile(ctx context.Context) (*entity.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (s *Service) PubThread(ctx context.Context, thread entity.ThreadInput) (*entity.Thread, error) {
	panic(fmt.Errorf("not implemented"))
}

func (s *Service) SearchThreads(
	ctx context.Context, tags []string, query entity.SliceQuery,
) (*entity.ThreadSlice, error) {
	panic(fmt.Errorf("not implemented"))
}

func (s *Service) GetThreadByID(ctx context.Context, id uid.UID) (*entity.Thread, error) {
	panic(fmt.Errorf("not implemented"))
}

func (s *Service) PubPost(ctx context.Context, post entity.PostInput) (*entity.Post, error) {
	panic(fmt.Errorf("not implemented"))
}

func (s *Service) GetPostByID(ctx context.Context, id uid.UID) (*entity.Post, error) {
	panic(fmt.Errorf("not implemented"))
}

func (s *Service) GetMainTags(ctx context.Context) ([]string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (s *Service) GetRecommendedTags(ctx context.Context) ([]string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (s *Service) SearchTags(ctx context.Context, query *string, limit *int) ([]*entity.Tag, error) {
	panic(fmt.Errorf("not implemented"))
}

func (s *Service) GetUnreadNotiCount(ctx context.Context) (*entity.UnreadNotiCount, error) {
	panic(fmt.Errorf("not implemented"))
}

func (s *Service) GetNotification(
	ctx context.Context, typeArg string, query entity.SliceQuery,
) (*entity.NotiSlice, error) {
	panic(fmt.Errorf("not implemented"))
}
