package uexky

import (
	"context"
	"log"

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

// NewUexky Make a new Uexky
func (p *Pool) NewUexky() *Uexky {
	log.Print("New Uexky!!!")
	u := &Uexky{
		Mongo: p.mongoPool.Copy(),
		Redis: NewRedis(p.redisPool),
	}
	return u
}

// Push an uexky object to context
func (p *Pool) Push(ctx context.Context) (context.Context, func()) {
	uexky := p.NewUexky()
	return context.WithValue(ctx, contextKeyUexky, uexky), uexky.Close
}

// Pop an uexky from context, if don't find uexky, will panic
func Pop(ctx context.Context) *Uexky {
	u, ok := ctx.Value(contextKeyUexky).(*Uexky)
	if !ok {
		log.Fatal("Can't find frame")
	}
	return u
}

// contextKey ...
type contextKey string

// contextKeys
const (
	contextKeyUexky = contextKey("frame")
)

// Uexky is context for http request
type Uexky struct {
	Mongo *Mongo
	Redis *Redis
	Auth  Auth
	Flow  Flow
}

// Close ...
func (u *Uexky) Close() {
	u.Mongo.Close()
}

// Auth ...
type Auth interface {
	IsSignedIn() bool
	RequireSignedIn() error
	Email() string
	CheckPriority(action, target string) (bool, error)
}
