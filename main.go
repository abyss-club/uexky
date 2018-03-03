package main

import (
	"log"
	"net/http"

	"github.com/CrowsT/uexky/api"
	"github.com/CrowsT/uexky/model"
)

func main() {
	if err := model.Dial("localhost"); err != nil {
		log.Fatal(err)
	}
	router := api.NewRouter()
	log.Fatal(http.ListenAndServe(":5000", router))
}
