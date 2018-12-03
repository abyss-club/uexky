package model

import (
	"log"

	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky/uexky"
)

// AuthInfo impl uexky.Auth
type AuthInfo struct {
	uexky *uexky.Uexky
	email string

	userID bson.ObjectId
	user   *User
}

// NewUexkyAuth make a new AuthInfo and add to Uexky
func NewUexkyAuth(uexky *uexky.Uexky, email string) *AuthInfo {
	if email != "" {
		log.Printf("Logged user %s", email)
	}
	ai := &AuthInfo{uexky: uexky, email: email}
	uexky.Auth = ai
	return ai
}

// IsSignedIn ...
func (ai *AuthInfo) IsSignedIn() bool {
	return ai.email != ""
}

// RequireSignedIn ...
func (ai *AuthInfo) RequireSignedIn() error {
	if !ai.IsSignedIn() {
		return errors.New("unauthorized, require signed in")
	}
	return nil
}

// Email ...
func (ai *AuthInfo) Email() string {
	return ai.email
}

// CheckPriority ...
func (ai *AuthInfo) CheckPriority(action, target string) (bool, error) {
	user, err := ai.GetUser()
	if err != nil {
		return false, err
	}
	return user.CheckPriority(ai.uexky, action, target)
}

// GetUser who signed in
func (ai *AuthInfo) GetUser() (*User, error) {
	if err := ai.RequireSignedIn(); err != nil {
		return nil, errors.New("unauthorized, require signed in")
	}
	if ai.user != nil {
		return ai.user, nil
	}
	return GetUserByEmail(ai.uexky, ai.email)
}
