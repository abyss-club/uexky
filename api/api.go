package api

// ContextKey ...
type ContextKey string

// ContextKeys
const (
	ContextKeyToken        = ContextKey("token")
	ContextKeyLoggedInUser = ContextKey("loggedIn")
	ContextKeyMongo        = ContextKey("mongo")
)
