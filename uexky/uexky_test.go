package uexky

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gitlab.com/abyss.club/uexky/lib/algo"
	"gitlab.com/abyss.club/uexky/lib/config"
	"gitlab.com/abyss.club/uexky/lib/errors"
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

var dbLock sync.Mutex

func initEnv(t *testing.T, mainTags ...string) (*Service, context.Context) {
	dbLock.Lock()
	defer dbLock.Unlock()
	if err := postgres.RebuildDB(); err != nil {
		t.Fatal(err)
	}
	db, err := postgres.NewDB()
	if err != nil {
		t.Fatalf("get database: %v", err)
	}
	txAdapter := &postgres.TxAdapter{DB: db}
	ctx := txAdapter.AttachDB(context.Background())
	service, err := InitUexkyService()
	if err != nil {
		t.Fatal(err)
	}
	if len(mainTags) != 0 {
		t.Logf("before set mainTags, service.Maintags = %v", service.GetMainTags(context.Background()))
		if err := service.SetMainTags(ctx, mainTags); err != nil {
			t.Fatal(errors.Wrapf(err, "SetMainTags(tags=%v)", mainTags))
		}
		t.Logf("after set mainTags, service.Maintags = %v", service.GetMainTags(context.Background()))
	}
	return service, txAdapter.AttachDB(ctx)
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
	ctx, err := service.AttachEmailUserToCtx(ctx, u.email)
	if err != nil {
		t.Fatal(errors.Wrap(err, "AttachEmailUserToCtx"))
	}
	user, err := service.Profile(ctx)
	if err != nil {
		t.Fatal(errors.Wrap(err, "Profile"))
	}
	if u.name != "" && user.Name == nil {
		user, err = service.SetUserName(ctx, u.name)
		if err != nil {
			t.Fatal(errors.Wrap(err, "SetUserName"))
		}
	}
	return user, ctx
}

func guestUser(t *testing.T, service *Service) (*entity.User, context.Context) {
	ctx := service.TxAdapter.AttachDB(context.Background())
	ctx, err := service.AttachGuestUserToCtx(ctx, uid.NewUID())
	if err != nil {
		t.Fatal(errors.Wrap(err, "AttachGuestUserToCtx"))
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
	mainTags := service.GetMainTags(ctx)
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

var tagSetCmp = cmp.Comparer(func(lh, rh []*entity.Tag) bool {
	sort.SliceStable(lh, func(i, j int) bool {
		return lh[i].Name < lh[j].Name
	})
	sort.SliceStable(rh, func(i, j int) bool {
		return rh[i].Name < rh[j].Name
	})
	return cmp.Equal(lh, rh)
})
