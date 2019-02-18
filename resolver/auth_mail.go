package resolver

import (
	"fmt"
	"log"
	"strings"
	"time"

	mailgun "github.com/mailgun/mailgun-go"
	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky-go/config"
	"gitlab.com/abyss.club/uexky-go/uexky"
	"gitlab.com/abyss.club/uexky-go/uuid64"
)

var mailClient mailgun.Mailgun

// Init ...
func Init() {
	mailClient = mailgun.NewMailgun(
		config.Config.Mail.Domain,
		config.Config.Mail.PrivateKey,
		config.Config.Mail.PublicKey,
	)
}

// 36 charactors Base64 token
var codeGenerator = uuid64.Generator{Sections: []uuid64.Section{
	&uuid64.RandomSection{Length: 10},
	&uuid64.CounterSection{Length: 4, Unit: time.Millisecond},
	&uuid64.TimestampSection{Length: 7, Unit: time.Millisecond},
	&uuid64.RandomSection{Length: 15},
}}

func isValidateEmail(mail string) bool {
	// TODO: use regular expression
	if strings.Index(mail, "@") != -1 {
		return true
	}
	return false
}

func authEmail(u *uexky.Uexky, email string) (string, error) {
	if !isValidateEmail(email) {
		return "", errors.New("Invalid Email Address")
	}
	code, err := codeGenerator.New()
	if err != nil {
		return "", err
	}
	// expired at 20 minutes.
	if _, err := u.Redis.Do("SET", code, email, "EX", 1200); err != nil {
		return "", errors.Wrap(err, "set code to redis")
	}
	return fmt.Sprintf("%s/auth/?code=%s", config.APIURLPrefix(), code), nil
}

func sendAuthMail(url, to string) error {
	msg := mailClient.NewMessage(
		fmt.Sprintf("auth@%s", config.Config.Mail.Domain),
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
