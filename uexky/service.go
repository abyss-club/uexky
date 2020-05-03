package uexky

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"
	"gitlab.com/abyss.club/uexky/lib/uid"
	"gitlab.com/abyss.club/uexky/uexky/entity"
)

type Repo interface{}

type Service struct {
	Repo  Repo
	User  *entity.UserService
	Forum *entity.ForumService
	Noti  *entity.NotiService
}

func (s *Service) SignInByEmail(ctx context.Context, email string) (bool, error) {
	return s.User.SignInByEmail(ctx, email)
}

func (s *Service) Profile(ctx context.Context) (*entity.User, error) {
	return s.User.GetUserRequirePermission(ctx, entity.ActionProfile)
}

func (s *Service) SetUserName(ctx context.Context, name string) (*entity.User, error) {
	user, err := s.User.GetUserRequirePermission(ctx, entity.ActionProfile)
	if err != nil {
		return nil, err
	}
	err = user.SetName(ctx, name)
	return user, err
}

func (s *Service) GetUserThreads(
	ctx context.Context, obj *entity.User, query entity.SliceQuery,
) (*entity.ThreadSlice, error) {
	user, err := s.User.GetUserRequirePermission(ctx, entity.ActionProfile)
	if err != nil {
		return nil, err
	}
	if obj == nil || obj.Email != user.Email {
		return nil, errors.New("permission denied")
	}
	return s.Forum.GetUserThreads(ctx, user, query)
}

func (s *Service) GetUserPosts(ctx context.Context, obj *entity.User, query entity.SliceQuery) (*entity.PostSlice, error) {
	user, err := s.User.GetUserRequirePermission(ctx, entity.ActionProfile)
	if err != nil {
		return nil, err
	}
	if obj == nil || obj.Email != user.Email {
		return nil, errors.New("permission denied")
	}
	return s.Forum.GetUserPosts(ctx, user, query)
}

func (s *Service) GetUserTags(ctx context.Context, obj *entity.User) ([]string, error) {
	user, err := s.User.GetUserRequirePermission(ctx, entity.ActionProfile)
	if err != nil {
		return nil, err
	}
	if obj == nil || obj.Email != user.Email {
		return nil, errors.New("permission denied")
	}
	return s.Forum.GetUserTags(ctx, user)
}

func (s *Service) SyncUserTags(ctx context.Context, tags []*string) (*entity.User, error) {
	user, err := s.User.GetUserRequirePermission(ctx, entity.ActionProfile)
	if err != nil {
		return nil, err
	}
	return user, s.Forum.SyncUserTags(ctx, user, tags)
}

func (s *Service) AddUserSubbedTag(ctx context.Context, tag string) (*entity.User, error) {
	user, err := s.User.GetUserRequirePermission(ctx, entity.ActionProfile)
	if err != nil {
		return nil, err
	}
	return user, s.Forum.AddUserSubbedTag(ctx, user, tag)
}

func (s *Service) DelUserSubbedTag(ctx context.Context, tag string) (*entity.User, error) {
	user, err := s.User.GetUserRequirePermission(ctx, entity.ActionProfile)
	if err != nil {
		return nil, err
	}
	return user, s.Forum.DelUserSubbedTag(ctx, user, tag)
}

func (s *Service) BanUser(ctx context.Context, postID *uid.UID, threadID *uid.UID) (bool, error) {
	user, err := s.User.GetUserRequirePermission(ctx, entity.ActionBanUser)
	if err != nil {
		return false, err
	}
	if postID != nil {
		post, err := s.Forum.GetPostByID(ctx, *postID)
		if err != nil {
			return false, err
		}
		return user.BanUser(ctx, post.Author, post.Anonymous)
	} else if threadID != nil {
		thread, err := s.Forum.GetThreadByID(ctx, *threadID)
		if err != nil {
			return false, err
		}
		return user.BanUser(ctx, thread.Author, thread.Anonymous)
	}
	return false, errors.New("must specified post id or thread id")
}

func (s *Service) BlockPost(ctx context.Context, postID uid.UID) (*entity.Post, error) {
	_, err := s.User.GetUserRequirePermission(ctx, entity.ActionBlockPost)
	if err != nil {
		return nil, err
	}
	post, err := s.Forum.GetPostByID(ctx, postID)
	if err != nil {
		return nil, err
	}
	return post, post.Block(ctx)
}

func (s *Service) LockThread(ctx context.Context, threadID uid.UID) (*entity.Thread, error) {
	_, err := s.User.GetUserRequirePermission(ctx, entity.ActionLockThread)
	if err != nil {
		return nil, err
	}
	thread, err := s.Forum.GetThreadByID(ctx, threadID)
	if err != nil {
		return nil, err
	}
	return thread, thread.Lock(ctx)
}

func (s *Service) BlockThread(ctx context.Context, threadID uid.UID) (*entity.Thread, error) {
	_, err := s.User.GetUserRequirePermission(ctx, entity.ActionBlockThread)
	if err != nil {
		return nil, err
	}
	thread, err := s.Forum.GetThreadByID(ctx, threadID)
	if err != nil {
		return nil, err
	}
	return thread, thread.Block(ctx)
}

func (s *Service) EditTags(
	ctx context.Context, threadID uid.UID, mainTag string, subTags []string,
) (*entity.Thread, error) {
	_, err := s.User.GetUserRequirePermission(ctx, entity.ActionEditTag)
	if err != nil {
		return nil, err
	}
	thread, err := s.Forum.GetThreadByID(ctx, threadID)
	if err != nil {
		return nil, err
	}
	return thread, thread.EditTags(ctx, mainTag, subTags)
}

func (s *Service) PubThread(ctx context.Context, thread entity.ThreadInput) (*entity.Thread, error) {
	return s.Forum.NewThread(ctx, thread)
}

func (s *Service) SearchThreads(
	ctx context.Context, tags []string, query entity.SliceQuery,
) (*entity.ThreadSlice, error) {
	return s.Forum.SearchThreads(ctx, tags, query)
}

func (s *Service) GetThreadByID(ctx context.Context, id uid.UID) (*entity.Thread, error) {
	return s.Forum.GetThreadByID(ctx, id)
}

func (s *Service) PubPost(ctx context.Context, post entity.PostInput) (*entity.Post, error) {
	res, err := s.Forum.NewPost(ctx, post)
	if err != nil {
		return nil, err
	}
	go func() {
		ctx := context.Background()
		if err := s.Noti.NewRepliedNoti(ctx, res.Thread, res.Post); err != nil {
			log.Error(err)
		}
		for _, qp := range res.QuotedPost {
			if err := s.Noti.NewQuotedNoti(ctx, res.Thread, res.Post, qp); err != nil {
				log.Error(err)
			}
		}
	}()
	return res.Post, nil
}

func (s *Service) GetPostByID(ctx context.Context, id uid.UID) (*entity.Post, error) {
	return s.Forum.GetPostByID(ctx, id)
}

func (s *Service) GetMainTags(ctx context.Context) ([]string, error) {
	return s.Forum.GetMainTags(ctx)
}

func (s *Service) GetRecommendedTags(ctx context.Context) ([]string, error) {
	return s.Forum.GetRecommendedTags(ctx)
}

func (s *Service) SearchTags(ctx context.Context, query *string, limit *int) ([]*entity.Tag, error) {
	return s.Forum.SearchTags(ctx, query, limit)
}

func (s *Service) GetUnreadNotiCount(ctx context.Context) (*entity.UnreadNotiCount, error) {
	return s.Noti.GetUnreadNotiCount(ctx)
}

func (s *Service) GetNotification(
	ctx context.Context, typeArg string, query entity.SliceQuery,
) (*entity.NotiSlice, error) {
	return s.Noti.GetNotification(ctx, typeArg, query)
}
