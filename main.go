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

func parseFlag() {
	flag.Parse()
	if configFile == "" {
		log.Fatal("Must specified config file")
	}
}

// Config for whole project, saved by json
type Config struct {
	APISchemaFile string `json:"api_schema"`
}

func readConfig() {
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
	parseFlag()
	readConfig()

	if err := model.Dial("localhost"); err != nil {
		log.Fatal(err)
	}

	router := api.NewRouter(config.APISchemaFile)
	log.Fatal(http.ListenAndServe(serve, router))
}
