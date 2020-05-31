package repo

import (
	"context"
	"reflect"
	"testing"
	"time"

	red "github.com/go-redis/redis/v7"
	"gitlab.com/abyss.club/uexky/uexky/entity"
)

func TestUserRepo_SetGetDelCodeEmail(t *testing.T) {
	type fields struct {
		Redis *red.Client
		Forum *ForumRepo
	}
	type args struct {
		ctx   context.Context
		email string
		code  string
		ex    time.Duration
		sleep time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
		setErr bool
		getErr bool
	}{
		{
			name: "normal",
			fields: fields{
				Redis: getRedis(t),
			},
			args: args{
				ctx:   context.Background(),
				email: "a@example.com",
				code:  "123",
				ex:    time.Minute,
				sleep: 0,
			},
			want: "a@example.com",
		},
		{
			name: "expired",
			fields: fields{
				Redis: getRedis(t),
			},
			args: args{
				ctx:   context.Background(),
				email: "b@example.com",
				code:  "456",
				ex:    time.Millisecond,
				sleep: 2 * time.Millisecond,
			},
			getErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserRepo{
				Redis: tt.fields.Redis,
				Forum: tt.fields.Forum,
			}
			if err := u.SetCode(tt.args.ctx, tt.args.email, tt.args.code, tt.args.ex); (err != nil) != tt.setErr {
				t.Errorf("UserRepo.SetCode() error = %v, wantErr %v", err, tt.setErr)
			}
			time.Sleep(tt.args.sleep)
			got, err := u.GetCodeEmail(tt.args.ctx, tt.args.code)
			if (err != nil) != tt.getErr {
				t.Errorf("UserRepo.GetCodeEmail() error = %v, wantErr %v", err, tt.getErr)
				return
			}
			if got != tt.want {
				t.Errorf("UserRepo.GetCodeEmail() = %v, want %v", got, tt.want)
			}
			if err := u.DelCode(tt.args.ctx, tt.args.code); err != nil {
				t.Errorf("UserRepo.DelCodeEmail() err = %v", err)
			}
			if _, err := u.GetCodeEmail(tt.args.ctx, tt.args.code); err != red.Nil {
				t.Errorf("UserRepo get email after delCode err = %v", err)
			}
		})
	}
}

func TestUserRepo_SetGetTokenEmail(t *testing.T) {
	type fields struct {
		Redis *red.Client
		Forum *ForumRepo
	}
	type args struct {
		ctx   context.Context
		email string
		tok   string
		ex    time.Duration
		sleep time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
		setErr bool
		getErr bool
	}{
		{
			name: "normal",
			fields: fields{
				Redis: getRedis(t),
			},
			args: args{
				ctx:   context.Background(),
				email: "a@example.com",
				tok:   "123",
				ex:    time.Minute,
			},
			want: "a@example.com",
		},
		{
			name: "expired",
			fields: fields{
				Redis: getRedis(t),
			},
			args: args{
				ctx:   context.Background(),
				email: "b@example.com",
				tok:   "456",
				ex:    1 * time.Millisecond,
				sleep: 2 * time.Millisecond,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserRepo{
				Redis: tt.fields.Redis,
				Forum: tt.fields.Forum,
			}
			if err := u.SetToken(tt.args.ctx, tt.args.email, tt.args.tok, tt.args.ex); (err != nil) != tt.setErr {
				t.Errorf("UserRepo.SetToken() error = %v, wantErr %v", err, tt.setErr)
			}
			time.Sleep(tt.args.sleep)
			got, err := u.GetTokenEmail(tt.args.ctx, tt.args.tok)
			if (err != nil) != tt.getErr {
				t.Errorf("UserRepo.GetTokenEmail() error = %v, wantErr %v", err, tt.getErr)
				return
			}
			if got != tt.want {
				t.Errorf("UserRepo.GetTokenEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserRepo_toEntityUser(t *testing.T) {
	type args struct {
		user *User
	}
	now := time.Now()
	u := &UserRepo{}
	name := "amy"
	tests := []struct {
		name string
		args args
		want *entity.User
	}{
		{
			"normal",
			args{
				user: &User{
					ID:                  1,
					CreatedAt:           now,
					UpdatedAt:           now,
					Email:               "a@example.com",
					Role:                "mod",
					LastReadSystemNoti:  2,
					LastReadRepliedNoti: 2,
					LastReadQuotedNoti:  2,
					Tags:                []string{"a", "b"},
				},
			},
			&entity.User{
				Email: "a@example.com",
				Role:  entity.RoleMod,
				Tags:  []string{"a", "b"},

				Repo: u,
				ID:   1,
				LastReadNoti: entity.LastReadNoti{
					SystemNoti:  2,
					RepliedNoti: 2,
					QuotedNoti:  2,
				},
			},
		},
		{
			"normal, with empty role",
			args{
				user: &User{
					ID:                  1,
					CreatedAt:           now,
					UpdatedAt:           now,
					Email:               "a@example.com",
					Name:                &name,
					LastReadSystemNoti:  2,
					LastReadRepliedNoti: 2,
					LastReadQuotedNoti:  2,
					Tags:                []string{"a", "b"},
				},
			},
			&entity.User{
				Email: "a@example.com",
				Name:  &name,
				Role:  entity.RoleNormal,
				Tags:  []string{"a", "b"},

				Repo: u,
				ID:   1,
				LastReadNoti: entity.LastReadNoti{
					SystemNoti:  2,
					RepliedNoti: 2,
					QuotedNoti:  2,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserRepo{}
			if got := u.toEntityUser(tt.args.user); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserRepo.toEntityUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserRepo_GetOrInsertUser(t *testing.T) {
	repo := &UserRepo{
		Forum: &ForumRepo{
			mainTags: []string{"A", "B"},
		},
	}
	ctx := getNewDBCtx(t)
	tests := []struct {
		name    string
		email   string
		want    *entity.User
		wantErr bool
	}{
		{
			name:  "new user",
			email: "a@example.com",
			want: &entity.User{
				Email: "a@example.com",
				Role:  entity.RoleNormal,
				Tags:  []string{"A", "B"},
				Repo:  repo,
				ID:    1,
			},
		},
		{
			name:  "get user",
			email: "a@example.com",
			want: &entity.User{
				Email: "a@example.com",
				Role:  entity.RoleNormal,
				Tags:  []string{"A", "B"},
				Repo:  repo,
				ID:    1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetOrInsertUser(ctx, tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepo.GetOrInsertUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserRepo.GetOrInsertUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserRepo_UpdateUser(t *testing.T) {
	repo := &UserRepo{
		Forum: &ForumRepo{
			mainTags: []string{"A", "B"},
		},
	}
	ctx := getNewDBCtx(t)
	email := "a@example.com"
	user, err := repo.GetOrInsertUser(ctx, email)
	if err != nil {
		t.Errorf("create user error: %v", err)
	}
	newRole := entity.RoleMod
	newName1 := "amy"
	type args struct {
		id     int
		update *entity.UserUpdate
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.User
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				id: user.ID,
				update: &entity.UserUpdate{
					Role: &newRole,
					Tags: []string{"B", "C"},
				},
			},
			want: &entity.User{
				Email: "a@example.com",
				Role:  entity.RoleMod,
				Tags:  []string{"B", "C"},
				Repo:  repo,
				ID:    1,
			},
		},
		{
			name: "set name",
			args: args{
				id: user.ID,
				update: &entity.UserUpdate{
					Name: &newName1,
				},
			},
			want: &entity.User{
				Email: "a@example.com",
				Name:  &newName1,
				Role:  entity.RoleMod,
				Tags:  []string{"B", "C"},
				Repo:  repo,
				ID:    1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := repo.UpdateUser(ctx, tt.args.id, tt.args.update); (err != nil) != tt.wantErr {
				t.Errorf("UserRepo.UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
			got, err := repo.GetOrInsertUser(ctx, user.Email)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepo.GetOrInsertUser() error = %+v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserRepo.GetOrInsertUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
