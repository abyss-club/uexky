package mw

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky/mgmt"
)

// RedisPool ...
var RedisPool = redis.Pool{
	Dial: func() (redis.Conn, error) {
		c, err := redis.DialURL(mgmt.Config.RedisURL)
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
		conn := RedisPool.Get()
		defer conn.Close()
		req = req.WithContext(context.WithValue(
			req.Context(), ContextKeyRedis, conn))
		handle(w, req, p)
	}
}

// GetRedis from context
func GetRedis(ctx context.Context) redis.Conn {
	c, ok := ctx.Value(ContextKeyRedis).(redis.Conn)
	if !ok {
		log.Fatal("Can't find redis in context")
	}
	return c
}

// SetCache ...
func SetCache(ctx context.Context, key string, value interface{}, expire int) error {
	var err error
	b, err := json.Marshal(value)
	if err != nil {
		log.Fatal(errors.Wrap(err, "marshall cache"))
	}
	if expire > 0 {
		_, err = GetRedis(ctx).Do("SET", key, string(b), "EX", expire)
	} else {
		_, err = GetRedis(ctx).Do("SET", key, string(b))
	}
	if err != nil {
		return errors.Wrap(err, "set cache")
	}
	return nil
}

// GetCache ...
func GetCache(ctx context.Context, key string, value interface{}) (bool, error) {
	vs, err := redis.String(GetRedis(ctx).Do("GET", key))
	if err == redis.ErrNil {
		return false, nil
	} else if err != nil {
		return false, errors.Wrap(err, "get cache")
	}
	if err := json.Unmarshal([]byte(vs), value); err != nil {
		log.Fatal(errors.Wrap(err, "unmarshall cache"))
	}
	return true, nil
}
