package mgmt

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/pkg/errors"
)

// Config for uexky
var Config struct {
	Mongo struct {
		URL string `json:"url"`
		DB  string `json:"db"`
	} `json:"mongo"`
	RedisURL string   `json:"redis_url"`
	MainTags []string `json:"main_tags"`
	Proto    string   `json:"proto"`
	Domain   struct {
		WEB string `json:"web"`
		API string `json:"api"`
	} `json:"domain"`
	Mail struct {
		Domain     string `json:"domain"`
		PrivateKey string `json:"private_key"`
		PublicKey  string `json:"public_key"`
	} `json:"mail"`
}

func setDefaultConfig() {
	Config.Mongo.URL = "localhost:27017"
	Config.Mongo.DB = "develop"
	Config.RedisURL = "redis://localhost:6379/0"
	Config.Proto = "https"
	Config.Domain.WEB = "abyss.club"
	Config.Domain.API = "api.abyss.club"
	Config.Mail.Domain = "mail.abyss.club"
}

// ReplaceConfigByEnv ...
func ReplaceConfigByEnv() {
	dbURL, found := os.LookupEnv("MONGO_URL")
	if found {
		Config.Mongo.URL = dbURL
	}

	redisURL, found := os.LookupEnv("REDIS_URL")
	if found {
		Config.RedisURL = redisURL
	}

	log.Printf("replaced config: %v", Config)
}

// WebURLPrefix ...
func WebURLPrefix() string {
	return fmt.Sprintf("%s://%s", Config.Proto, Config.Domain.WEB)
}

// APIURLPrefix ...
func APIURLPrefix() string {
	return fmt.Sprintf("%s://%s", Config.Proto, Config.Domain.API)
}

// LoadConfig from file
func LoadConfig(filename string) {
	setDefaultConfig()

	if filename != "" {
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Fatal(errors.Wrap(err, "Read config file error"))
		}
		if err := json.Unmarshal(b, &Config); err != nil {
			log.Fatal(errors.Wrap(err, "Read config error"))
		}
	}
	ReplaceConfigByEnv()

	log.Printf("load config: %v", Config)
}
