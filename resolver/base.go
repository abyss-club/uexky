package resolver

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/julienschmidt/httprouter"
	"gitlab.com/abyss.club/uexky/model"
)

// handle:

type graphqlParams struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
}

// GraphQLHandle ...
func GraphQLHandle() httprouter.Handle {
	schema := graphql.MustParseSchema(Schema, &Resolver{})
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

// Resolver for graphql
type Resolver struct{}

// types:

// SliceInfoResolver ...
type SliceInfoResolver struct {
	SliceInfo *model.SliceInfo
}

// FirstCursor ...
func (si *SliceInfoResolver) FirstCursor(ctx context.Context) (string, error) {
	return si.SliceInfo.FirstCursor, nil
}

// LastCursor ...
func (si *SliceInfoResolver) LastCursor(ctx context.Context) (string, error) {
	return si.SliceInfo.LastCursor, nil
}

// SliceQuery for api, different from model.SliceQuery
type SliceQuery struct {
	Before *string
	After  *string
	Limit  int32
}

// Parse to model.SliceQuery
func (sq *SliceQuery) Parse(reverse bool) (*model.SliceQuery, error) {
	if (sq.Before == nil && sq.After == nil) ||
		(sq.Before != nil && sq.After != nil) {
		return nil, errors.New("Invalid query")
	}
	if sq.Before != nil {
		return &model.SliceQuery{
			Limit:  int(sq.Limit),
			Desc:   !reverse,
			Cursor: *(sq.Before),
		}, nil
	}
	return &model.SliceQuery{
		Limit:  int(sq.Limit),
		Desc:   reverse,
		Cursor: *(sq.After),
	}, nil
}
