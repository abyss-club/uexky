package uexky

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky/config"
)

// NewRedisPool ...
func NewRedisPool() *redis.Pool {
	return &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(config.Config.RedisURL)
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
		IdleTimeout: time.Second * 240,
	}
}

// Redis ...
type Redis struct {
	pool *redis.Pool
}

// NewRedis ...
func NewRedis(pool *redis.Pool) *Redis {
	return &Redis{pool: pool}
}

// Do ...
func (r *Redis) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	conn := r.pool.Get()
	defer conn.Close()
	return conn.Do(commandName, args...)
}

// SetCache ...
func SetCache(u *Uexky, key string, value interface{}, expire int) error {
	var err error
	b, err := json.Marshal(value)
	if err != nil {
		log.Fatal(errors.Wrap(err, "marshall cache"))
	}
	if expire > 0 {
		_, err = u.Redis.Do("SET", key, string(b), "EX", expire)
	} else {
		_, err = u.Redis.Do("SET", key, string(b))
	}
	if err != nil {
		return errors.Wrap(err, "set cache")
	}
	return nil
}

// GetCache ...
func GetCache(u *Uexky, key string, value interface{}) (bool, error) {
	vs, err := redis.String(u.Redis.Do("GET", key))
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
