package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gitlab.com/abyss.club/uexky/mgmt"
	"gitlab.com/abyss.club/uexky/mw"
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
	mongo := mw.ConnectMongodb()
	resolver.Init()
	handler := httprouter.New()
	handler.POST(
		"/graphql/",
		mw.WithRedis(mw.WithMongo(mw.WithAuth(resolver.GraphQLHandle()), mongo)),
	)
	handler.GET("/auth/", mw.WithMongo(mw.AuthHandle, mongo))
	return handler
}

func main() {
	loadConfig()
	router := newRouter()
	log.Print("start to serve")
	log.Fatal(http.ListenAndServe(serve, router))
}
