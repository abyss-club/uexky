package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/CrowsT/uexky/api"
	"github.com/CrowsT/uexky/model"
	"github.com/pkg/errors"
)

var (
	configFile string
	config     *Config
)

func init() {
	flag.StringVar(&configFile, "c", "", "config file")
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

	router := api.NewRouter()
	log.Fatal(http.ListenAndServe(":5000", router))
}
