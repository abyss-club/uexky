package model

import (
	"testing"

	"gitlab.com/abyss.club/uexky/uexky"
)

func TestAuthInfo(t *testing.T) {
	t.Log("not signed in")
	{
		u := uexkyPool.NewUexky()
		defer u.Close()
		ai := NewUexkyAuth(u, "")
		uexky.NewUexkyFlow(u, "1", "2")

		if ai.IsSignedIn() != false {
			t.Error("Not signed in, but get true")
		}
		if ai.RequireSignedIn() == nil {
			t.Error("RequireSignedIn should be error")
		}
		if ai.Email() != "" {
			t.Errorf("Email should be empty, but get %s", ai.Email())
		}
		// if ai.CheckPriority("") == true {
		// 	t.Error("CheckPriority must be false", ai.CheckPriority(""))
		// }
		if _, err := ai.GetUser(); err == nil {
			t.Error("GetUser should be error")
		}
	}
	{
		u := uexkyPool.NewUexky()
		defer u.Close()
		ai := NewUexkyAuth(u, mockUsers[0].Email)
		uexky.NewUexkyFlow(u, "1", "2")

		if ai.IsSignedIn() != true {
			t.Error("Signed in, but get false")
		}
		if err := ai.RequireSignedIn(); err != nil {
			t.Errorf("RequireSignedIn() error = %v, want nil", err)
		}
		if email := ai.Email(); email != mockUsers[0].Email {
			t.Errorf("Email() = %s, want %s", email, mockUsers[0].Email)
		}
		// if ai.CheckPriority("") == true {
		// 	t.Error("CheckPriority must be false", ai.CheckPriority(""))
		// }
		if user, err := ai.GetUser(); err != nil {
			t.Errorf("GetUser shouldn't be error, but get %s", err)
		} else if user.ID != mockUsers[0].ID {
			t.Errorf("GetUser = %v, want %v", user, mockUsers[0])
		}
	}
}
