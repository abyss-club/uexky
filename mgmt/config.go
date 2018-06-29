package mgmt

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/pkg/errors"
)

// Config for uexky
var Config struct {
	Mongo struct {
		URI string `json:"mongo_uri"`
		DB  string `json:"db"`
	} `json:"mongo"`
	RedisURI string
	MainTags []string `json:"main_tags"`
	Proto    string
	Domain   struct {
		WEB string
		API string
	}
}

func setDefaultConfig() {
	Config.Mongo.URI = "localhost:27017"
	Config.Mongo.DB = "develop"
	Config.RedisURI = "redis://localhost:6379/0"
	Config.Proto = "https"
	Config.Domain.WEB = "abyss.club"
	Config.Domain.API = "api.abyss.club"
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
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Read config file error"))
	}
	if err := json.Unmarshal(b, &Config); err != nil {
		log.Fatal(errors.Wrap(err, "Read config error"))
	}
	log.Printf("load config: %v", Config)
}
