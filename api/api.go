package api

// ContextKey ...
type ContextKey string

// ContextKeys
const (
	ContextKeyToken = ContextKey("token")
	ContextKeyEmail = ContextKey("loggedIn")
	ContextKeyUser  = ContextKey("user")
	ContextKeyMongo = ContextKey("mongo")
	ContextKeyRedis = ContextKey("redis")
)
