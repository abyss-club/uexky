package model

import (
	"gitlab.com/abyss.club/uexky/config"
	"gitlab.com/abyss.club/uexky/mongo"
)

// collection name for model
const (
	colleUser         = "users"
	colleAID          = "anonymous_ids"
	collePost         = "posts"
	colleThread       = "threads"
	colleTag          = "tags"
	colleNotification = "notification"
)

// Models
var (
	UserModel   = &mongo.Model{C: "users"}
	AIDModel    = &mongo.Model{C: "anonymous_ids"}
	PostModel   = &mongo.Model{C: "posts"}
	ThreadModel = &mongo.Model{C: "threads"}
	TagModel    = &mongo.Model{C: "tags"}
	NotiModel   = &mongo.Model{C: "noti"}
)

// Register all models
func Register(client *mongo.Client) error {
	db := config.Config.Mongo.DB
	UserModel = client.RegisterModel(db, "users", nil)
	return nil
}
