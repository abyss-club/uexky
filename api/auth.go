package api

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky/model"
	"gitlab.com/abyss.club/uexky/uuid64"
)

var redisConn redis.Conn

func init() {
	c, err := redis.DialURL("redis://localhost:6379/0")
	if err != nil {
		log.Fatal(errors.Wrap(err, "Connect to redis"))
	}
	redisConn = c
}

const (
	apiHostName  = "https://api.abyss.club"
	siteHostName = "https://abyss.club"
)

// 36 charactors Base64 token
var codeGenerator = uuid64.Generator{Sections: []uuid64.Section{
	&uuid64.RandomSection{Length: 10},
	&uuid64.CounterSection{Length: 4, Unit: time.Millisecond},
	&uuid64.TimestampSection{Length: 7, Unit: time.Millisecond},
	&uuid64.RandomSection{Length: 15},
}}

// 24 charactors Base64 token
var tokenGenerator = uuid64.Generator{Sections: []uuid64.Section{
	&uuid64.RandomSection{Length: 10},
	&uuid64.CounterSection{Length: 2, Unit: time.Millisecond},
	&uuid64.TimestampSection{Length: 7, Unit: time.Millisecond},
	&uuid64.RandomSection{Length: 5},
}}

func authEmail(email string) string {
	code, err := codeGenerator.New()
	if err != nil {
		log.Fatal(err)
	}
	if _, err := redisConn.Do("SET", code, email, "EX 600"); err != nil {
		log.Fatal(errors.Wrap(err, "set code to redis"))
	}
	return fmt.Sprintf("%s/auth/code?=%s", apiHostName, code)
}

func authCode(code string) (string, error) {
	email, err := redis.String(redisConn.Do("GET", code))
	if err == redis.ErrNil {
		return "", errors.New("Invalid code")
	} else if err != nil {
		return "", errors.Wrap(err, "Get code from redis")
	}
	account, err := model.FindAccountByEmail(context.Background(), email)
	if err != nil {
		return "", errors.Wrap(err, "find account")
	}
	token, err := tokenGenerator.New()
	if err != nil {
		return "", errors.Wrap(err, "gen token")
	}
	if _, err := redisConn.Do("SET", token, account.ID.Hex(), "EX 86400"); err != nil {
		return "", errors.Wrap(err, "set code to redis")
	}
	return token, nil
}

func sendAuthMail(code string) error {
	return nil // TODO:
}

func authToken(token string) (string, error) {
	idStr, err := redis.String(redisConn.Do("GET", token))
	if err == redis.ErrNil {
		return "", errors.New("Invalid Token")
	} else if err != nil {
		return "", errors.Wrap(err, "Get token from redis")
	}
	return bson.ObjectIdHex(idStr), nil
}
