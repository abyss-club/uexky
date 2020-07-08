package uexky

import (
	"context"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky/lib/algo"
	"gitlab.com/abyss.club/uexky/lib/uid"
	"gitlab.com/abyss.club/uexky/uexky/entity"
)

func TestService_LoginByEmail(t *testing.T) {
	email := "user1@example.com"
	service, err := InitDevService()
	getNewDBCtx(t)
	if err != nil {
		t.Fatal(err)
	}
	user, _ := loginUser(t, service, testUser{email: email})
	wantUser := &entity.User{
		Email: email,
		Role:  entity.RoleNormal,
		Repo:  service.User.Repo,
		ID:    1,
	}
	if !reflect.DeepEqual(user, wantUser) {
		t.Errorf("want user %+v, bug got %+v", wantUser, user)
	}
}

func TestService_SetUserName(t *testing.T) {
	service, err := InitDevService()
	if err != nil {
		t.Fatal(err)
	}
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
				ctx:   getNewDBCtx(t),
				name:  "tom",
				email: "tom@example.com",
			},
			wantName: "tom",
			wantErr:  false,
		},
		{
			name: "already has name",
			args: args{
				ctx:   getNewDBCtx(t),
				name:  "tom2",
				email: "tom@example.com",
			},
			wantName: "tom",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, ctx := loginUser(t, service, testUser{email: tt.args.email})
			gotUser, err := service.SetUserName(ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.SetUserName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUser.Name == nil || *gotUser.Name != tt.wantName {
				t.Errorf("Service.SetUserName() = %v, want %v", gotUser.Name, tt.wantName)
			}
			if user.Name != gotUser.Name {
				t.Errorf("user name not sync origin object")
			}
		})
	}
}

func TestService_GetUserThreads(t *testing.T) {
	service, err := InitDevService()
	if err != nil {
		t.Fatal(err)
	}
	ctx := getNewDBCtx(t)
	mainTags := []string{"MainA", "MainB", "MainC"}
	if err := service.SetMainTags(ctx, mainTags); err != nil {
		t.Fatal(err)
	}
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
				if diff := cmp.Diff(thread, tt.want.Threads[i], forumRepoComp); diff != "" {
					t.Errorf("Service.GetUserThreads().Threads[%v] diff: %s", i, diff)
				}
			}
		})
	}
}

func TestService_GetUserPosts(t *testing.T) {
	service, err := InitDevService()
	if err != nil {
		t.Fatal(err)
	}
	ctx := getNewDBCtx(t)
	mainTags := []string{"MainA", "MainB", "MainC"}
	if err := service.SetMainTags(ctx, mainTags); err != nil {
		t.Fatal(err)
	}
	var threads []*entity.Thread
	for i := 0; i < 2; i++ {
		thread, _ := pubThread(t, service, testUser{email: "thread@example.com"})
		if err != nil {
			t.Fatal(err)
		}
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
				if diff := cmp.Diff(post, tt.want.Posts[i], forumRepoComp); diff != "" {
					t.Errorf("Service.GetUserPosts().Posts[%v] diff: %s", i, diff)
				}
			}
		})
	}
}

func TestService_SyncUserTags(t *testing.T) {
	ctx := getNewDBCtx(t)
	service, err := InitDevService()
	if err != nil {
		t.Fatal(err)
	}
	if err := service.SetMainTags(ctx, []string{"MainA", "MainB", "MainC"}); err != nil {
		t.Fatal(err)
	}
	_, ctx = loginUser(t, service, testUser{name: "a", email: "a@example.com"})
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
	ctx := getNewDBCtx(t)
	service, err := InitDevService()
	if err != nil {
		t.Fatal(err)
	}
	if err := service.SetMainTags(ctx, []string{"MainA", "MainB", "MainC"}); err != nil {
		t.Fatal(err)
	}
	_, ctx = loginUser(t, service, testUser{name: "a", email: "a@example.com"})
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
	ctx := getNewDBCtx(t)
	service, err := InitDevService()
	if err != nil {
		t.Fatal(err)
	}
	if err := service.SetMainTags(ctx, []string{"MainA", "MainB", "MainC"}); err != nil {
		t.Fatal(err)
	}
	_, ctx = loginUser(t, service, testUser{name: "a", email: "a@example.com"})
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
	service, err := InitDevService()
	if err != nil {
		t.Fatal(err)
	}
	ctx := getNewDBCtx(t)
	if err := service.SetMainTags(ctx, []string{"MainA", "MainB", "MainC"}); err != nil {
		t.Fatal(err)
	}
	thread, _ := pubThread(t, service, testUser{email: "t@example.com"})
	post, _ := pubPost(t, service, testUser{email: "p@example.com"}, thread.ID)
	mod, _ := loginUser(t, service, testUser{email: "mod@example.com"})
	if err := service.User.Repo.UpdateUser(ctx, mod.ID, &entity.UserUpdate{
		Role: (*entity.Role)(algo.NullString(string(entity.RoleMod))),
	}); err != nil {
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

func TestService_BlockPost(t *testing.T) {
	service, err := InitDevService()
	if err != nil {
		t.Fatal(err)
	}
	ctx := getNewDBCtx(t)
	if err := service.SetMainTags(ctx, []string{"MainA", "MainB", "MainC"}); err != nil {
		t.Fatal(err)
	}
	thread, _ := pubThread(t, service, testUser{email: "t@example.com"})

	oriPost, _ := pubPost(t, service, testUser{email: "p@example.com"}, thread.ID)
	mod, _ := loginUser(t, service, testUser{email: "mod@example.com"})
	if err := service.User.Repo.UpdateUser(ctx, mod.ID, &entity.UserUpdate{
		Role: (*entity.Role)(algo.NullString(string(entity.RoleMod))),
	}); err != nil {
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
		if diff := cmp.Diff(post, oriPost, forumRepoComp); diff != "" {
			t.Errorf("BlockPost() post matched: %s", diff)
		}
	})
	t.Run("check post in database", func(t *testing.T) {
		post, err := service.GetPostByID(modCtx, oriPost.ID)
		if err != nil {
			t.Fatal(errors.Wrap(err, "GetPostByID"))
		}
		oriPost.CreatedAt = post.CreatedAt
		if diff := cmp.Diff(post, oriPost, forumRepoComp); diff != "" {
			t.Errorf("BlockPost() post matched: %s", diff)
		}
	})
}

func TestService_LockThread(t *testing.T) {
	service, err := InitDevService()
	if err != nil {
		t.Fatal(err)
	}
	ctx := getNewDBCtx(t)
	if err := service.SetMainTags(ctx, []string{"MainA", "MainB", "MainC"}); err != nil {
		t.Fatal(err)
	}
	oriThread, _ := pubThread(t, service, testUser{email: "t@example.com"})

	mod, _ := loginUser(t, service, testUser{email: "mod@example.com"})
	if err := service.User.Repo.UpdateUser(ctx, mod.ID, &entity.UserUpdate{
		Role: (*entity.Role)(algo.NullString(string(entity.RoleMod))),
	}); err != nil {
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
		if diff := cmp.Diff(thread, oriThread, forumRepoComp); diff != "" {
			t.Errorf("LockThread() not matched: %s", diff)
		}
	})
	t.Run("check thread in database", func(t *testing.T) {
		thread, err := service.GetThreadByID(modCtx, oriThread.ID)
		if err != nil {
			t.Fatal(errors.Wrap(err, "GetThreadByID"))
		}
		oriThread.CreatedAt = thread.CreatedAt
		if diff := cmp.Diff(thread, oriThread, forumRepoComp); diff != "" {
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
	service, err := InitDevService()
	if err != nil {
		t.Fatal(err)
	}
	ctx := getNewDBCtx(t)
	if err := service.SetMainTags(ctx, []string{"MainA", "MainB", "MainC"}); err != nil {
		t.Fatal(err)
	}
	oriThread, _ := pubThread(t, service, testUser{email: "t@example.com"})

	mod, _ := loginUser(t, service, testUser{email: "mod@example.com"})
	if err := service.User.Repo.UpdateUser(ctx, mod.ID, &entity.UserUpdate{
		Role: (*entity.Role)(algo.NullString(string(entity.RoleMod))),
	}); err != nil {
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
		if diff := cmp.Diff(thread, oriThread, forumRepoComp); diff != "" {
			t.Errorf("LockThread() not matched: %s", diff)
		}
	})
	t.Run("check thread in database", func(t *testing.T) {
		thread, err := service.GetThreadByID(modCtx, oriThread.ID)
		if err != nil {
			t.Fatal(errors.Wrap(err, "GetThreadByID"))
		}
		oriThread.CreatedAt = thread.CreatedAt
		if diff := cmp.Diff(thread, oriThread, forumRepoComp); diff != "" {
			t.Errorf("LockThread() not matched: %s", diff)
		}
	})
}

func TestService_EditTags(t *testing.T) {
	service, err := InitDevService()
	if err != nil {
		t.Fatal(err)
	}
	ctx := getNewDBCtx(t)
	if err := service.SetMainTags(ctx, []string{"MainA", "MainB", "MainC"}); err != nil {
		t.Fatal(err)
	}
	oriThread, _ := pubThread(t, service, testUser{email: "t@example.com"})

	mod, _ := loginUser(t, service, testUser{email: "mod@example.com"})
	if err := service.User.Repo.UpdateUser(ctx, mod.ID, &entity.UserUpdate{
		Role: (*entity.Role)(algo.NullString(string(entity.RoleMod))),
	}); err != nil {
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
		if diff := cmp.Diff(thread, oriThread, forumRepoComp); diff != "" {
			t.Errorf("EditTags() not matched: %s", diff)
		}
	})
	t.Run("check thread in database", func(t *testing.T) {
		thread, err := service.GetThreadByID(modCtx, oriThread.ID)
		if err != nil {
			t.Fatal(errors.Wrap(err, "GetThreadByID"))
		}
		oriThread.CreatedAt = thread.CreatedAt
		if diff := cmp.Diff(thread, oriThread, forumRepoComp); diff != "" {
			t.Errorf("EditTags() not matched: %s", diff)
		}
	})
}

func TestService_PubThread(t *testing.T) {
	service, err := InitDevService()
	if err != nil {
		t.Fatal(err)
	}
	getNewDBCtx(t)
	user, ctx := loginUser(t, service, testUser{email: "a@example.com", name: "a"})
	if err := service.SetMainTags(ctx, []string{"MainA", "MainB", "MainC"}); err != nil {
		t.Fatal(err)
	}
	type args struct {
		ctx    context.Context
		thread entity.ThreadInput
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.Thread
		wantErr bool
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
				},
			},
			want: &entity.Thread{
				Anonymous: true,
				Title:     nil,
				Content:   "content",
				MainTag:   "MainA",
				SubTags:   []string{"SubA", "SubB"},
				Repo:      service.Forum.Repo,
				AuthorObj: entity.Author{
					UserID: user.ID,
				},
			},
		},
		{
			name: "pub thread with user name",
			args: args{
				ctx,
				entity.ThreadInput{
					Anonymous: false,
					Content:   "content",
					MainTag:   "MainA",
					SubTags:   []string{"SubA", "SubB", "SubC"},
				},
			},
			want: &entity.Thread{
				Anonymous: false,
				Title:     nil,
				Content:   "content",
				MainTag:   "MainA",
				SubTags:   []string{"SubA", "SubB", "SubC"},
				Repo:      service.Forum.Repo,
				AuthorObj: entity.Author{
					UserID:   user.ID,
					UserName: user.Name,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.PubThread(tt.args.ctx, tt.args.thread)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.PubThread() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tt.want.ID = got.ID
			tt.want.CreatedAt = got.CreatedAt
			tt.want.AuthorObj.AnonymousID = got.AuthorObj.AnonymousID
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.PubThread() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_SearchThreads(t *testing.T) {
	service, err := InitDevService()
	if err != nil {
		t.Fatal(err)
	}
	ctx := getNewDBCtx(t)
	mainTags := []string{"MainA", "MainB", "MainC"}
	if err := service.SetMainTags(ctx, mainTags); err != nil {
		t.Fatal(err)
	}
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
	_, ctx = loginUser(t, service, testUser{email: "a@example.com"})
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
				if diff := cmp.Diff(thread, tt.want.Threads[i], forumRepoComp); diff != "" {
					t.Errorf("Service.SearchThreads().Threads[%v] diff: %s", i, diff)
				}
			}
		})
	}
}

func TestService_PubPost(t *testing.T) {
	service, err := InitDevService()
	if err != nil {
		t.Fatal(err)
	}
	ctx := getNewDBCtx(t)
	mainTags := []string{"MainA", "MainB", "MainC"}
	if err := service.SetMainTags(ctx, mainTags); err != nil {
		t.Fatal(err)
	}
	thread, _ := pubThread(t, service, testUser{email: "a@example.com", name: "a"})
	post, _ := pubPost(t, service, testUser{email: "a@example.com", name: "a"}, thread.ID)
	user1, userCtx1 := loginUser(t, service, testUser{email: "p1@example.com"})
	user2, userCtx2 := loginUser(t, service, testUser{email: "p2@example.com", name: "p2"})
	type args struct {
		ctx  context.Context
		post entity.PostInput
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.Post
		wantErr bool
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
				Anonymous: true,
				Content:   "content1",
				Repo:      service.Forum.Repo,
				Data: entity.PostData{
					ThreadID: thread.ID,
					Author: entity.Author{
						UserID: user1.ID,
					},
					QuotePosts: []*entity.Post{},
				},
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
				Anonymous: false,
				Content:   "content2",
				Repo:      service.Forum.Repo,
				Data: entity.PostData{
					ThreadID: thread.ID,
					Author: entity.Author{
						UserID:   user2.ID,
						UserName: user2.Name,
					},
					QuotePosts: []*entity.Post{},
				},
			},
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
				Anonymous: true,
				Content:   "content3",
				Repo:      service.Forum.Repo,
				Data: entity.PostData{
					ThreadID: thread.ID,
					Author: entity.Author{
						UserID: user1.ID,
					},
					QuoteIDs:   []uid.UID{post.ID},
					QuotePosts: []*entity.Post{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.PubPost(tt.args.ctx, tt.args.post)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.PubPost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tt.want.ID = got.ID
			tt.want.CreatedAt = got.CreatedAt
			if tt.args.post.Anonymous {
				tt.want.Data.Author.AnonymousID = got.Data.Author.AnonymousID
			}
			if diff := cmp.Diff(got, tt.want, forumRepoComp); diff != "" {
				t.Errorf("Service.PubPost() missmatch %s", diff)
			}
			if got.Author() != tt.want.Author() {
				t.Errorf("Service.PubPost().Author = %s, want = %s", got.Author(), tt.want.Author())
			}
			if len(tt.want.Data.QuoteIDs) != 0 {
				quoted, err := got.Quotes(tt.args.ctx)
				if err != nil {
					t.Error(errors.Wrap(err, "Quotes()"))
				}
				if len(quoted) == 0 {
					t.Error("should have a quoted post")
				} else if diff := cmp.Diff(quoted[0], post, forumRepoComp, timeCmp); diff != "" {
					t.Errorf("Service.PubPost().Quotes() missmatch: %s", diff)
				}
			}
		})
	}
}

func TestService_SearchTags(t *testing.T) {
	service, err := InitDevService()
	if err != nil {
		t.Fatal(err)
	}
	ctx := getNewDBCtx(t)
	if err := service.SetMainTags(ctx, []string{"MainA", "MainB", "MainC"}); err != nil {
		t.Fatal(err)
	}
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
	service, err := InitDevService()
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
	service, err := InitDevService()
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		ctx   context.Context
		query entity.SliceQuery
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.NotiSlice
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetNotification(tt.args.ctx, tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetNotification() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.GetNotification() = %v, want %v", got, tt.want)
			}
		})
	}
}
