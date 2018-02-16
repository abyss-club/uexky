package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	res := map[string]string{"hello": "world"}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func main() {
	router := httprouter.New()
	router.GET("/", index)
	log.Fatal(http.ListenAndServe(":5000", router))
}
