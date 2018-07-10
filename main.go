package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gitlab.com/abyss.club/uexky/api"
	"gitlab.com/abyss.club/uexky/mgmt"
	"gitlab.com/abyss.club/uexky/model"
	"gitlab.com/abyss.club/uexky/resolver"
)

var (
	configFile string
	serve      string
)

func init() {
	flag.StringVar(&configFile, "c", "", "config file")
	flag.StringVar(&serve, "s", ":5000", "server address")
}

func loadConfig() {
	flag.Parse()
	if configFile == "" {
		log.Fatal("Must specified config file")
	}
	mgmt.LoadConfig(configFile)
}

func newRouter() http.Handler {
	mongo := api.ConnectMongodb()
	resolver.Init()
	handler := httprouter.New()
	handler.POST(
		"/graphql/",
		api.WithRedis(api.WithMongo(api.WithAuth(resolver.GraphQLHandle()), mongo)),
	)
	handler.GET("/auth/", api.WithMongo(api.AuthHandle))
	return handler
}

func main() {
	loadConfig()

	if err := model.Init(); err != nil {
		log.Fatal(err)
	}

	router := newRouter()
	log.Print("start to serve")
	log.Fatal(http.ListenAndServe(serve, router))
}
