package uexky

import (
	"context"
	"log"
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
	user, _, err := loginUser(service, testUser{email: email})
	if err != nil {
		t.Fatal(err)
	}
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
			user, ctx, err := loginUser(service, testUser{email: tt.args.email})
			if err != nil {
				t.Fatal(err)
			}
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
		thread, _, err := pubThread(service, tu)
		if err != nil {
			t.Fatal(err)
		}
		threads = append(threads, thread)
	}
	user, ctx, err := loginUser(service, testUser{email: "a@example.com"})
	if err != nil {
		t.Fatal(err)
	}
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
		thread, _, err := pubThread(service, testUser{email: "thread@example.com"})
		if err != nil {
			t.Fatal(err)
		}
		threads = append(threads, thread)
	}
	tu := testUser{email: "post@example.com", name: " post"}
	var posts []*entity.Post
	for i := 0; i < 10; i++ {
		post, _, err := pubPost(service, tu, threads[i%2].ID, nil)
		if err != nil {
			t.Fatal(err)
		}
		posts = append(posts, post)
	}
	user, ctx, err := loginUser(service, tu)
	if err != nil {
		log.Fatal(err)
	}
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

func TestService_GetUserTags(t *testing.T) {
	service, err := InitDevService()
	if err != nil {
		t.Fatal(err)
	}
	ctx := getNewDBCtx(t)
	if err := service.SetMainTags(ctx, []string{"MainA", "MainB", "MainC"}); err != nil {
		t.Fatal(err)
	}
	user, ctx, err := loginUser(service, testUser{name: "a", email: "a@example.com"})
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		ctx context.Context
		obj *entity.User
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "new user's tags",
			args: args{
				ctx: ctx,
				obj: user,
			},
			want:    []string{"MainA", "MainB", "MainC"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetUserTags(tt.args.ctx, tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetUserTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.GetUserTags() = %v, want %v", got, tt.want)
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
	_, ctx, err = loginUser(service, testUser{name: "a", email: "a@example.com"})
	if err != nil {
		t.Fatal(err)
	}
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
			tags, err := service.GetUserTags(ctx, got)
			if err != nil {
				t.Error(errors.Wrap(err, "GetUserTags"))
			}
			if diff := cmp.Diff(tags, tt.wantTags); diff != "" {
				t.Errorf("Service.SyncUserTags() = %v, wantTags %v", tags, tt.wantTags)
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
	_, ctx, err = loginUser(service, testUser{name: "a", email: "a@example.com"})
	if err != nil {
		t.Fatal(err)
	}
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
			tags, err := service.GetUserTags(ctx, got)
			if err != nil {
				t.Error(errors.Wrap(err, "GetUserTags"))
			}
			if diff := cmp.Diff(tags, tt.wantTags); diff != "" {
				t.Errorf("Service.AddUserSubbedTag() = %v, wantTags %v", tags, tt.wantTags)
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
	_, ctx, err = loginUser(service, testUser{name: "a", email: "a@example.com"})
	if err != nil {
		t.Fatal(err)
	}
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
			tags, err := service.GetUserTags(ctx, got)
			if err != nil {
				t.Error(errors.Wrap(err, "GetUserTags"))
			}
			if diff := cmp.Diff(tags, tt.wantTags); diff != "" {
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
	type args struct {
		ctx      context.Context
		postID   *uid.UID
		threadID *uid.UID
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.BanUser(tt.args.ctx, tt.args.postID, tt.args.threadID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.BanUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Service.BanUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_BlockPost(t *testing.T) {
	service, err := InitDevService()
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		ctx    context.Context
		postID uid.UID
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.Post
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.BlockPost(tt.args.ctx, tt.args.postID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.BlockPost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.BlockPost() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_LockThread(t *testing.T) {
	service, err := InitDevService()
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		ctx      context.Context
		threadID uid.UID
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.Thread
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.LockThread(tt.args.ctx, tt.args.threadID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.LockThread() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.LockThread() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_BlockThread(t *testing.T) {
	service, err := InitDevService()
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		ctx      context.Context
		threadID uid.UID
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.Thread
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.BlockThread(tt.args.ctx, tt.args.threadID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.BlockThread() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.BlockThread() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_EditTags(t *testing.T) {
	service, err := InitDevService()
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		ctx      context.Context
		threadID uid.UID
		mainTag  string
		subTags  []string
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.Thread
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.EditTags(tt.args.ctx, tt.args.threadID, tt.args.mainTag, tt.args.subTags)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.EditTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.EditTags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_PubThread(t *testing.T) {
	service, err := InitDevService()
	if err != nil {
		t.Fatal(err)
	}
	getNewDBCtx(t)
	user, ctx, err := loginUser(service, testUser{email: "a@example.com", name: "a"})
	if err != nil {
		t.Fatal(err)
	}
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
		var err error
		switch {
		case i < 3:
			thread, _, err = pubThreadWithTags(service, tu, "MainA", []string{"SubA"})
		case i < 6:
			thread, _, err = pubThreadWithTags(service, tu, "MainA", []string{"SubA", "SubB"})
		default:
			thread, _, err = pubThreadWithTags(service, tu, "MainC", []string{"SubB", "SubC"})
		}
		if err != nil {
			t.Fatal(err)
		}
		threads = append(threads, thread)
	}
	_, ctx, err = loginUser(service, testUser{email: "a@example.com"})
	if err != nil {
		t.Fatal(err)
	}
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

func TestService_GetThreadByID(t *testing.T) {
	service, err := InitDevService()
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		ctx context.Context
		id  uid.UID
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.Thread
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetThreadByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetThreadByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.GetThreadByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_PubPost(t *testing.T) {
	/*
		service, err := InitDevService()
		if err != nil {
			t.Fatal(err)
		}
		ctx := getNewDBCtx(t)
		mainTags := []string{"MainA", "MainB", "MainC"}
		if err := service.SetMainTags(ctx, mainTags); err != nil {
			t.Fatal(err)
		}
		thread, ctx, err := pubThread(service, testUser{email: "a@example.com", name: "a"})
		type args struct {
			email string
			name  *string
			post  entity.PostInput
		}
		tests := []struct {
			name    string
			args    args
			want    *entity.Post
			wantErr bool
		}{}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				user, ctx := loginUser(service, testUser{email: tt.args.email})
				if tt.args.name != nil && user.name == nil {
					if _, err := service.SetUserName(ctx, *tt.args.name); err != nil {
						t.Fatal(err)
					}
				}
				got, err := service.PubPost(tt.args.ctx, tt.args.post)
				if (err != nil) != tt.wantErr {
					t.Errorf("Service.PubPost() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Service.PubPost() = %v, want %v", got, tt.want)
				}
			})
		}
	*/
}

func TestService_GetPostByID(t *testing.T) {
	service, err := InitDevService()
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		ctx context.Context
		id  uid.UID
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.Post
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetPostByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetPostByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.GetPostByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_SetMainTags(t *testing.T) {
	service, err := InitDevService()
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		ctx  context.Context
		tags []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := service.SetMainTags(tt.args.ctx, tt.args.tags); (err != nil) != tt.wantErr {
				t.Errorf("Service.SetMainTags() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_GetMainTags(t *testing.T) {
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
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetMainTags(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetMainTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.GetMainTags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_GetRecommendedTags(t *testing.T) {
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
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetRecommendedTags(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetRecommendedTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.GetRecommendedTags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_SearchTags(t *testing.T) {
	service, err := InitDevService()
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		ctx   context.Context
		query *string
		limit *int
	}
	tests := []struct {
		name    string
		args    args
		want    []*entity.Tag
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.SearchTags(tt.args.ctx, tt.args.query, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.SearchTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.SearchTags() = %v, want %v", got, tt.want)
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
		want    *entity.UnreadNotiCount
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
		ctx     context.Context
		typeArg entity.NotiType
		query   entity.SliceQuery
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
			got, err := service.GetNotification(tt.args.ctx, tt.args.typeArg, tt.args.query)
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

func TestService_GetSystemNotiHasRead(t *testing.T) {
	service, err := InitDevService()
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		ctx context.Context
		obj *entity.SystemNoti
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetSystemNotiHasRead(tt.args.ctx, tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetSystemNotiHasRead() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Service.GetSystemNotiHasRead() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_GetSystemNotiContent(t *testing.T) {
	service, err := InitDevService()
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		ctx context.Context
		obj *entity.SystemNoti
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetSystemNotiContent(tt.args.ctx, tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetSystemNotiContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Service.GetSystemNotiContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_GetRepliedNotiHasRead(t *testing.T) {
	service, err := InitDevService()
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		ctx context.Context
		obj *entity.RepliedNoti
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetRepliedNotiHasRead(tt.args.ctx, tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetRepliedNotiHasRead() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Service.GetRepliedNotiHasRead() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_GetRepliedNotiThread(t *testing.T) {
	service, err := InitDevService()
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		ctx context.Context
		obj *entity.RepliedNoti
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.Thread
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetRepliedNotiThread(tt.args.ctx, tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetRepliedNotiThread() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.GetRepliedNotiThread() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_GetRepliedNotiRepliers(t *testing.T) {
	service, err := InitDevService()
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		ctx context.Context
		obj *entity.RepliedNoti
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetRepliedNotiRepliers(tt.args.ctx, tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetRepliedNotiRepliers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.GetRepliedNotiRepliers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_GetQuotedNotiHasRead(t *testing.T) {
	service, err := InitDevService()
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		ctx context.Context
		obj *entity.QuotedNoti
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetQuotedNotiHasRead(tt.args.ctx, tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetQuotedNotiHasRead() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Service.GetQuotedNotiHasRead() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_GetQuotedNotiThread(t *testing.T) {
	service, err := InitDevService()
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		ctx context.Context
		obj *entity.QuotedNoti
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.Thread
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetQuotedNotiThread(tt.args.ctx, tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetQuotedNotiThread() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.GetQuotedNotiThread() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_GetQuotedNotiQuotedPost(t *testing.T) {
	service, err := InitDevService()
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		ctx context.Context
		obj *entity.QuotedNoti
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.Post
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetQuotedNotiQuotedPost(tt.args.ctx, tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetQuotedNotiQuotedPost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.GetQuotedNotiQuotedPost() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_GetQuotedNotiPost(t *testing.T) {
	service, err := InitDevService()
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		ctx context.Context
		obj *entity.QuotedNoti
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.Post
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetQuotedNotiPost(tt.args.ctx, tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetQuotedNotiPost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.GetQuotedNotiPost() = %v, want %v", got, tt.want)
			}
		})
	}
}
