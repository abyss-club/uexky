package auth

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gitlab.com/abyss.club/uexky/lib/config"
	"gitlab.com/abyss.club/uexky/lib/errors"
	"gitlab.com/abyss.club/uexky/mocks"
)

func TestMain(m *testing.M) {
	if err := config.Load(""); err != nil {
		log.Fatalf("load config: %v", err)
	}
	fmt.Printf("run test in config: %#v\n", config.Get())
	os.Exit(m.Run())
}

func TestService_EmailSignFlow(t *testing.T) {
	s, err := InitMockAuthService()
	if err != nil {
		t.Fatal(err)
	}

	// args
	ctx := context.Background()
	userEmail := "a@example.com"
	redirectTo := "/thread=1"

	// TrySignInByEmail

	gotCode, err := s.TrySignInByEmail(ctx, userEmail, redirectTo)
	if err != nil {
		t.Fatalf("Service.TrySignInByEmail() error = %v, wantErr %v", err, false)
		return
	}
	mail := s.Mail.(*mocks.MailAdapter).LastMail
	if mail.To != userEmail {
		t.Fatalf("Service.TrySignInByEmail(), mail send to = %v, want %v", mail.To, userEmail)
	}
	if !strings.Contains(mail.Text, string(gotCode)) {
		t.Fatalf("Service.TrySignInByEmail(), mail text should contains code")
	}
	if !strings.Contains(mail.Text, redirectTo) {
		t.Fatalf("Service.TrySignInByEmail(), mail text should contains redirectTo")
	}

	// SignInByCode

	token, err := s.SignInByCode(ctx, gotCode)
	if err != nil {
		t.Fatalf("Service.SignInByCode() err = %+v", err)
	}
	if diff := cmp.Diff(token.User, UserInfo{Email: userEmail, IsGuest: false}); diff != "" {
		t.Fatalf("Service.SignInByCode(), token.User mismatch: %s", diff)
	}
	_, err = s.SignInByCode(ctx, gotCode)
	if err == nil && !errors.Is(err, errors.NotFound) {
		t.Fatalf("Service.SignInByCode() again, should get NotFound err, but got = %+v", err)
	}

	// GetToken

	gotToken, err := s.GetToken(ctx, token.Tok)
	if err != nil {
		t.Fatalf("Service.GetToken() err = %+v", err)
	}
	if diff := cmp.Diff(token, gotToken); diff != "" {
		t.Fatalf("Service.GetToken(), mismatch: %s", diff)
	}
}

func TestService_GuestUserFlow(t *testing.T) {
	s, err := InitMockAuthService()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	// SignInGuestUser

	token, err := s.SignInGuestUser(ctx)
	if err != nil {
		t.Fatalf("Service.SignInGuestUser() err = %+v", err)
	}
	if !token.User.IsGuest || token.User.Email != "" || token.User.UserID == 0 {
		t.Fatalf("Service.SignInGuestUser() must be guest user, but get %+v", token.User)
	}

	// GetToken

	gotToken, err := s.GetToken(ctx, token.Tok)
	if err != nil {
		t.Fatalf("Service.GetToken() err = %+v", err)
	}
	if diff := cmp.Diff(token, gotToken); diff != "" {
		t.Fatalf("Service.GetToken(), mismatch: %s", diff)
	}
}
