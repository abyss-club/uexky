package uexky

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky/lib/algo"
	"gitlab.com/abyss.club/uexky/lib/config"
	"gitlab.com/abyss.club/uexky/lib/postgres"
	"gitlab.com/abyss.club/uexky/lib/uid"
	"gitlab.com/abyss.club/uexky/uexky/entity"
)

func TestMain(m *testing.M) {
	if err := config.Load(""); err != nil {
		log.Fatalf("load config: %v", err)
	}
	fmt.Printf("run test in config: %#v\n", config.Get())
	os.Exit(m.Run())
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

func loginUser(t *testing.T, service *Service, u testUser) (*entity.User, context.Context) {
	if u.email != "" {
		return signedInUser(t, service, u)
	}
	return guestUser(t, service)
}

func signedInUser(t *testing.T, service *Service, u testUser) (*entity.User, context.Context) {
	ctx := service.TxAdapter.AttachDB(context.Background())
	code, err := service.TrySignInByEmail(ctx, u.email)
	if err != nil {
		t.Fatal(errors.Wrap(err, "TrySignInByEmail"))
	}
	token, err := service.SignInByCode(ctx, string(code))
	if err != nil {
		t.Fatal(errors.Wrap(err, "SignInByCode"))
	}
	userCtx, _, err := service.CtxWithUserByToken(ctx, token.Tok)
	if err != nil {
		t.Fatal(errors.Wrap(err, "CtxWithUserByToken"))
	}
	user, err := service.Profile(userCtx)
	if err != nil {
		t.Fatal(errors.Wrap(err, "Profile"))
	}
	if u.name != "" && user.Name == nil {
		var err error
		if user, err = service.SetUserName(userCtx, u.name); err != nil {
			t.Fatal(errors.Wrap(err, "SetUserName"))
		}
	}
	return user, userCtx
}

func guestUser(t *testing.T, service *Service) (*entity.User, context.Context) {
	ctx := service.TxAdapter.AttachDB(context.Background())
	ctx, _, err := service.CtxWithUserByToken(ctx, "")
	if err != nil {
		t.Fatal(errors.Wrap(err, "CtxWithUserByToken"))
	}
	user, err := service.Profile(ctx)
	if err != nil {
		t.Fatal(errors.Wrap(err, "Profile"))
	}
	return user, ctx
}

func pubThread(t *testing.T, service *Service, u testUser) (*entity.Thread, context.Context) {
	var user *entity.User
	var ctx context.Context
	user, ctx = loginUser(t, service, u)
	mainTags, err := service.GetMainTags(ctx)
	if err != nil {
		t.Fatal(err)
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
	if err != nil {
		t.Fatal(err)
	}
	return thread, ctx
}

func pubThreadWithTags(t *testing.T, service *Service, u testUser, mainTag string, subTags []string) (*entity.Thread, context.Context) {
	user, ctx := loginUser(t, service, u)
	_, err := service.GetMainTags(ctx)
	if err != nil {
		t.Fatal(err)
	}
	input := entity.ThreadInput{
		Anonymous: rand.Intn(2) == 0,
		Content:   uid.RandomBase64Str(50),
		MainTag:   mainTag,
		SubTags:   subTags,
		Title:     algo.NullString(uid.RandomBase64Str(10)),
	}
	if user.Name == nil {
		input.Anonymous = true
	}
	thread, err := service.PubThread(ctx, input)
	if err != nil {
		t.Fatal(err)
	}
	return thread, ctx
}

func pubPost(t *testing.T, service *Service, u testUser, threadID uid.UID, quotedIds ...uid.UID) (*entity.Post, context.Context) {
	var user *entity.User
	var ctx context.Context
	user, ctx = loginUser(t, service, u)
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
	if err != nil {
		t.Fatal(err)
	}
	return post, ctx
}

var forumRepoComp = cmp.Comparer(func(lh, rh entity.ForumRepo) bool {
	return (lh == nil && rh == nil) || (lh != nil && rh != nil)
})

var timeCmp = cmp.Comparer(func(lh, rh time.Time) bool {
	delta := lh.Sub(rh)
	return delta > (-10*time.Millisecond) && delta < (10*time.Millisecond)
})

var tagSetCmp = cmp.Comparer(func(lh, rh []*entity.Tag) bool {
	sort.SliceStable(lh, func(i, j int) bool {
		return lh[i].Name < lh[j].Name
	})
	sort.SliceStable(rh, func(i, j int) bool {
		return rh[i].Name < rh[j].Name
	})
	return cmp.Equal(lh, rh)
})