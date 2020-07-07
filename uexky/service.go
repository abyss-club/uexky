package uexky

import (
	"context"
	"errors"

	"gitlab.com/abyss.club/uexky/lib/uid"
	"gitlab.com/abyss.club/uexky/uexky/adapter"
	"gitlab.com/abyss.club/uexky/uexky/entity"
)

type Service struct {
	User      *entity.UserService
	Forum     *entity.ForumService
	Noti      *entity.NotiService
	TxAdapter adapter.Tx
}

func (s *Service) TrySignInByEmail(ctx context.Context, email string) (entity.Code, error) {
	return s.User.TrySignInByEmail(ctx, email)
}

func (s *Service) SignInByCode(ctx context.Context, code string) (entity.Token, error) {
	return s.User.SignInByCode(ctx, code)
}

func (s *Service) CtxWithUserByToken(ctx context.Context, tok string) (context.Context, error) {
	return s.User.CtxWithUserByToken(ctx, tok)
}

func (s *Service) Profile(ctx context.Context) (*entity.User, error) {
	return s.User.RequirePermission(ctx, entity.ActionProfile)
}

func (s *Service) SetUserName(ctx context.Context, name string) (*entity.User, error) {
	user, err := s.User.RequirePermission(ctx, entity.ActionProfile)
	if err != nil {
		return nil, err
	}
	err = user.SetName(ctx, name)
	return user, err
}

func (s *Service) GetUserThreads(
	ctx context.Context, obj *entity.User, query entity.SliceQuery,
) (*entity.ThreadSlice, error) {
	user, err := s.User.RequirePermission(ctx, entity.ActionProfile)
	if err != nil {
		return nil, err
	}
	if obj == nil || obj.Email != user.Email {
		return nil, errors.New("permission denied")
	}
	return s.Forum.GetUserThreads(ctx, user, query)
}

func (s *Service) GetUserPosts(ctx context.Context, obj *entity.User, query entity.SliceQuery) (*entity.PostSlice, error) {
	user, err := s.User.RequirePermission(ctx, entity.ActionProfile)
	if err != nil {
		return nil, err
	}
	if obj == nil || obj.Email != user.Email {
		return nil, errors.New("permission denied")
	}
	return s.Forum.GetUserPosts(ctx, user, query)
}

func (s *Service) SyncUserTags(ctx context.Context, tags []string) (*entity.User, error) {
	user, err := s.User.RequirePermission(ctx, entity.ActionProfile)
	if err != nil {
		return nil, err
	}
	return user, user.SyncTags(ctx, user, tags)
}

func (s *Service) AddUserSubbedTag(ctx context.Context, tag string) (*entity.User, error) {
	user, err := s.User.RequirePermission(ctx, entity.ActionProfile)
	if err != nil {
		return nil, err
	}
	return user, user.AddSubbedTag(ctx, user, tag)
}

func (s *Service) DelUserSubbedTag(ctx context.Context, tag string) (*entity.User, error) {
	user, err := s.User.RequirePermission(ctx, entity.ActionProfile)
	if err != nil {
		return nil, err
	}
	return user, user.DelSubbedTag(ctx, user, tag)
}

func (s *Service) BanUser(ctx context.Context, postID *uid.UID, threadID *uid.UID) (bool, error) {
	_, err := s.User.RequirePermission(ctx, entity.ActionBanUser)
	if err != nil {
		return false, err
	}
	if postID != nil {
		post, err := s.Forum.GetPostByID(ctx, *postID)
		if err != nil {
			return false, err
		}
		return s.User.BanUser(ctx, post.Data.Author.UserID)
	} else if threadID != nil {
		thread, err := s.Forum.GetThreadByID(ctx, *threadID)
		if err != nil {
			return false, err
		}
		return s.User.BanUser(ctx, thread.AuthorObj.UserID)
	}
	return false, errors.New("must specified post id or thread id")
}

func (s *Service) BlockPost(ctx context.Context, postID uid.UID) (*entity.Post, error) {
	_, err := s.User.RequirePermission(ctx, entity.ActionBlockPost)
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
	_, err := s.User.RequirePermission(ctx, entity.ActionLockThread)
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
	_, err := s.User.RequirePermission(ctx, entity.ActionBlockThread)
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
	_, err := s.User.RequirePermission(ctx, entity.ActionEditTag)
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
	if err := s.TxAdapter.Begin(ctx); err != nil {
		return nil, err
	}
	user, err := s.User.RequirePermission(ctx, entity.ActionPubThread)
	if err != nil {
		return nil, s.TxAdapter.Rollback(ctx, err)
	}
	newThread, err := s.Forum.NewThread(ctx, user, thread)
	if err != nil {
		return nil, s.TxAdapter.Rollback(ctx, err)
	}
	return newThread, s.TxAdapter.Commit(ctx)
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
	if err := s.TxAdapter.Begin(ctx); err != nil {
		return nil, err
	}
	user, err := s.User.RequirePermission(ctx, entity.ActionPubPost)
	if err != nil {
		return nil, s.TxAdapter.Rollback(ctx, err)
	}
	res, err := s.Forum.NewPost(ctx, user, post)
	if err != nil {
		return nil, s.TxAdapter.Rollback(ctx, err)
	}
	// go func() {
	// 	ctx := context.Background()
	// 	ctx = s.TxAdapter.AttachDB(ctx)
	// 	if err := s.Noti.NewRepliedNoti(ctx, res.Thread, res.Post); err != nil {
	// 		log.Error(err)
	// 		return
	// 	}
	// 	quotePosts, err := res.Post.Quotes(ctx)
	// 	if err != nil {
	// 		log.Error(err)
	// 		return
	// 	}
	// 	for _, qp := range quotePosts {
	// 		if err := s.Noti.NewQuotedNoti(ctx, res.Thread, res.Post, qp); err != nil {
	// 			log.Error(err)
	// 		}
	// 	}
	// }()
	return res.Post, s.TxAdapter.Commit(ctx)
}

func (s *Service) GetPostByID(ctx context.Context, id uid.UID) (*entity.Post, error) {
	return s.Forum.GetPostByID(ctx, id)
}

func (s *Service) SetMainTags(ctx context.Context, tags []string) error {
	return s.Forum.SetMainTags(ctx, tags)
}

func (s *Service) GetMainTags(ctx context.Context) ([]string, error) {
	return s.Forum.GetMainTags(ctx)
}

func (s *Service) GetRecommendedTags(ctx context.Context) ([]string, error) {
	return s.Forum.GetMainTags(ctx)
}

func (s *Service) SearchTags(ctx context.Context, query *string, limit *int) ([]*entity.Tag, error) {
	return s.Forum.SearchTags(ctx, query, limit)
}

func (s *Service) GetUnreadNotiCount(ctx context.Context) (int, error) {
	user, err := s.User.RequirePermission(ctx, entity.ActionProfile)
	if err != nil {
		return 0, err
	}
	return s.Noti.GetUnreadNotiCount(ctx, user)
}

func (s *Service) GetNotification(ctx context.Context, query entity.SliceQuery) (*entity.NotiSlice, error) {
	user, err := s.User.RequirePermission(ctx, entity.ActionProfile)
	if err != nil {
		return nil, err
	}
	return s.Noti.GetNotification(ctx, user, query)
}

func (s *Service) GetNotificationHasRead(ctx context.Context, obj *entity.Notification) (bool, error) {
	user, err := s.User.RequirePermission(ctx, entity.ActionProfile)
	if err != nil {
		return false, err
	}
	return obj.HasRead(user), nil
}
