package devtools

import (
	"context"
	"fmt"
	"math/rand"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.com/abyss.club/uexky/lib/algo"
	"gitlab.com/abyss.club/uexky/lib/config"
	"gitlab.com/abyss.club/uexky/lib/errors"
	"gitlab.com/abyss.club/uexky/lib/uid"
	"gitlab.com/abyss.club/uexky/uexky"
	"gitlab.com/abyss.club/uexky/uexky/entity"
)

type mockThreadsOpt struct {
	userCount    int
	threadCount  int
	maxPostCount int
	minPostCount int
}

var mockFlags mockThreadsOpt

func init() {
	mockDataCmd.PersistentFlags().IntVar(&mockFlags.userCount, "user", 5, "user count")
	mockDataCmd.PersistentFlags().IntVar(&mockFlags.threadCount, "thread", 100, "thread count")
	mockDataCmd.PersistentFlags().IntVar(&mockFlags.maxPostCount, "max", 100, "max posts count in thread")
	mockDataCmd.PersistentFlags().IntVar(&mockFlags.minPostCount, "min", 0, "min posts count in thread")
}

var mockDataCmd = &cobra.Command{
	Use:   "mock --user N[5] --thread N[100] --max N[100] --min[0]",
	Short: "mock test data",
	Run: func(cmd *cobra.Command, args []string) {
		if err := mockThreads(&mockFlags); err != nil {
			log.Fatalf("%+v", err)
		}
	},
}

func createUser(service *uexky.Service) (*entity.User, error) {
	ctx := service.TxAdapter.AttachDB(context.Background())
	email := fmt.Sprintf("%s@%s", uid.RandomBase64Str(16), config.Get().Server.Domain)
	ctx, err := service.AttachEmailUserToCtx(ctx, email)
	if err != nil {
		return nil, errors.Wrap(err, "AttachEmailUserToCtx")
	}
	return entity.GetCurrentUser(ctx), nil
}

func mockThreads(opt *mockThreadsOpt) error {
	service, err := uexky.InitUexkyService()
	ctx := service.TxAdapter.AttachDB(context.Background())
	if err != nil {
		return errors.Wrap(err, "init service")
	}
	mainTags := service.GetMainTags(ctx)
	if len(mainTags) == 0 {
		return errors.Internal.New("no main tags, set main tags first")
	}
	var users []*entity.User
	for i := 0; i < opt.userCount; i++ {
		user, err := createUser(service)
		if err != nil {
			return err
		}
		users = append(users, user)
	}
	var subTags []string
	for i := 0; i < opt.threadCount; i++ {
		t := fmt.Sprintf("st:%s", uid.RandomBase64Str(5))
		subTags = append(subTags, t)
	}
	var wg sync.WaitGroup
	var errs []error
	wg.Add(opt.threadCount)
	workers := make(chan struct{}, 16)
	for i := 0; i < opt.threadCount; i++ {
		workers <- struct{}{}
		go func(i int) {
			defer func() {
				<-workers
				wg.Done()
			}()
			if len(errs) != 0 {
				return
			}
			pc := opt.minPostCount + rand.Intn(1+opt.maxPostCount-opt.minPostCount)
			if err := makeThread(service, users, mainTags, subTags, pc, i); err != nil {
				log.Error(err)
				errs = append(errs, errors.Wrapf(err, "make thread %v", i+1))
				return
			}
			fmt.Println("create thread: ", i+1)
		}(i)
	}
	wg.Wait()
	if len(errs) != 0 {
		return errors.Internal.New("create threads error")
	}
	return nil
}

func makeThread(
	service *uexky.Service, users []*entity.User, mainTags []string, subTags []string, postCount int, index int,
) error {
	input := &entity.ThreadInput{
		Anonymous: rand.Intn(2) == 1,
		Content:   uid.RandomBase64Str(200),
		MainTag:   mainTags[rand.Intn(len(mainTags))],
	}
	subTagsCount := rand.Intn(4)
	for i := 0; i < subTagsCount; i++ {
		t := subTags[rand.Intn(len(subTags))]
		input.SubTags = append(input.SubTags, t)
	}
	if rand.Intn(2) == 1 {
		input.Title = algo.NullString(uid.RandomBase64Str(20))
	}
	user := users[rand.Intn(len(users))]
	var err error
	ctx := service.TxAdapter.AttachDB(context.Background())
	ctx = user.AttachContext(ctx)
	thread, err := service.PubThread(ctx, *input)
	if err != nil {
		return errors.Wrap(err, "create thread")
	}
	var posts []*entity.Post
	for i := 0; i < postCount; i++ {
		qCount := quotedCount()
		var qids []uid.UID
		for i := 0; i < len(posts) && i < qCount; i++ {
			qids = append(qids, posts[rand.Intn(len(posts))].ID)
		}
		post, err := makePost(service, users, thread, qids)
		if err != nil {
			return errors.Wrapf(err, "create thread %v post %v", index+1, i+1)
		}
		log.Infof("create thread %v post %v", index+1, i+1)
		posts = append(posts, post)
	}
	return nil
}

func quotedCount() int {
	w := rand.Intn(10)
	switch {
	case w < 5: // 50%
		return 0
	case w < 7: // 20%
		return 1
	case w < 9:
		return 2 // 10%
	default:
		return 3 // 10%
	}
}

func makePost(
	service *uexky.Service, users []*entity.User, thread *entity.Thread, quotedIds []uid.UID,
) (*entity.Post, error) {
	input := &entity.PostInput{
		ThreadID:  thread.ID,
		Anonymous: rand.Intn(2) == 1,
		Content:   uid.RandomBase64Str(200),
		QuoteIds:  quotedIds,
	}
	user := users[rand.Intn(len(users))]
	var err error
	ctx := service.TxAdapter.AttachDB(context.Background())
	ctx = user.AttachContext(ctx)
	post, err := service.PubPost(ctx, *input)
	if err != nil {
		return nil, errors.Wrap(err, "create post")
	}
	return post, nil
}
