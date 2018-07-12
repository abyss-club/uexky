package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky/mgmt"
)

var redisPool = redis.Pool{
	Dial: func() (redis.Conn, error) {
		c, err := redis.DialURL(mgmt.Config.RedisURI)
		if err != nil {
			log.Fatal(errors.Wrap(err, "Connect to redis"))
		}
		return c, nil
	},
	TestOnBorrow: func(c redis.Conn, t time.Time) error {
		if time.Since(t) < time.Minute {
			return nil
		}
		_, err := c.Do("PING")
		return err
	},
	MaxIdle:     128,
	MaxActive:   1024,
	Wait:        true,
	IdleTimeout: time.Second * 60,
}

// WithRedis ...
func WithRedis(handle httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		conn := redisPool.Get()
		defer conn.Close()
		req = req.WithContext(context.WithValue(
			req.Context(), ContextKeyRedis, conn))
		handle(w, req, p)
	}
}

// GetRedis from context
func GetRedis(ctx context.Context) redis.Conn {
	c, ok := ctx.Value(ContextKeyMongo).(redis.Conn)
	if !ok {
		log.Fatal("Can't find mongodb in context")
	}
	return c
}
