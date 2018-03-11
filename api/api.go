package api

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/CrowsT/uexky/model"
	graphql "github.com/graph-gophers/graphql-go"
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
	handler.POST("/graphql/", withToken(graphqlHandle(schema)))
	return handler
}

type graphqlParams struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
}

func graphqlHandle(schema *graphql.Schema) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		params := graphqlParams{}
		if err := json.NewDecoder(req.Body).Decode(&params); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		response := schema.Exec(req.Context(), params.Query, params.OperationName, params.Variables)
		responseJSON, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(responseJSON)
	}
}

func withToken(handle httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		token := req.Header.Get("Access-Token")
		if token != "" && len(token) != 24 {
			http.Error(w, "Invalid Token Format", http.StatusForbidden)
			return
		}
		ctx := context.WithValue(req.Context(), model.CtxTokenKey{}, token)
		req.WithContext(ctx)
		handle(w, req, p)
	}
}
