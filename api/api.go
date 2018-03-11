package api

import (
	"io/ioutil"
	"log"
	"net/http"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/julienschmidt/httprouter"
)

// Resolver for graphql
type Resolver struct {
}

// NewRouter make router with all apis
func NewRouter(schemaFile string) http.Handler {
	b, err := ioutil.ReadFile(schemaFile)
	if err != nil {
		log.Fatal(err)
	}
	schema := graphql.MustParseSchema(string(b), &Resolver{})

	handler := httprouter.New()
	handler.Handler("POST", "/graphql/", &relay.Handler{Schema: schema})
	return handler
}
