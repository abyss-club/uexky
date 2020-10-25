package auth

/*

func TestService_CtxWithUserByToken(t *testing.T) {
	service, ctx := initEnv(t)
	type args struct {
		ctx   context.Context
		email string
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.User
		wantErr bool
	}{
		{
			name: "signed in user",
			args: args{
				ctx:   ctx,
				email: "user1@example.com",
			},
			want: &entity.User{
				Email: algo.NullString("user1@example.com"),
				Role:  entity.RoleNormal,
				Repo:  service.User.Repo,
			},
		},
		{
			name: "guest user first login",
			args: args{
				ctx: ctx,
			},
			want: &entity.User{
				Role: entity.RoleGuest,
				Repo: service.User.Repo,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var signUpToken *entity.Token
			var signUpUser *entity.User
			var gotCtx context.Context
			var gotToken *entity.Token
			var err error
			if tt.args.email != "" {
				code, err := service.TrySignInByEmail(tt.args.ctx, tt.args.email, "")
				if err != nil {
					t.Fatal(errors.Wrap(err, "TrySignInByEmail"))
				}
				signUpToken, err = service.SignInByCode(ctx, string(code))
				if err != nil {
					t.Fatal(errors.Wrap(err, "SignInByCode"))
				}
			} else {
				var err error
				var signUpCtx context.Context
				signUpCtx, signUpToken, err = service.CtxWithUserByToken(ctx, "")
				if err != nil {
					t.Fatal(errors.Wrap(err, "CtxWithUserByToken"))
				}
				signUpUser, err = service.Profile(signUpCtx)
				if err != nil {
					t.Fatal(errors.Wrap(err, "Profile"))
				}
			}
			gotCtx, gotToken, err = service.CtxWithUserByToken(tt.args.ctx, signUpToken.Tok)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.CtxWithUserByToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotUser, err := service.Profile(gotCtx)
			if err != nil {
				t.Fatal(errors.Wrap(err, "Profile"))
			}
			tt.want.ID = gotUser.ID
			if !reflect.DeepEqual(gotUser, tt.want) {
				t.Errorf("want user %+v, bug got %+v", tt.want, gotUser)
			}
			if !reflect.DeepEqual(signUpToken, gotToken) {
				t.Errorf("token when signed up is %+v, bug got %+v", signUpToken, gotToken)
			}
			if tt.args.email == "" {
				if !reflect.DeepEqual(signUpUser, gotUser) {
					t.Errorf("guest user, 2nd time login = %+v, want equal to first time: %+v", gotUser, signUpUser)
				}
			}
		})
	}
}
*/
