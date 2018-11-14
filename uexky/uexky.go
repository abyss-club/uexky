package uexky

import (
	"context"
	"log"

	"github.com/globalsign/mgo/bson"
	"github.com/gomodule/redigo/redis"
)

// Pool store global data for Frame
type Pool struct {
	mongoPool *Mongo
	redisPool *redis.Pool
}

// InitPool ...
func InitPool() *Pool {
	return &Pool{
		mongoPool: ConnectMongodb(),
		redisPool: NewRedisPool(),
	}
}

// ContextKey ...
type ContextKey string

// ContextKeys
const (
	ContextKeyUexky = ContextKey("frame")
)

// Push an uexky object to context
func (p *Pool) Push(
	ctx context.Context, auth Auth, flow *Flow,
) (context.Context, func()) {
	frame := &Uexky{
		Mongo: p.mongoPool.Copy(),
		Redis: p.redisPool.Get(),
		Auth:  auth,
		Flow:  flow,
	}
	return context.WithValue(ctx, ContextKeyUexky, frame), frame.Close
}

// Pop an uexky from context, if don't find uexky, will panic
func Pop(ctx context.Context) *Uexky {
	u, ok := ctx.Value(ContextKeyUexky).(*Uexky)
	if !ok {
		log.Fatal("Can't find frame")
	}
	return u
}

// Uexky is context for http request
type Uexky struct {
	Mongo *Mongo
	Redis redis.Conn
	Auth  Auth
	Flow  *Flow
}

// Close ...
func (u *Uexky) Close() {
	u.Mongo.Close()
	u.Redis.Close()
	u.Auth = nil
	u.Flow = nil
}

// Auth ...
type Auth interface {
	IsSignedIn() bool
	RequireSignedIn() error
	Email() string
	ID() bson.ObjectId
	CheckPriority(action string) bool
}
