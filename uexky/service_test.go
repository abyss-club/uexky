package uexky

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"gitlab.com/abyss.club/uexky/lib/algo"
	"gitlab.com/abyss.club/uexky/lib/errors"
	"gitlab.com/abyss.club/uexky/lib/uid"
	"gitlab.com/abyss.club/uexky/uexky/entity"
)

func TestService_SetUserName(t *testing.T) {
	service, ctx := initEnv(t)
	type args struct {
		ctx   context.Context
		name  string
		email string
	}
	tests := []struct {
		name     string
		args     args
		wantName string
		wantErr  bool
	}{
		{
			name: "normal",
			args: args{
				ctx:   ctx,
				name:  "tom",
				email: "tom@example.com",
			},
			wantName: "tom",
			wantErr:  false,
		},
		{
			name: "already has name",
			args: args{
				ctx:   ctx,
				name:  "tom2",
				email: "tom@example.com",
			},
			wantName: "tom",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ctx := loginUser(t, service, testUser{email: tt.args.email})
			gotUser, err := service.SetUserName(ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.SetUserName() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if gotUser == nil || gotUser.Name == nil || *gotUser.Name != tt.wantName {
				t.Errorf("Service.SetUserName() = %+v, want user.name=%v", gotUser, tt.wantName)
			}
		})
	}
}

func TestService_GetUserThreads(t *testing.T) {
	mainTags := []string{"MainA", "MainB", "MainC"}
	service, _ := initEnv(t, mainTags...)

	var threads []*entity.Thread
	tu := testUser{email: "a@example.com", name: "a"}
	for i := 0; i < 10; i++ {
		thread, _ := pubThread(t, service, tu)
		threads = append(threads, thread)
	}
	user, ctx := loginUser(t, service, testUser{email: "a@example.com"})
	type args struct {
		ctx   context.Context
		obj   *entity.User
		query entity.SliceQuery
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.ThreadSlice
		wantErr bool
	}{
		{
			name: "first 5 threads",
			args: args{
				ctx: ctx,
				obj: user,
				query: entity.SliceQuery{
					After: algo.NullString(""),
					Limit: 5,
				},
			},
			want: &entity.ThreadSlice{
				Threads: []*entity.Thread{
					threads[9], threads[8], threads[7], threads[6], threads[5],
				},
				SliceInfo: &entity.SliceInfo{
					FirstCursor: threads[9].ID.ToBase64String(),
					LastCursor:  threads[5].ID.ToBase64String(),
					HasNext:     true,
				},
			},
		},
		{
			name: "next 5 threads",
			args: args{
				ctx: ctx,
				obj: user,
				query: entity.SliceQuery{
					After: algo.NullString(threads[5].ID.ToBase64String()),
					Limit: 5,
				},
			},
			want: &entity.ThreadSlice{
				Threads: []*entity.Thread{
					threads[4], threads[3], threads[2], threads[1], threads[0],
				},
				SliceInfo: &entity.SliceInfo{
					FirstCursor: threads[4].ID.ToBase64String(),
					LastCursor:  threads[0].ID.ToBase64String(),
					HasNext:     false,
				},
			},
		},
		{
			name: "last 5 threads",
			args: args{
				ctx: ctx,
				obj: user,
				query: entity.SliceQuery{
					Before: algo.NullString(""),
					Limit:  5,
				},
			},
			want: &entity.ThreadSlice{
				Threads: []*entity.Thread{
					threads[4], threads[3], threads[2], threads[1], threads[0],
				},
				SliceInfo: &entity.SliceInfo{
					FirstCursor: threads[4].ID.ToBase64String(),
					LastCursor:  threads[0].ID.ToBase64String(),
					HasNext:     true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetUserThreads(tt.args.ctx, tt.args.obj, tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetUserThreads() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got.SliceInfo, tt.want.SliceInfo); diff != "" {
				t.Errorf("Service.GetUserThreads().SliceInfo diff: %s", diff)
			}
			if len(got.Threads) != len(tt.want.Threads) {
				t.Errorf("Service.GetUserThreads().len(Threads) = %v, want %v", len(got.Threads), len(tt.want.Threads))
			}
			for i, thread := range got.Threads {
				thread.CreatedAt = tt.want.Threads[i].CreatedAt
				if diff := cmp.Diff(thread, tt.want.Threads[i]); diff != "" {
					t.Errorf("Service.GetUserThreads().Threads[%v] diff: %s", i, diff)
				}
			}
		})
	}
}

func TestService_GetUserPosts(t *testing.T) {
	mainTags := []string{"MainA", "MainB", "MainC"}
	service, _ := initEnv(t, mainTags...)

	var threads []*entity.Thread
	for i := 0; i < 2; i++ {
		thread, _ := pubThread(t, service, testUser{email: "thread@example.com"})
		threads = append(threads, thread)
	}
	tu := testUser{email: "post@example.com", name: " post"}
	var posts []*entity.Post
	for i := 0; i < 10; i++ {
		post, _ := pubPost(t, service, tu, threads[i%2].ID)
		posts = append(posts, post)
	}
	user, ctx := loginUser(t, service, tu)
	type args struct {
		ctx   context.Context
		obj   *entity.User
		query entity.SliceQuery
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.PostSlice
		wantErr bool
	}{
		{
			name: "first 5 posts",
			args: args{
				ctx: ctx,
				obj: user,
				query: entity.SliceQuery{
					After: algo.NullString(""),
					Limit: 5,
				},
			},
			want: &entity.PostSlice{
				Posts: []*entity.Post{
					posts[9], posts[8], posts[7], posts[6], posts[5],
				},
				SliceInfo: &entity.SliceInfo{
					FirstCursor: posts[9].ID.ToBase64String(),
					LastCursor:  posts[5].ID.ToBase64String(),
					HasNext:     true,
				},
			},
		},
		{
			name: "next 5 posts",
			args: args{
				ctx: ctx,
				obj: user,
				query: entity.SliceQuery{
					After: algo.NullString(posts[5].ID.ToBase64String()),
					Limit: 5,
				},
			},
			want: &entity.PostSlice{
				Posts: []*entity.Post{
					posts[4], posts[3], posts[2], posts[1], posts[0],
				},
				SliceInfo: &entity.SliceInfo{
					FirstCursor: posts[4].ID.ToBase64String(),
					LastCursor:  posts[0].ID.ToBase64String(),
					HasNext:     false,
				},
			},
		},
		{
			name: "last 5 posts",
			args: args{
				ctx: ctx,
				obj: user,
				query: entity.SliceQuery{
					Before: algo.NullString(""),
					Limit:  5,
				},
			},
			want: &entity.PostSlice{
				Posts: []*entity.Post{
					posts[4], posts[3], posts[2], posts[1], posts[0],
				},
				SliceInfo: &entity.SliceInfo{
					FirstCursor: posts[4].ID.ToBase64String(),
					LastCursor:  posts[0].ID.ToBase64String(),
					HasNext:     true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetUserPosts(tt.args.ctx, tt.args.obj, tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetUserPosts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got.SliceInfo, tt.want.SliceInfo); diff != "" {
				t.Errorf("Service.GetUserPosts().SliceInfo diff: %s", diff)
			}
			if len(got.Posts) != len(tt.want.Posts) {
				t.Errorf("Service.GetUserPosts().len(Posts) = %v, want %v", len(got.Posts), len(tt.want.Posts))
			}
			for i, post := range got.Posts {
				post.CreatedAt = tt.want.Posts[i].CreatedAt
				if diff := cmp.Diff(post, tt.want.Posts[i]); diff != "" {
					t.Errorf("Service.GetUserPosts().Posts[%v] diff: %s", i, diff)
				}
			}
		})
	}
}

func TestService_SyncUserTags(t *testing.T) {
	mainTags := []string{"MainA", "MainB", "MainC"}
	service, _ := initEnv(t, mainTags...)

	_, ctx := loginUser(t, service, testUser{name: "a", email: "a@example.com"})
	type args struct {
		ctx  context.Context
		tags []string
	}
	tests := []struct {
		name     string
		args     args
		wantTags []string
		wantErr  bool
	}{
		{
			name: "sync 3 tags",
			args: args{
				ctx:  ctx,
				tags: []string{"MainA", "SubA", "SubB"},
			},
			wantTags: []string{"MainA", "SubA", "SubB"},
			wantErr:  false,
		},
		{
			name: "sync tags to add and del",
			args: args{
				ctx:  ctx,
				tags: []string{"SubA", "SubB", "SubC"},
			},
			wantTags: []string{"SubA", "SubB", "SubC"},
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.SyncUserTags(tt.args.ctx, tt.args.tags)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.SyncUserTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got.Tags, tt.wantTags); diff != "" {
				t.Errorf("Service.SyncUserTags() = %v, wantTags %v", got.Tags, tt.wantTags)
			}
		})
	}
}

func TestService_AddUserSubbedTag(t *testing.T) {
	mainTags := []string{"MainA", "MainB", "MainC"}
	service, _ := initEnv(t, mainTags...)

	_, ctx := loginUser(t, service, testUser{name: "a", email: "a@example.com"})
	type args struct {
		ctx context.Context
		tag string
	}
	tests := []struct {
		name     string
		args     args
		wantTags []string
		wantErr  bool
	}{
		{
			name:     "add 1 tag",
			args:     args{ctx: ctx, tag: "subA"},
			wantTags: []string{"MainA", "MainB", "MainC", "subA"},
			wantErr:  false,
		},
		{
			name:     "add 1 duplicated tag",
			args:     args{ctx: ctx, tag: "MainB"},
			wantTags: []string{"MainA", "MainB", "MainC", "subA"},
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.AddUserSubbedTag(tt.args.ctx, tt.args.tag)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.AddUserSubbedTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got.Tags, tt.wantTags); diff != "" {
				t.Errorf("Service.AddUserSubbedTag() = %v, wantTags %v", got.Tags, tt.wantTags)
			}
		})
	}
}

func TestService_DelUserSubbedTag(t *testing.T) {
	mainTags := []string{"MainA", "MainB", "MainC"}
	service, _ := initEnv(t, mainTags...)

	_, ctx := loginUser(t, service, testUser{name: "a", email: "a@example.com"})
	type args struct {
		ctx context.Context
		tag string
	}
	tests := []struct {
		name     string
		args     args
		wantTags []string
		wantErr  bool
	}{
		{
			name: "del 1 tag",
			args: args{
				ctx: ctx,
				tag: "MainA",
			},
			wantTags: []string{"MainB", "MainC"},
			wantErr:  false,
		},
		{
			name: "del 1 unexists tag",
			args: args{
				ctx: ctx,
				tag: "subC",
			},
			wantTags: []string{"MainB", "MainC"},
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.DelUserSubbedTag(tt.args.ctx, tt.args.tag)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.DelUserSubbedTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got.Tags, tt.wantTags); diff != "" {
				t.Errorf("Service.DelUserSubbedTag() diff = %v", diff)
			}
		})
	}
}

func TestService_BanUser(t *testing.T) {
	mainTags := []string{"MainA", "MainB", "MainC"}
	service, ctx := initEnv(t, mainTags...)

	thread, _ := pubThread(t, service, testUser{email: "t@example.com"})
	post, _ := pubPost(t, service, testUser{email: "p@example.com"}, thread.ID)
	mod, _ := loginUser(t, service, testUser{email: "mod@example.com"})
	mod.Role = entity.RoleMod
	_, err := service.Repo.User.Update(ctx, mod)
	if err != nil {
		t.Fatal(err)
	}
	_, modCtx := loginUser(t, service, testUser{email: "mod@example.com"})
	type args struct {
		ctx      context.Context
		postID   *uid.UID
		threadID *uid.UID
	}
	tests := []struct {
		name       string
		args       args
		checkEmail string
		wantBanned bool
		wantErr    bool
	}{
		{
			name: "ban user by post id",
			args: args{
				ctx:    modCtx,
				postID: &post.ID,
			},
			checkEmail: "p@example.com",
			wantBanned: true,
		},
		{
			name: "ban user by thread id",
			args: args{
				ctx:      modCtx,
				threadID: &thread.ID,
			},
			checkEmail: "t@example.com",
			wantBanned: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.BanUser(tt.args.ctx, tt.args.postID, tt.args.threadID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.BanUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			user, _ := loginUser(t, service, testUser{email: tt.checkEmail})
			if err != nil {
				t.Error(errors.Wrap(err, "get user profile"))
			}
			if (user.Role == entity.RoleBanned) != tt.wantBanned {
				t.Errorf("user role = %v, wantBanned %v", user.Role, tt.wantBanned)
			}
		})
	}
}

func TestService_PubThread(t *testing.T) {
	mainTags := []string{"MainA", "MainB", "MainC"}
	service, _ := initEnv(t, mainTags...)

	user, ctx := loginUser(t, service, testUser{email: "a@example.com", name: "a"})
	gUser, gCtx := loginUser(t, service, testUser{})
	type args struct {
		ctx    context.Context
		thread entity.ThreadInput
	}
	tests := []struct {
		name      string
		args      args
		want      *entity.Thread
		wantErrIs error
	}{
		{
			name: "anonymous signed in user",
			args: args{
				ctx,
				entity.ThreadInput{
					Anonymous: true,
					Content:   "content",
					MainTag:   "MainA",
					SubTags:   []string{"SubA", "SubB"},
					Title:     algo.NullString("title"),
				},
			},
			want: &entity.Thread{
				Author: &entity.Author{
					Anonymous: true,
					UserID:    user.ID,
					Author:    user.ID.ToBase64String(),
				},
				Title:   algo.NullString("title"),
				Content: "content",
				MainTag: "MainA",
				SubTags: []string{"SubA", "SubB"},
			},
		},
		{
			name: "pub thread with user name",
			args: args{
				ctx,
				entity.ThreadInput{
					Anonymous: false,
					Content:   "content1",
					MainTag:   "MainA",
					SubTags:   []string{"SubA", "SubB", "SubC"},
				},
			},
			want: &entity.Thread{
				Author: &entity.Author{
					UserID:    user.ID,
					Anonymous: false,
					Author:    *user.Name,
				},
				Title:   nil,
				Content: "content1",
				MainTag: "MainA",
				SubTags: []string{"SubA", "SubB", "SubC"},
			},
		},
		{
			name: "pub duplicated thread",
			args: args{
				ctx,
				entity.ThreadInput{
					Anonymous: false,
					Content:   "content1",
					MainTag:   "MainA",
					SubTags:   []string{"SubA", "SubB", "SubC"},
				},
			},
			wantErrIs: errors.Duplicated.New(),
		},
		{
			name: "guest user",
			args: args{
				gCtx,
				entity.ThreadInput{
					Anonymous: true,
					Content:   "content",
					MainTag:   "MainA",
					SubTags:   []string{"SubA", "SubB", "SubC"},
				},
			},
			want: &entity.Thread{
				Author: &entity.Author{
					UserID:    gUser.ID,
					Guest:     true,
					Anonymous: true,
					Author:    gUser.ID.ToBase64String(),
				},
				Title:   nil,
				Content: "content",
				MainTag: "MainA",
				SubTags: []string{"SubA", "SubB", "SubC"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.PubThread(tt.args.ctx, tt.args.thread)
			if (tt.wantErrIs == nil && err != nil) || (tt.wantErrIs != nil && !errors.Is(err, tt.wantErrIs)) {
				t.Errorf("Service.PubThread() error = %v, wantErr %v", err, tt.wantErrIs)
			}
			if err != nil {
				return
			}
			tt.want.ID = got.ID
			tt.want.CreatedAt = got.CreatedAt
			tt.want.LastPostID = got.LastPostID
			if tt.want.Author.Anonymous {
				tt.want.Author.Author = got.Author.Author
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("Service.PubThread() missmatch: %s", diff)
			}
			if got.LastPostID != got.ID {
				t.Errorf("New thread's last_psot_id should equal to id")
			}
		})
	}
}

func TestService_LockThread(t *testing.T) {
	mainTags := []string{"MainA", "MainB", "MainC"}
	service, ctx := initEnv(t, mainTags...)

	oriThread, _ := pubThread(t, service, testUser{email: "t@example.com"})

	mod, _ := loginUser(t, service, testUser{email: "mod@example.com"})
	mod.Role = entity.RoleMod
	_, err := service.Repo.User.Update(ctx, mod)
	if err != nil {
		t.Fatal(err)
	}
	oriThread.Locked = true
	_, modCtx := loginUser(t, service, testUser{email: "mod@example.com"})

	t.Run("check thread in memory", func(t *testing.T) {
		thread, err := service.LockThread(modCtx, oriThread.ID)
		if err != nil {
			t.Fatal(errors.Wrap(err, "LockThread"))
		}
		oriThread.CreatedAt = thread.CreatedAt
		if diff := cmp.Diff(thread, oriThread); diff != "" {
			t.Errorf("LockThread() not matched: %s", diff)
		}
	})
	t.Run("check thread in database", func(t *testing.T) {
		thread, err := service.GetThreadByID(modCtx, oriThread.ID)
		if err != nil {
			t.Fatal(errors.Wrap(err, "GetThreadByID"))
		}
		oriThread.CreatedAt = thread.CreatedAt
		if diff := cmp.Diff(thread, oriThread); diff != "" {
			t.Errorf("LockThread() not matched: %s", diff)
		}
	})
	t.Run("thread locked", func(t *testing.T) {
		_, ctx := loginUser(t, service, testUser{email: "p@example.com"})
		input := entity.PostInput{
			ThreadID:  oriThread.ID,
			Anonymous: true,
			Content:   uid.RandomBase64Str(50),
		}
		_, err := service.PubPost(ctx, input)
		if err == nil {
			t.Errorf("locked thread still can send new post")
		}
	})
}

func TestService_BlockThread(t *testing.T) {
	mainTags := []string{"MainA", "MainB", "MainC"}
	service, ctx := initEnv(t, mainTags...)

	oriThread, _ := pubThread(t, service, testUser{email: "t@example.com"})

	mod, _ := loginUser(t, service, testUser{email: "mod@example.com"})
	mod.Role = entity.RoleMod
	_, err := service.Repo.User.Update(ctx, mod)
	if err != nil {
		t.Fatal(err)
	}
	oriThread.Blocked = true
	oriThread.Content = entity.BlockedContent
	_, modCtx := loginUser(t, service, testUser{email: "mod@example.com"})

	t.Run("check thread in memory", func(t *testing.T) {
		thread, err := service.BlockThread(modCtx, oriThread.ID)
		if err != nil {
			t.Fatal(errors.Wrap(err, "BlockThread"))
		}
		oriThread.CreatedAt = thread.CreatedAt
		if diff := cmp.Diff(thread, oriThread); diff != "" {
			t.Errorf("LockThread() not matched: %s", diff)
		}
	})
	t.Run("check thread in database", func(t *testing.T) {
		thread, err := service.GetThreadByID(modCtx, oriThread.ID)
		if err != nil {
			t.Fatal(errors.Wrap(err, "GetThreadByID"))
		}
		oriThread.CreatedAt = thread.CreatedAt
		if diff := cmp.Diff(thread, oriThread); diff != "" {
			t.Errorf("LockThread() not matched: %s", diff)
		}
	})
}

func TestService_EditTags(t *testing.T) {
	mainTags := []string{"MainA", "MainB", "MainC"}
	service, ctx := initEnv(t, mainTags...)

	oriThread, _ := pubThread(t, service, testUser{email: "t@example.com"})

	mod, _ := loginUser(t, service, testUser{email: "mod@example.com"})
	mod.Role = entity.RoleMod
	_, err := service.Repo.User.Update(ctx, mod)
	if err != nil {
		t.Fatal(err)
	}
	oriThread.MainTag = "MainC"
	oriThread.SubTags = []string{"SubC", "SubB", "SubA"}
	_, modCtx := loginUser(t, service, testUser{email: "mod@example.com"})

	t.Run("check thread in memory", func(t *testing.T) {
		thread, err := service.EditTags(modCtx, oriThread.ID, oriThread.MainTag, oriThread.SubTags)
		if err != nil {
			t.Fatal(errors.Wrap(err, "EditTags"))
		}
		oriThread.CreatedAt = thread.CreatedAt
		if diff := cmp.Diff(thread, oriThread); diff != "" {
			t.Errorf("EditTags() not matched: %s", diff)
		}
	})
	t.Run("check thread in database", func(t *testing.T) {
		thread, err := service.GetThreadByID(modCtx, oriThread.ID)
		if err != nil {
			t.Fatal(errors.Wrap(err, "GetThreadByID"))
		}
		oriThread.CreatedAt = thread.CreatedAt
		if diff := cmp.Diff(thread, oriThread); diff != "" {
			t.Errorf("EditTags() not matched: %s", diff)
		}
	})
}

func TestService_SearchThreads(t *testing.T) {
	mainTags := []string{"MainA", "MainB", "MainC"}
	service, _ := initEnv(t, mainTags...)

	var threads []*entity.Thread
	tu := testUser{email: "a@example.com", name: "a"}
	for i := 0; i < 10; i++ {
		var thread *entity.Thread
		switch {
		case i < 3:
			thread, _ = pubThreadWithTags(t, service, tu, "MainA", []string{"SubA"})
		case i < 6:
			thread, _ = pubThreadWithTags(t, service, tu, "MainA", []string{"SubA", "SubB"})
		default:
			thread, _ = pubThreadWithTags(t, service, tu, "MainC", []string{"SubB", "SubC"})
		}
		threads = append(threads, thread)
	}
	_, ctx := loginUser(t, service, testUser{email: "a@example.com"})
	type args struct {
		ctx   context.Context
		tags  []string
		query entity.SliceQuery
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.ThreadSlice
		wantErr bool
	}{
		{
			name: "first 5 threads with empty array as tag",
			args: args{
				ctx:  ctx,
				tags: []string{""},
				query: entity.SliceQuery{
					After: algo.NullString(""),
					Limit: 5,
				},
			},
			want: &entity.ThreadSlice{
				Threads: []*entity.Thread{},
				SliceInfo: &entity.SliceInfo{
					FirstCursor: "",
					LastCursor:  "",
					HasNext:     false,
				},
			},
		},
		{
			name: "first 5 threads with a maintag",
			args: args{
				ctx:  ctx,
				tags: []string{"MainA"},
				query: entity.SliceQuery{
					After: algo.NullString(""),
					Limit: 5,
				},
			},
			want: &entity.ThreadSlice{
				Threads: []*entity.Thread{
					threads[5], threads[4], threads[3], threads[2], threads[1],
				},
				SliceInfo: &entity.SliceInfo{
					FirstCursor: threads[5].ID.ToBase64String(),
					LastCursor:  threads[1].ID.ToBase64String(),
					HasNext:     true,
				},
			},
		},
		{
			name: "first 5 threads with two tags",
			args: args{
				ctx:  ctx,
				tags: []string{"MainC", "SubB"},
				query: entity.SliceQuery{
					After: algo.NullString(""),
					Limit: 5,
				},
			},
			want: &entity.ThreadSlice{
				Threads: []*entity.Thread{
					threads[9], threads[8], threads[7], threads[6], threads[5],
				},
				SliceInfo: &entity.SliceInfo{
					FirstCursor: threads[9].ID.ToBase64String(),
					LastCursor:  threads[5].ID.ToBase64String(),
					HasNext:     true,
				},
			},
		},
		{
			name: "last 5 threads with only subtag",
			args: args{
				ctx:  ctx,
				tags: []string{"SubC"},
				query: entity.SliceQuery{
					Before: algo.NullString(""),
					Limit:  5,
				},
			},
			want: &entity.ThreadSlice{
				Threads: []*entity.Thread{
					threads[9], threads[8], threads[7], threads[6],
				},
				SliceInfo: &entity.SliceInfo{
					FirstCursor: threads[9].ID.ToBase64String(),
					LastCursor:  threads[6].ID.ToBase64String(),
					HasNext:     false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.SearchThreads(tt.args.ctx, tt.args.tags, tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.SearchThreads() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got.SliceInfo, tt.want.SliceInfo); diff != "" {
				t.Errorf("Service.SearchThreads().SliceInfo diff: %s", diff)
			}
			if len(got.Threads) != len(tt.want.Threads) {
				t.Errorf("Service.SearchThreads().len(Threads) = %v, want %v", len(got.Threads), len(tt.want.Threads))
			}
			for i, thread := range got.Threads {
				thread.CreatedAt = tt.want.Threads[i].CreatedAt
				if diff := cmp.Diff(thread, tt.want.Threads[i]); diff != "" {
					t.Errorf("Service.SearchThreads().Threads[%v] diff: %s", i, diff)
				}
			}
		})
	}
}

func TestService_GetThreadReplies(t *testing.T) {
	mainTags := []string{"MainA", "MainB", "MainC"}
	service, _ := initEnv(t, mainTags...)

	thread, ctx := pubThread(t, service, testUser{email: "t@example", name: "a"})
	var posts []*entity.Post
	for i := 0; i < 10; i++ {
		post, _ := pubPost(t, service, testUser{email: "p@example"}, thread.ID)
		posts = append(posts, post)
	}

	type args struct {
		ctx    context.Context
		thread *entity.Thread
		query  entity.SliceQuery
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.PostSlice
		wantErr bool
	}{
		{
			name: "first 5",
			args: args{
				ctx:    ctx,
				thread: thread,
				query: entity.SliceQuery{
					After: algo.NullString(""),
					Limit: 5,
				},
			},
			want: &entity.PostSlice{
				Posts: []*entity.Post{
					posts[0], posts[1], posts[2], posts[3], posts[4],
				},
				SliceInfo: &entity.SliceInfo{
					FirstCursor: posts[0].ID.ToBase64String(),
					LastCursor:  posts[4].ID.ToBase64String(),
					HasNext:     true,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetThreadReplies(tt.args.ctx, tt.args.thread, tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetThreadReplies() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.GetThreadReplies() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_GetThreadReplyCount(t *testing.T) {
	mainTags := []string{"MainA", "MainB", "MainC"}
	service, _ := initEnv(t, mainTags...)

	thread, ctx := pubThread(t, service, testUser{email: "t@example", name: "a"})
	for i := 0; i < 10; i++ {
		pubPost(t, service, testUser{email: "p@example"}, thread.ID)
	}

	type args struct {
		ctx    context.Context
		thread *entity.Thread
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				ctx:    ctx,
				thread: thread,
			},
			want:    10,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetThreadReplyCount(tt.args.ctx, tt.args.thread)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetThreadReplyCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Service.GetThreadReplyCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_GetThreadCatalog(t *testing.T) {
	mainTags := []string{"MainA", "MainB", "MainC"}
	service, _ := initEnv(t, mainTags...)

	thread, ctx := pubThread(t, service, testUser{email: "t@example", name: "a"})
	var catalog []*entity.ThreadCatalogItem
	for i := 0; i < 10; i++ {
		post, _ := pubPost(t, service, testUser{email: "p@example"}, thread.ID)
		catalog = append(catalog, &entity.ThreadCatalogItem{
			PostID:    post.ID,
			CreatedAt: post.CreatedAt,
		})
	}

	type args struct {
		ctx    context.Context
		thread *entity.Thread
	}
	tests := []struct {
		name    string
		args    args
		want    []*entity.ThreadCatalogItem
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				ctx:    ctx,
				thread: thread,
			},
			want:    catalog,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetThreadCatalog(tt.args.ctx, tt.args.thread)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetThreadCatalog() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.GetThreadCatalog() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_PubPost(t *testing.T) {
	mainTags := []string{"MainA", "MainB", "MainC"}
	service, _ := initEnv(t, mainTags...)

	thread, _ := pubThread(t, service, testUser{email: "a@example.com", name: "a"})
	post, _ := pubPost(t, service, testUser{email: "a@example.com", name: "a"}, thread.ID)
	user1, userCtx1 := loginUser(t, service, testUser{email: "p1@example.com"})
	user2, userCtx2 := loginUser(t, service, testUser{email: "p2@example.com", name: "p2"})
	userG, userCtxG := loginUser(t, service, testUser{})
	type args struct {
		ctx  context.Context
		post entity.PostInput
	}
	type quotedChecker struct {
		quotedCount int
	}
	tests := []struct {
		name        string
		args        args
		want        *entity.Post
		wantErrIs   error
		checkQuoted *quotedChecker
	}{
		{
			name: "anonymous signed in user",
			args: args{
				ctx: userCtx1,
				post: entity.PostInput{
					ThreadID:  thread.ID,
					Anonymous: true,
					Content:   "content1",
				},
			},
			want: &entity.Post{
				Author: &entity.Author{
					UserID:    user1.ID,
					Anonymous: true,
					Author:    user1.ID.ToBase64String(),
				},
				Content:  "content1",
				ThreadID: thread.ID,
			},
		},
		{
			name: "pub post with user name",
			args: args{
				ctx: userCtx2,
				post: entity.PostInput{
					ThreadID:  thread.ID,
					Anonymous: false,
					Content:   "content2",
				},
			},
			want: &entity.Post{
				Author: &entity.Author{
					UserID:    user2.ID,
					Anonymous: false,
					Author:    *user2.Name,
				},
				Content:  "content2",
				ThreadID: thread.ID,
			},
		},
		{
			name: "check duplicated",
			args: args{
				ctx: userCtx2,
				post: entity.PostInput{
					ThreadID:  thread.ID,
					Anonymous: false,
					Content:   "content2",
				},
			},
			wantErrIs: errors.Duplicated.New(),
		},
		{
			name: "pub post with quoted post",
			args: args{
				ctx: userCtx1,
				post: entity.PostInput{
					ThreadID:  thread.ID,
					Anonymous: true,
					Content:   "content3",
					QuoteIds:  []uid.UID{post.ID},
				},
			},
			want: &entity.Post{
				Author: &entity.Author{
					UserID:    user1.ID,
					Anonymous: true,
					Author:    user1.ID.ToBase64String(),
				},
				Content:  "content3",
				ThreadID: thread.ID,
				QuoteIDs: []uid.UID{post.ID},
			},
			checkQuoted: &quotedChecker{
				quotedCount: 1,
			},
		},
		{
			name: "pub post with guest user",
			args: args{
				ctx: userCtxG,
				post: entity.PostInput{
					ThreadID:  thread.ID,
					Anonymous: true,
					Content:   "contentG",
				},
			},
			want: &entity.Post{
				Author: &entity.Author{
					UserID:    userG.ID,
					Guest:     true,
					Anonymous: true,
					Author:    userG.ID.ToBase64String(),
				},
				Content:  "contentG",
				ThreadID: thread.ID,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.PubPost(tt.args.ctx, tt.args.post)
			if (tt.wantErrIs == nil && err != nil) || (tt.wantErrIs != nil && !errors.Is(err, tt.wantErrIs)) {
				t.Errorf("Service.PubThread() error = %v, wantErr %v", err, tt.wantErrIs)
			}
			if err != nil {
				return
			}
			tt.want.ID = got.ID
			tt.want.CreatedAt = got.CreatedAt
			if tt.args.post.Anonymous {
				tt.want.Author.Author = got.Author.Author
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("Service.PubPost() missmatch %s", diff)
			}
			if len(tt.want.QuoteIDs) != 0 {
				quoted, err := service.Repo.Post.QuotedPosts(tt.args.ctx, got)
				if err != nil {
					t.Error(errors.Wrap(err, "Quotes()"))
				}
				if len(quoted) == 0 {
					t.Error("should have a quoted post")
				} else if diff := cmp.Diff(quoted[0], post); diff != "" {
					t.Errorf("Service.PubPost().Quotes() missmatch: %s", diff)
				}
			}
			if tt.checkQuoted != nil {
				gotCount, err := service.Repo.Post.QuotedCount(tt.args.ctx, post)
				if err != nil {
					t.Errorf("quotedPost.QuotedCount error: %v", err)
				}
				if gotCount != tt.checkQuoted.quotedCount {
					t.Errorf("Post(%v).QuotedCount()=%v, want=%v", post, gotCount, tt.checkQuoted.quotedCount)
				}
			}
		})
	}
}

func TestService_BlockPost(t *testing.T) {
	mainTags := []string{"MainA", "MainB", "MainC"}
	service, ctx := initEnv(t, mainTags...)

	thread, _ := pubThread(t, service, testUser{email: "t@example.com"})
	oriPost, _ := pubPost(t, service, testUser{email: "p@example.com"}, thread.ID)
	mod, _ := loginUser(t, service, testUser{email: "mod@example.com"})
	mod.Role = entity.RoleMod
	_, err := service.Repo.User.Update(ctx, mod)
	if err != nil {
		t.Fatal(err)
	}
	oriPost.Blocked = true
	oriPost.Content = entity.BlockedContent
	_, modCtx := loginUser(t, service, testUser{email: "mod@example.com"})

	t.Run("check post in memory", func(t *testing.T) {
		post, err := service.BlockPost(modCtx, oriPost.ID)
		if err != nil {
			t.Fatal(errors.Wrap(err, "BlockPost"))
		}
		oriPost.CreatedAt = post.CreatedAt
		if diff := cmp.Diff(post, oriPost); diff != "" {
			t.Errorf("BlockPost() post matched: %s", diff)
		}
	})
	t.Run("check post in database", func(t *testing.T) {
		post, err := service.GetPostByID(modCtx, oriPost.ID)
		if err != nil {
			t.Fatal(errors.Wrap(err, "GetPostByID"))
		}
		oriPost.CreatedAt = post.CreatedAt
		if diff := cmp.Diff(post, oriPost); diff != "" {
			t.Errorf("BlockPost() post matched: %s", diff)
		}
	})
}

func TestService_GetPostQuotedPosts(t *testing.T) {
	mainTags := []string{"MainA", "MainB", "MainC"}
	service, _ := initEnv(t, mainTags...)

	thread, _ := pubThread(t, service, testUser{email: "t@example", name: "a"})
	post1, _ := pubPost(t, service, testUser{email: "1@example"}, thread.ID)
	post2, _ := pubPost(t, service, testUser{email: "2@example"}, thread.ID)
	post3, _ := pubPost(t, service, testUser{email: "3@example"}, thread.ID)
	post4, ctx := pubPost(t, service, testUser{email: "4@example"}, thread.ID, post2.ID, post1.ID, post3.ID)

	type args struct {
		ctx  context.Context
		post *entity.Post
	}
	tests := []struct {
		name    string
		args    args
		want    []*entity.Post
		wantErr bool
	}{
		{
			name: "quote 3, test order",
			args: args{
				ctx:  ctx,
				post: post4,
			},
			want:    []*entity.Post{post2, post1, post3},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetPostQuotedPosts(tt.args.ctx, tt.args.post)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetPostQuotedPosts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.GetPostQuotedPosts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_GetPostQuotedCount(t *testing.T) {
	mainTags := []string{"MainA", "MainB", "MainC"}
	service, _ := initEnv(t, mainTags...)

	thread, _ := pubThread(t, service, testUser{email: "t@example", name: "a"})
	post1, ctx := pubPost(t, service, testUser{email: "1@example"}, thread.ID)
	pubPost(t, service, testUser{email: "2@example"}, thread.ID, post1.ID)
	pubPost(t, service, testUser{email: "3@example"}, thread.ID, post1.ID)

	type args struct {
		ctx  context.Context
		post *entity.Post
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				ctx:  ctx,
				post: post1,
			},
			want:    2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetPostQuotedCount(tt.args.ctx, tt.args.post)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetPostQuotedCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Service.GetPostQuotedCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_SearchTags(t *testing.T) {
	mainTags := []string{"MainA", "MainB", "MainC"}
	service, ctx := initEnv(t, mainTags...)

	pubThreadWithTags(t, service, testUser{email: "a@example.com"}, "MainA", []string{"Sub11", "Sub21"})
	pubThreadWithTags(t, service, testUser{email: "a@example.com"}, "MainB", []string{"Sub12", "Sub22"})
	pubThreadWithTags(t, service, testUser{email: "a@example.com"}, "MainC", []string{"Sub13", "Sub23"})
	pubThreadWithTags(t, service, testUser{email: "a@example.com"}, "MainA", []string{"Sub14", "Sub24"})
	pubThreadWithTags(t, service, testUser{email: "a@example.com"}, "MainB", []string{"Sub15", "Sub25"})
	pubThreadWithTags(t, service, testUser{email: "a@example.com"}, "MainC", []string{"Sub16", "Sub26"})
	type args struct {
		ctx   context.Context
		query *string
		limit *int
	}
	tests := []struct {
		name        string
		args        args
		ignoreOrder bool
		want        []*entity.Tag
		wantErr     bool
	}{
		{
			name: "search all tags",
			args: args{
				ctx:   ctx,
				query: nil,
				limit: algo.NullInt(9),
			},
			want: []*entity.Tag{
				{Name: "MainA", IsMain: true},
				{Name: "MainB", IsMain: true},
				{Name: "MainC", IsMain: true},
				{Name: "Sub14", IsMain: false},
				{Name: "Sub15", IsMain: false},
				{Name: "Sub16", IsMain: false},
				{Name: "Sub24", IsMain: false},
				{Name: "Sub25", IsMain: false},
				{Name: "Sub26", IsMain: false},
			},
			ignoreOrder: true,
		},
		{
			name: "search mainTags tags",
			args: args{
				ctx:   ctx,
				query: algo.NullString("Main"),
				limit: algo.NullInt(10),
			},
			want: []*entity.Tag{
				{Name: "MainC", IsMain: true},
				{Name: "MainB", IsMain: true},
				{Name: "MainA", IsMain: true},
			},
		},
		{
			name: "search sub tags",
			args: args{
				ctx:   ctx,
				query: algo.NullString("ub1"),
				limit: algo.NullInt(10),
			},
			want: []*entity.Tag{
				{Name: "Sub16", IsMain: false},
				{Name: "Sub15", IsMain: false},
				{Name: "Sub14", IsMain: false},
				{Name: "Sub13", IsMain: false},
				{Name: "Sub12", IsMain: false},
				{Name: "Sub11", IsMain: false},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.SearchTags(tt.args.ctx, tt.args.query, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.SearchTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.ignoreOrder {
				if diff := cmp.Diff(got, tt.want); diff != "" {
					t.Errorf("Service.SearchTags() missmatch: %s", diff)
				}
			} else {
				if diff := cmp.Diff(got, tt.want, tagSetCmp); diff != "" {
					t.Errorf("Service.SearchTags() missmatch: %s", diff)
				}
			}
		})
	}
}

func TestService_GetUnreadNotiCount(t *testing.T) {
	mainTags := []string{"MainA", "MainB", "MainC"}
	service, ctx := initEnv(t, mainTags...)

	user, userCtx := loginUser(t, service, testUser{email: "a@example.com"})
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name       string
		args       args
		beforeTest func(t *testing.T)
		want       int
		wantErr    bool
	}{
		{
			name: "1 system noti(welcome)",
			args: args{
				ctx: userCtx,
			},
			want: 1,
		},
		{
			name: "1 global system noti",
			args: args{
				ctx: userCtx,
			},
			beforeTest: func(t *testing.T) {
				noti, err := entity.NewSystemNoti("welcome!", "welcome to abyss", entity.SendToGroup(entity.AllUser))
				if err != nil {
					t.Fatal(err)
				}
				if err := service.Repo.Noti.Insert(ctx, noti); err != nil {
					t.Fatal(err)
				}
			},
			want: 2,
		},
		{
			name: "count after read",
			args: args{
				ctx: userCtx,
			},
			beforeTest: func(t *testing.T) {
				if _, err := service.GetNotifications(userCtx, entity.SliceQuery{
					After: algo.NullString(""),
					Limit: 5,
				}); err != nil {
					t.Fatal(err)
				}
			},
			want: 0,
		},
		{
			name: "new reply noti",
			args: args{
				ctx: userCtx,
			},
			beforeTest: func(t *testing.T) {
				thread, _ := pubThread(t, service, testUser{email: *user.Email})
				pubPost(t, service, testUser{email: "p@example.com"}, thread.ID)
				time.Sleep(100 * time.Millisecond)
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.beforeTest != nil {
				tt.beforeTest(t)
			}
			got, err := service.GetUnreadNotiCount(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetUnreadNotiCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.GetUnreadNotiCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_GetNotification(t *testing.T) {
	mainTags := []string{"MainA", "MainB", "MainC"}
	service, ctx := initEnv(t, mainTags...)

	user, _ := loginUser(t, service, testUser{email: "t@example.com"}) // One Welcome Noti
	thread, _ := pubThread(t, service, testUser{email: *user.Email})

	type args struct {
		email string
		query entity.SliceQuery
	}
	tests := []struct {
		name       string
		args       args
		beforeTest func(t *testing.T, want *entity.NotiSlice)
		want       *entity.NotiSlice
		wantErr    bool
	}{
		{
			name: "2 system noti 1 replied noti",
			args: args{
				email: *user.Email,
				query: entity.SliceQuery{After: algo.NullString(""), Limit: 5},
			},
			beforeTest: func(t *testing.T, want *entity.NotiSlice) {
				noti, err := entity.NewSystemNoti("Hi everyone!", "Let's Party!", entity.SendToGroup(entity.AllUser))
				if err != nil {
					t.Fatal(err)
				}
				if err := service.Repo.Noti.Insert(ctx, noti); err != nil {
					t.Fatal(err)
				}
				post, _ := pubPost(t, service, testUser{email: "p1@example.com"}, thread.ID)
				content := want.Notifications[0].Content.(entity.RepliedNoti)
				content.FirstReplyID = post.ID
				want.Notifications[0].Content = content
				time.Sleep(100 * time.Millisecond)
			},
			want: &entity.NotiSlice{
				SliceInfo: &entity.SliceInfo{
					HasNext: false,
				},
				Notifications: []*entity.Notification{
					{
						Type: entity.NotiTypeReplied,
						Content: entity.RepliedNoti{
							Thread: &entity.ThreadOutline{
								ID:      thread.ID,
								Title:   thread.Title,
								Content: thread.Content,
								MainTag: thread.MainTag,
								SubTags: thread.SubTags,
							},
							NewRepliesCount: 1,
						},
						Key:       fmt.Sprintf("replied:%s", thread.ID.ToBase64String()),
						Receivers: []entity.Receiver{entity.SendToUser(user.ID)},
					},
					{
						Type: entity.NotiTypeSystem,
						Content: entity.SystemNoti{
							Title:   "Hi everyone!",
							Content: "Let's Party!",
						},
						Receivers: []entity.Receiver{entity.SendToGroup(entity.AllUser)},
					},
					{
						Type: entity.NotiTypeSystem,
						Content: entity.SystemNoti{
							Title:   entity.WelcomeTitle,
							Content: entity.WelcomeContent,
						},
						Receivers: []entity.Receiver{entity.SendToUser(user.ID)},
					},
				},
			},
		},
		{
			name: "3 quoted noti, 1 replied",
			args: args{
				email: *user.Email,
				query: entity.SliceQuery{After: algo.NullString(""), Limit: 4},
			},
			beforeTest: func(t *testing.T, want *entity.NotiSlice) {
				qp1, _ := pubPost(t, service, testUser{email: *user.Email}, thread.ID)
				qp2, _ := pubPost(t, service, testUser{email: *user.Email}, thread.ID)
				p1, _ := pubPost(t, service, testUser{email: "p1@example.com"}, thread.ID, qp1.ID)
				p2, _ := pubPost(t, service, testUser{email: "p1@example.com"}, thread.ID, qp1.ID, qp2.ID)
				writeWant := func(noti *entity.Notification, q, p *entity.Post) {
					noti.Content = entity.QuotedNoti{
						ThreadID: thread.ID,
						QuotedPost: &entity.PostOutline{
							ID: q.ID,
							Author: &entity.Author{
								// user id won't return to frontend
								Anonymous: q.Author.Anonymous,
								Author:    q.Author.Author,
							},
							Content: q.Content,
						},
						Post: &entity.PostOutline{
							ID: p.ID,
							Author: &entity.Author{
								// user id won't return to frontend
								Anonymous: p.Author.Anonymous,
								Author:    p.Author.Author,
							},
							Content: p.Content,
						},
					}
					noti.Key = fmt.Sprintf("quoted:%s:%s", q.ID.ToBase64String(), p.ID.ToBase64String())
				}
				writeWant(want.Notifications[0], qp2, p2)
				writeWant(want.Notifications[1], qp1, p2)
				content := want.Notifications[2].Content.(entity.RepliedNoti)
				content.FirstReplyID = p1.ID
				content.NewRepliesCount = 2
				want.Notifications[2].Content = content
				writeWant(want.Notifications[3], qp1, p1)
			},
			want: &entity.NotiSlice{
				SliceInfo: &entity.SliceInfo{
					HasNext: true,
				},
				Notifications: []*entity.Notification{
					{ // 0
						Type:      entity.NotiTypeQuoted,
						Receivers: []entity.Receiver{entity.SendToUser(user.ID)},
					},
					{ // 1
						Type:      entity.NotiTypeQuoted,
						Receivers: []entity.Receiver{entity.SendToUser(user.ID)},
					},
					{ // 2
						Type: entity.NotiTypeReplied,
						Content: entity.RepliedNoti{
							Thread: &entity.ThreadOutline{
								ID:      thread.ID,
								Title:   thread.Title,
								Content: thread.Content,
								MainTag: thread.MainTag,
								SubTags: thread.SubTags,
							},
							NewRepliesCount: 1,
						},
						Key:       fmt.Sprintf("replied:%s", thread.ID.ToBase64String()),
						Receivers: []entity.Receiver{entity.SendToUser(user.ID)},
					},
					{ // 3
						Type:      entity.NotiTypeQuoted,
						Receivers: []entity.Receiver{entity.SendToUser(user.ID)},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.beforeTest != nil {
				tt.beforeTest(t, tt.want)
			}
			user, ctx := loginUser(t, service, testUser{email: tt.args.email})
			t.Logf("get user context: %#v, %#v", user, ctx)
			got, err := service.GetNotifications(ctx, tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetNotification() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// sliceInfo
			wantSliceInfo := &entity.SliceInfo{HasNext: tt.want.SliceInfo.HasNext}
			if len(got.Notifications) > 0 {
				wantSliceInfo.FirstCursor = got.Notifications[0].SortKey.ToBase64String()
				wantSliceInfo.LastCursor = got.Notifications[len(got.Notifications)-1].SortKey.ToBase64String()
			}
			if diff := cmp.Diff(got.SliceInfo, wantSliceInfo); diff != "" {
				t.Errorf("Service.GetNotification().SliceInfo missmatch: %s", diff)
			}
			// notifications.count
			if len(got.Notifications) != len(tt.want.Notifications) {
				t.Errorf(
					"Service.GetNotification().count missmatch, got=%v, want=%v",
					len(got.Notifications), tt.want.Notifications,
				)
			}
			// notifications
			for i := range got.Notifications {
				g, w := got.Notifications[i], tt.want.Notifications[i]
				w.EventTime, w.SortKey = g.EventTime, g.SortKey
				if w.Key == "" {
					w.Key = g.Key
				}
				if diff := cmp.Diff(g, w); diff != "" {
					t.Errorf("Service.GetNotification().Notifications[%v] missmatch: %s", i, diff)
				}
			}
		})
	}
}
