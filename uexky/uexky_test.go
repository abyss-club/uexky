package uexky

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky/lib/algo"
	"gitlab.com/abyss.club/uexky/lib/config"
	"gitlab.com/abyss.club/uexky/lib/postgres"
	"gitlab.com/abyss.club/uexky/lib/redis"
	"gitlab.com/abyss.club/uexky/lib/uid"
	"gitlab.com/abyss.club/uexky/mocks"
	"gitlab.com/abyss.club/uexky/repo"
	"gitlab.com/abyss.club/uexky/uexky/entity"
)

func TestMain(m *testing.M) {
	if err := config.Load(""); err != nil {
		log.Fatalf("load config: %v", err)
	}
	fmt.Printf("run test in config: %#v\n", config.Get())
	os.Exit(m.Run())
}

// func getRedis(t *testing.T) *red.Client {
// 	rc, err := redis.NewClient()
// 	if err != nil {
// 		t.Fatalf("connect redis: %v", err)
// 	}
// 	return rc
// }

// copy from wire.InitDevService
// TODO: move service wire to here
func getService() (*Service, error) {
	client, err := redis.NewClient()
	if err != nil {
		return nil, err
	}
	forumRepo := &repo.ForumRepo{}
	userRepo := &repo.UserRepo{
		Redis: client,
		Forum: forumRepo,
	}
	mailAdapter := &mocks.MailAdapter{}
	userService := &entity.UserService{
		Repo: userRepo,
		Mail: mailAdapter,
	}
	forumService := &entity.ForumService{
		Repo: forumRepo,
	}
	notiRepo := &repo.NotiRepo{}
	notiService := &entity.NotiService{
		Repo: notiRepo,
	}
	db, err := postgres.NewDB()
	if err != nil {
		return nil, err
	}
	txAdapter := &postgres.TxAdapter{
		DB: db,
	}
	service := &Service{
		User:      userService,
		Forum:     forumService,
		Noti:      notiService,
		TxAdapter: txAdapter,
	}
	return service, nil
}

func getNewDBCtx(t *testing.T) context.Context {
	if err := postgres.RebuildDB(); err != nil {
		t.Fatal(err)
	}
	db, err := postgres.NewDB()
	if err != nil {
		t.Fatalf("get database: %v", err)
	}
	txAdapter := &postgres.TxAdapter{DB: db}
	ctx := context.Background()
	return txAdapter.AttachDB(ctx)
}

type testUser struct {
	email string
	name  string
}

func loginUser(service *Service, u testUser) (*entity.User, context.Context, error) {
	ctx := service.TxAdapter.AttachDB(context.Background())
	code, err := service.TrySignInByEmail(ctx, u.email)
	if err != nil {
		return nil, nil, errors.Wrap(err, "TrySignInByEmail")
	}
	token, err := service.SignInByCode(ctx, string(code))
	if err != nil {
		return nil, nil, errors.Wrap(err, "SignInByCode")
	}
	userCtx, err := service.CtxWithUserByToken(ctx, token.Tok)
	if err != nil {
		return nil, nil, errors.Wrap(err, "CtxWithUserByToken")
	}
	user, err := service.Profile(userCtx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Profile")
	}
	if u.name != "" && user.Name == nil {
		var err error
		if user, err = service.SetUserName(userCtx, u.name); err != nil {
			return nil, nil, errors.Wrap(err, "SetUserName")
		}
	}
	return user, userCtx, nil
}

var forumRepoComp = cmp.Comparer(func(lh, rh entity.ForumRepo) bool {
	return (lh == nil && rh == nil) || (lh != nil && rh != nil)
})

func pubThread(service *Service, u testUser) (*entity.Thread, context.Context, error) {
	user, ctx, err := loginUser(service, u)
	if err != nil {
		return nil, nil, err
	}
	mainTags, err := service.GetMainTags(ctx)
	if err != nil {
		return nil, nil, err
	}
	input := entity.ThreadInput{
		Anonymous: rand.Intn(2) == 0,
		Content:   uid.RandomBase64Str(50),
		MainTag:   mainTags[rand.Intn(len(mainTags))],
		SubTags:   []string{uid.RandomBase64Str(6), uid.RandomBase64Str(7)},
		Title:     algo.NullString(uid.RandomBase64Str(10)),
	}
	if user.Name == nil {
		input.Anonymous = true
	}
	thread, err := service.PubThread(ctx, input)
	return thread, ctx, err
}

func pubPost(service *Service, u testUser, threadID uid.UID, quotedIds []uid.UID) (*entity.Post, context.Context, error) {
	user, ctx, err := loginUser(service, u)
	if err != nil {
		return nil, nil, err
	}
	input := entity.PostInput{
		ThreadID:  threadID,
		Anonymous: rand.Intn(2) == 0,
		Content:   uid.RandomBase64Str(50),
		QuoteIds:  quotedIds,
	}
	if user.Name == nil {
		input.Anonymous = true
	}
	post, err := service.PubPost(ctx, input)
	return post, ctx, err
}
