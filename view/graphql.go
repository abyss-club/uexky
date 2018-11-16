package view

import (
	"encoding/json"
	"net/http"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/julienschmidt/httprouter"
	"gitlab.com/abyss.club/uexky/resolver"
)

// GraphQLHandle ...
func GraphQLHandle() httprouter.Handle {
	resolver.Init()
	schema := graphql.MustParseSchema(resolver.Schema, &resolver.Resolver{})
	return func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		params := graphqlParams{}
		if err := json.NewDecoder(req.Body).Decode(&params); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		response := schema.Exec(
			req.Context(), params.Query, params.OperationName, params.Variables,
		)
		responseJSON, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(responseJSON)
	}
}

type graphqlParams struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
}
