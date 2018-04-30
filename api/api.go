package api

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/julienschmidt/httprouter"
	"github.com/nanozuki/uexky/model"
)

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

// Resolver for graphql
type Resolver struct{}

// Query:

// Account resolve query 'account'
func (r *Resolver) Account(ctx context.Context) (*AccountResolver, error) {
	account, err := model.GetAccount(ctx)
	return &AccountResolver{account}, err
}

// ThreadSlice ...
func (r *Resolver) ThreadSlice(ctx context.Context, args struct {
	Limit int
	Tags  *[]string
	After *string
}) (
	*ThreadSliceResolver, error,
) {
	after := ""
	if args.After != nil {
		after = *args.After
	}
	tags := []string{}
	if args.Tags != nil {
		tags = *args.Tags
	}

	sq := &model.SliceQuery{Limit: args.Limit, After: after}
	threads, sliceInfo, err := model.GetThreadsByTags(ctx, tags, sq)
	if err != nil {
		return nil, err
	}

	var trs []*ThreadResolver
	for _, t := range threads {
		trs = append(trs, &ThreadResolver{Thread: t})
	}
	sir := &SliceInfoResolver{SliceInfo: sliceInfo}
	return &ThreadSliceResolver{threads: trs, sliceInfo: sir}, nil
}

// Thread ...
func (r *Resolver) Thread(ctx context.Context, args struct{ ID string }) (*ThreadResolver, error) {
	th, err := model.FindThread(ctx, args.ID)
	if err != nil {
		return nil, err
	}
	if th == nil {
		return nil, nil
	}
	return &ThreadResolver{Thread: th}, nil
}

// Mutation:

// AddAccount resolve mutation 'addAccount'
func (r *Resolver) AddAccount(ctx context.Context) (*AccountResolver, error) {
	account, err := model.NewAccount(ctx)
	return &AccountResolver{account}, err
}

// AddName ...
func (r *Resolver) AddName(ctx context.Context, args struct{ Name string }) (*AccountResolver, error) {
	account, err := model.GetAccount(ctx)
	if err != nil {
		return nil, nil
	}
	if err := account.AddName(ctx, args.Name); err != nil {
		return nil, err
	}
	return &AccountResolver{Account: account}, nil
}
