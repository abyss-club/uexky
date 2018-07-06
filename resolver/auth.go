package resolver

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/gomodule/redigo/redis"
	mailgun "github.com/mailgun/mailgun-go"
	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky/mgmt"
	"gitlab.com/abyss.club/uexky/model"
	"gitlab.com/abyss.club/uexky/uuid64"
)

// RedisConn ...
var RedisConn redis.Conn
var mailClient mailgun.Mailgun

// InitRedis ...
func InitRedis() {
	c, err := redis.DialURL(mgmt.Config.RedisURI)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Connect to redis"))
	}
	RedisConn = c
	mailClient = mailgun.NewMailgun(
		mgmt.Config.Mail.Domain, mgmt.Config.Mail.PrivateKey,
		mgmt.Config.Mail.PublicKey,
	)
}

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

func authEmail(email string) (string, error) {
	code, err := codeGenerator.New()
	if err != nil {
		return "", err
	}
	if _, err := RedisConn.Do("SET", code, email, "EX", 3600); err != nil {
		return "", errors.Wrap(err, "set code to redis")
	}
	return fmt.Sprintf("%s/auth/?code=%s", mgmt.APIURLPrefix(), code), nil
}

// AuthCode ...
func AuthCode(code string) (string, error) {
	email, err := redis.String(RedisConn.Do("GET", code))
	if err == redis.ErrNil {
		return "", errors.New("Invalid code")
	} else if err != nil {
		return "", errors.Wrap(err, "Get code from redis")
	}
	user, err := model.GetUserByEmail(context.Background(), email)
	if err != nil {
		return "", errors.Wrap(err, "find user")
	}
	token, err := tokenGenerator.New()
	if err != nil {
		return "", errors.Wrap(err, "gen token")
	}
	if _, err := RedisConn.Do("SET", token, user.ID.Hex(), "EX", 600); err != nil {
		return "", errors.Wrap(err, "set token to redis")
	}
	return token, nil
}

func sendAuthMail(url, to string) error {
	msg := mailClient.NewMessage(
		fmt.Sprintf("auth@%s", mgmt.Config.Mail.Domain),
		"点击登入 Abyss!",
		fmt.Sprintf("点击此链接进入 Abyss：%s", url),
		to,
	)
	msg.SetHtml(fmt.Sprintf(`
<html>
    <head>
        <meta charset="utf-8">
        <title>点击登入 Abyss!</title>
    </head>
    <body>
        <p>点击 <a href="%s">此链接</a> 进入 Abyss</p>
    </body>
</html>`, url))
	res, id, err := mailClient.Send(msg)
	if err != nil {
		return errors.Wrap(err, "Send Auth Email")
	}
	log.Printf("Send Email to %s, id = %s, res = %s", to, id, res)
	return nil
}

// AuthToken ...
func AuthToken(token string) (bson.ObjectId, error) {
	idStr, err := redis.String(RedisConn.Do("GET", token))
	if err == redis.ErrNil {
		return "", nil
	} else if err != nil {
		return "", errors.Wrap(err, "Get token from redis")
	}
	if !bson.IsObjectIdHex(idStr) {
		return "", nil // Can't find valid user.
	}
	return bson.ObjectIdHex(idStr), nil
}

func isValidateEmail(mail string) bool {
	// TODO: use regular expression
	if strings.Index(mail, "@") != -1 {
		return true
	}
	return false
}
