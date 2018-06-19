package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/nanozuki/uexky/api"
	"github.com/nanozuki/uexky/model"
	"github.com/pkg/errors"
)

var (
	configFile string
	serve      string
	config     *Config
)

func init() {
	flag.StringVar(&configFile, "c", "", "config file")
	flag.StringVar(&serve, "s", ":5000", "server address")
}

// Config for whole project, saved by json
type Config struct {
	Mongo struct {
		URI string `json:"mongo_uri"`
		DB  string `json:"db"`
	} `json:"mongo"`
	MainTags []string `json:"main_tags"`
}

func loadConfig() {
	flag.Parse()
	if configFile == "" {
		log.Fatal("Must specified config file")
	}
	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Read config error"))
	}
	c := Config{}
	if err := json.Unmarshal(b, &c); err != nil {
		log.Fatal(errors.Wrap(err, "Read config error"))
	}
	config = &c
}

func main() {
	loadConfig()

	if err := model.Init(config.Mongo.URI, config.Mongo.DB, config.MainTags); err != nil {
		log.Fatal(err)
	}

	router := api.NewRouter()
	log.Print("start to serve")
	log.Fatal(http.ListenAndServe(serve, router))
}
