package uexky

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"gitlab.com/abyss.club/uexky/adapter"
	"gitlab.com/abyss.club/uexky/lib/algo"
	"gitlab.com/abyss.club/uexky/lib/config"
	"gitlab.com/abyss.club/uexky/lib/errors"
	"gitlab.com/abyss.club/uexky/lib/uid"
	"gitlab.com/abyss.club/uexky/uexky/entity"
)

type Service struct {
	TxAdapter adapter.Tx
	Repo      *entity.Repo
}

func NewService(tx adapter.Tx, repo *entity.Repo) (*Service, error) {
	s := &Service{TxAdapter: tx, Repo: repo}
	if err := loadMainTags(s); err != nil {
		return nil, err
	}
	return s, nil
}

func loadMainTags(s *Service) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
	defer cancel()
	ctx = s.TxAdapter.AttachDB(ctx)
	mainTags, err := s.Repo.Tag.GetMainTags(ctx)
	if err != nil {
		return errors.Wrap(err, "get main tags from db")
	}
	config.SetMainTags(mainTags)
	return nil
}

// ---- User Part ----

func (s *Service) AttachEmailUserToCtx(ctx context.Context, email string) (context.Context, error) {
	user, err := s.Repo.User.GetByEmail(ctx, email)
	if err != nil {
		if !errors.Is(err, errors.NotFound) {
			return nil, errors.Wrap(err, "User.GetByEmail")
		}

		// new signed in user
		user = entity.NewSignedInUser(email)
		user, err = s.Repo.User.Insert(ctx, user)
		if err != nil {
			return nil, errors.Wrap(err, "Create New User")
		}
		if err := s.NewNotiOnNewUser(ctx, user); err != nil {
			log.Error(err, "NewNotiOnNewUser")
		}
	}
	return user.AttachContext(ctx), nil
}

func (s *Service) AttachGuestUserToCtx(ctx context.Context, id uid.UID) (context.Context, error) {
	user, err := s.Repo.User.GetGuestByID(ctx, id)
	if err != nil {
		if !errors.Is(err, errors.NotFound) {
			return nil, errors.Wrap(err, "User.GetByID")
		}
		user = entity.NewGuestUser(id)
		user, err = s.Repo.User.Insert(ctx, user)
		if err != nil {
			return nil, errors.Wrap(err, "Create New Guest User")
		}
	}
	return user.AttachContext(ctx), nil
}

func (s *Service) Profile(ctx context.Context) (*entity.User, error) {
	if err := Cost(ctx, 1); err != nil {
		return nil, err
	}
	user := entity.GetCurrentUser(ctx)
	if err := user.RequirePermission(entity.ActionProfile); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) SetUserName(ctx context.Context, name string) (*entity.User, error) {
	if err := Cost(ctx, 1); err != nil {
		return nil, err
	}
	user := entity.GetCurrentUser(ctx)
	if err := user.RequirePermission(entity.ActionProfile); err != nil {
		return nil, err
	}
	if err := user.SetName(name); err != nil {
		return nil, err
	}
	return s.Repo.User.Update(ctx, user)
}

func (s *Service) GetUserThreads(
	ctx context.Context, obj *entity.User, query entity.SliceQuery,
) (*entity.ThreadSlice, error) {
	if err := Cost(ctx, query.Limit); err != nil {
		return nil, err
	}
	user := entity.GetCurrentUser(ctx)
	if err := user.RequirePermission(entity.ActionProfile); err != nil {
		return nil, err
	}
	if obj == nil || obj.Email != user.Email {
		return nil, errors.Permission.New("permission denied")
	}
	return s.Repo.User.ThreadSlice(ctx, user, query)
}

func (s *Service) GetUserPosts(ctx context.Context, obj *entity.User, query entity.SliceQuery) (*entity.PostSlice, error) {
	if err := Cost(ctx, query.Limit); err != nil {
		return nil, err
	}
	user := entity.GetCurrentUser(ctx)
	if err := user.RequirePermission(entity.ActionProfile); err != nil {
		return nil, err
	}
	if obj == nil || obj.Email != user.Email {
		return nil, errors.Permission.New("permission denied")
	}
	return s.Repo.User.PostSlice(ctx, user, query)
}

func (s *Service) SyncUserTags(ctx context.Context, tags []string) (*entity.User, error) {
	if err := Cost(ctx, 1); err != nil {
		return nil, err
	}
	user := entity.GetCurrentUser(ctx)
	if err := user.RequirePermission(entity.ActionProfile); err != nil {
		return nil, err
	}
	user.SetTags(tags)
	return s.Repo.User.Update(ctx, user)
}

func (s *Service) AddUserSubbedTag(ctx context.Context, tag string) (*entity.User, error) {
	if err := Cost(ctx, 1); err != nil {
		return nil, err
	}
	user := entity.GetCurrentUser(ctx)
	if err := user.RequirePermission(entity.ActionProfile); err != nil {
		return nil, err
	}
	user.AddTag(tag)
	return s.Repo.User.Update(ctx, user)
}

func (s *Service) DelUserSubbedTag(ctx context.Context, tag string) (*entity.User, error) {
	if err := Cost(ctx, 1); err != nil {
		return nil, err
	}
	user := entity.GetCurrentUser(ctx)
	if err := user.RequirePermission(entity.ActionProfile); err != nil {
		return nil, err
	}
	user.DelTag(tag)
	return s.Repo.User.Update(ctx, user)
}

func (s *Service) BanUser(ctx context.Context, postID *uid.UID, threadID *uid.UID) (bool, error) {
	user := entity.GetCurrentUser(ctx)
	if err := user.RequirePermission(entity.ActionBanUser); err != nil {
		return false, err
	}
	var targetID uid.UID
	switch {
	case postID != nil:
		post, err := s.Repo.Post.GetByID(ctx, *postID)
		if err != nil {
			return false, err
		}
		targetID = post.Author.UserID
	case threadID != nil:
		thread, err := s.Repo.Thread.GetByID(ctx, *threadID)
		if err != nil {
			return false, err
		}
		targetID = thread.Author.UserID
	default:
		return false, errors.BadParams.New("must specified post id or thread id")
	}
	target, err := s.Repo.User.GetByID(ctx, targetID)
	if err != nil {
		if errors.Is(err, errors.NotFound) {
			return false, nil
		}
		return false, errors.Wrap(err, "User.GetByID")
	}
	target.Ban()
	if _, err := s.Repo.User.Update(ctx, target); err != nil {
		return false, errors.Wrap(err, "User.Update")
	}
	return true, nil
}

// ---- Thread Part ----

func (s *Service) PubThread(ctx context.Context, thread entity.ThreadInput) (*entity.Thread, error) {
	if err := Cost(ctx, 1); err != nil {
		return nil, err
	}
	err := s.Repo.Thread.CheckIfDuplicated(ctx, thread.Title, thread.Content)
	if err != nil {
		return nil, errors.Wrap(err, "Thread.CheckIfDuplicated")
	}
	var newThread *entity.Thread
	err = s.TxAdapter.WithTx(ctx, func() error {
		user := entity.GetCurrentUser(ctx)
		if err := user.RequirePermission(entity.ActionPubThread); err != nil {
			return errors.Wrapf(err, "PubThread(thread=%+v)", thread)
		}
		t, err := entity.NewThread(user, thread)
		if err != nil {
			return err
		}
		t, err = s.Repo.Thread.Insert(ctx, t)
		if err != nil {
			return errors.Wrapf(err, "PubThread(thread=%+v)", thread)
		}
		newThread = t
		return nil
	})
	return newThread, err
}

func (s *Service) LockThread(ctx context.Context, threadID uid.UID) (*entity.Thread, error) {
	if err := Cost(ctx, 1); err != nil {
		return nil, err
	}
	user := entity.GetCurrentUser(ctx)
	if err := user.RequirePermission(entity.ActionLockThread); err != nil {
		return nil, err
	}
	thread, err := s.Repo.Thread.GetByID(ctx, threadID)
	if err != nil {
		return nil, err
	}
	thread.Lock()
	return s.Repo.Thread.Update(ctx, thread)
}

func (s *Service) BlockThread(ctx context.Context, threadID uid.UID) (*entity.Thread, error) {
	if err := Cost(ctx, 1); err != nil {
		return nil, err
	}
	user := entity.GetCurrentUser(ctx)
	if err := user.RequirePermission(entity.ActionBlockThread); err != nil {
		return nil, err
	}
	thread, err := s.Repo.Thread.GetByID(ctx, threadID)
	if err != nil {
		return nil, err
	}
	thread.Block()
	return s.Repo.Thread.Update(ctx, thread)
}

func (s *Service) EditTags(
	ctx context.Context, threadID uid.UID, mainTag string, subTags []string,
) (*entity.Thread, error) {
	if err := Cost(ctx, 1); err != nil {
		return nil, err
	}
	user := entity.GetCurrentUser(ctx)
	if err := user.RequirePermission(entity.ActionEditTag); err != nil {
		return nil, err
	}
	thread, err := s.Repo.Thread.GetByID(ctx, threadID)
	if err != nil {
		return nil, err
	}
	if err := thread.EditTags(mainTag, subTags); err != nil {
		return nil, err
	}
	return s.Repo.Thread.Update(ctx, thread)
}

func (s *Service) SearchThreads(
	ctx context.Context, tags []string, query entity.SliceQuery,
) (*entity.ThreadSlice, error) {
	if err := Cost(ctx, query.Limit); err != nil {
		return nil, err
	}
	return s.Repo.Thread.FindSlice(ctx, &entity.ThreadsSearch{Tags: tags}, query)
}

func (s *Service) GetThreadByID(ctx context.Context, id uid.UID) (*entity.Thread, error) {
	if err := Cost(ctx, 1); err != nil {
		return nil, err
	}
	return s.Repo.Thread.GetByID(ctx, id)
}

func (s *Service) GetThreadReplies(ctx context.Context, thread *entity.Thread, query entity.SliceQuery) (*entity.PostSlice, error) {
	if err := Cost(ctx, query.Limit); err != nil {
		return nil, err
	}
	postSlice, err := s.Repo.Thread.Replies(ctx, thread, query)
	if err != nil {
		return nil, errors.Wrap(err, "Thread.Replies")
	}
	return postSlice, err
}

func (s *Service) GetThreadReplyCount(ctx context.Context, thread *entity.Thread) (int, error) {
	count, err := s.Repo.Thread.ReplyCount(ctx, thread)
	return count, errors.Wrap(err, "Thread.ReplyCount")
}

func (s *Service) GetThreadCatalog(ctx context.Context, thread *entity.Thread) ([]*entity.ThreadCatalogItem, error) {
	catalogs, err := s.Repo.Thread.Catalog(ctx, thread)
	return catalogs, errors.Wrap(err, "Thread.Catalog")
}

// ---- Post Part ----

func (s *Service) PubPost(ctx context.Context, input entity.PostInput) (*entity.Post, error) {
	if err := Cost(ctx, 1); err != nil {
		return nil, err
	}
	user := entity.GetCurrentUser(ctx)
	if err := user.RequirePermission(entity.ActionPubPost); err != nil {
		return nil, errors.Wrapf(err, "PubPost(input=%+v)", input)
	}
	if err := s.Repo.Post.CheckIfDuplicated(ctx, user.ID, input.Content); err != nil {
		return nil, errors.Wrap(err, "Thread.CheckIfDuplicated")
	}
	var post *entity.Post
	var thread *entity.Thread
	err := s.TxAdapter.WithTx(ctx, func() error {
		var err error
		thread, err = s.Repo.Thread.GetByID(ctx, input.ThreadID)
		if err != nil {
			return errors.Wrap(err, "find thread")
		}
		var aid string
		if input.Anonymous {
			aid, err = s.Repo.Thread.PostAID(ctx, thread, user)
			if err != nil {
				return errors.Wrap(err, "Thread.PostAID")
			}
		}
		post, err = entity.NewPost(&input, user, thread, aid)
		if err != nil {
			return errors.Wrap(err, "NewPost")
		}
		post, err = s.Repo.Post.Insert(ctx, post)
		if err != nil {
			return errors.Wrapf(err, "PubPost(input=%+v)", input)
		}
		quotedPost, err := s.Repo.Post.QuotedPosts(ctx, post)
		if err != nil {
			return err
		}
		s.NewNotiOnNewPost(ctx, user, thread, post, quotedPost)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (s *Service) BlockPost(ctx context.Context, postID uid.UID) (*entity.Post, error) {
	if err := Cost(ctx, 1); err != nil {
		return nil, err
	}
	user := entity.GetCurrentUser(ctx)
	if err := user.RequirePermission(entity.ActionBlockPost); err != nil {
		return nil, err
	}
	post, err := s.Repo.Post.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}
	post.Block()
	return s.Repo.Post.Update(ctx, post)
}

func (s *Service) GetPostQuotedPosts(ctx context.Context, post *entity.Post) ([]*entity.Post, error) {
	quotes, err := s.Repo.Post.QuotedPosts(ctx, post)
	if err != nil {
		return nil, errors.Wrap(err, "Post.QuotedPosts")
	}
	return quotes, nil
}

func (s *Service) GetPostQuotedCount(ctx context.Context, post *entity.Post) (int, error) {
	count, err := s.Repo.Post.QuotedCount(ctx, post)
	if err != nil {
		return 0, errors.Wrap(err, "Post.QuotedCount")
	}
	return count, nil
}

func (s *Service) GetPostByID(ctx context.Context, id uid.UID) (*entity.Post, error) {
	if err := Cost(ctx, 1); err != nil {
		return nil, err
	}
	return s.Repo.Post.GetByID(ctx, id)
}

// ---- Tag Part ----

func (s *Service) SetMainTags(ctx context.Context, tags []string) error {
	if len(config.GetMainTags()) != 0 {
		return errors.BadParams.New("already have main tags")
	}
	if err := s.Repo.Tag.SetMainTags(ctx, tags); err != nil {
		return errors.Wrap(err, "set main tags to db")
	}
	config.SetMainTags(tags)
	return nil
}

func (s *Service) GetMainTags(ctx context.Context) []string {
	return config.GetMainTags()
}

func (s *Service) GetRecommendedTags(ctx context.Context) []string {
	return config.GetMainTags()
}

func (s *Service) SearchTags(ctx context.Context, query *string, limit *int) ([]*entity.Tag, error) {
	search := &entity.TagSearch{
		Text:  algo.NullToString(query),
		Limit: algo.NullToIntDefault(limit, 10),
	}
	if err := Cost(ctx, search.Limit); err != nil {
		return nil, err
	}
	return s.Repo.Tag.Search(ctx, search)
}

// ---- Notification Part ----

func (s *Service) GetUnreadNotiCount(ctx context.Context) (int, error) {
	user := entity.GetCurrentUser(ctx)
	if err := user.RequirePermission(entity.ActionProfile); err != nil {
		return 0, err
	}
	return s.Repo.Noti.GetUnreadCount(ctx, user)
}

func (s *Service) GetNotifications(ctx context.Context, query entity.SliceQuery) (*entity.NotiSlice, error) {
	if err := Cost(ctx, query.Limit); err != nil {
		return nil, err
	}
	user := entity.GetCurrentUser(ctx)
	if err := user.RequirePermission(entity.ActionProfile); err != nil {
		return nil, err
	}
	slice, err := s.Repo.Noti.GetSlice(ctx, user, query)
	if err != nil {
		return nil, errors.Wrapf(err, "GetNotification(user=%+v, query=%+v)", user, query)
	}
	if len(slice.Notifications) > 0 {
		lastRead := slice.Notifications[0].SortKey
		user.UpdateReadID(lastRead)
		if _, err := s.Repo.User.Update(ctx, user); err != nil {
			return nil, errors.Wrapf(err, "GetNotification(user=%+v, query=%+v)", user, query)
		}
	}
	return slice, nil
}

func (s *Service) NewNotiOnNewUser(ctx context.Context, user *entity.User) error {
	noti, err := entity.NewWelcomeNoti(user)
	if err != nil {
		return err
	}
	return s.Repo.Noti.Insert(ctx, noti)
}

func (s *Service) NewNotiOnNewPost(ctx context.Context, user *entity.User, thread *entity.Thread, post *entity.Post, quotedPosts []*entity.Post) {
	if user.ID != thread.Author.UserID {
		if err := s.newRepliedNoti(ctx, user, thread, post); err != nil {
			log.Error(err)
		}
	}
	for _, qp := range quotedPosts {
		if user.ID != qp.Author.UserID {
			if err := s.newQuotedNoti(ctx, thread, post, qp); err != nil {
				log.Error(err)
			}
		}
	}
}

func (s *Service) newRepliedNoti(ctx context.Context, user *entity.User, thread *entity.Thread, reply *entity.Post) error {
	key := entity.RepliedNotiKey(thread)
	prev, err := s.Repo.Noti.GetByKey(ctx, thread.Author.UserID, key)
	if err != nil {
		if errors.Is(err, errors.NotFound) {
			prev = nil
		} else {
			return errors.Wrap(err, "find prev replied noti")
		}
	}
	if prev == nil {
		noti := entity.NewRepliedNoti(user, thread, reply)
		if noti == nil {
			return nil
		}
		return s.Repo.Noti.Insert(ctx, noti)
	}
	prev.AddReply(user, thread, reply)
	return s.Repo.Noti.UpdateContent(ctx, prev)
}

func (s *Service) newQuotedNoti(ctx context.Context, thread *entity.Thread, post *entity.Post, quotedPost *entity.Post) error {
	noti := entity.NewQuotedNoti(thread, post, quotedPost)
	err := s.Repo.Noti.Insert(ctx, noti)
	return errors.Wrapf(err, "NewQuotedNoti(thread=%+v, post=%+v, quotedPost=%+v)", thread, post, quotedPost)
}
