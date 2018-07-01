package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/globalsign/mgo/bson"
	graphql "github.com/graph-gophers/graphql-go"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky/mgmt"
	"gitlab.com/abyss.club/uexky/model"
)

// NewRouter make router with all apis
func NewRouter() http.Handler {
	initRedis()
	schema := graphql.MustParseSchema(schema, &Resolver{})
	handler := httprouter.New()
	handler.POST("/graphql/", withAuth(graphqlHandle(schema)))
	handler.GET("/auth/", authHandle) // TODO: when user is already logged in
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

func withAuth(handle httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		tokenCookie, err := req.Cookie("token")
		if err != nil { // err must be ErrNoCookie,  non-login user, do noting
			handle(w, req, p)
			return
		}

		accountID, err := authToken(tokenCookie.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		if accountID != "" {
			log.Printf("Logged user %v", accountID)
			ctx := context.WithValue(req.Context(), model.ContextLoggedInAccount, accountID)
			req = req.WithContext(ctx)
		}
		handle(w, req, p)
	}
}

func authHandle(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	code := req.URL.Query().Get("code")
	if code == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("缺乏必要信息"))
		return
	}
	token, err := authCode(code)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("验证信息错误，或已失效"))
		return
	}

	// TODO: delete code in redis
	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		Domain:   mgmt.Config.Domain.WEB,
		Secure:   true,
		HttpOnly: true,
	}
	w.WriteHeader(http.StatusMovedPermanently)
	http.SetCookie(w, cookie)
	w.Header().Set("Location", mgmt.APIURLPrefix())
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

// Post ...
func (r *Resolver) Post(ctx context.Context, args struct{ ID string }) (*PostResolver, error) {
	post, err := model.FindPost(ctx, args.ID)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, nil
	}
	return &PostResolver{Post: post}, nil
}

// Uexky ...
func (r *Resolver) Uexky(ctx context.Context) (*UexkyResolver, error) {
	return &UexkyResolver{}, nil
}

// UexkyResolver ...
type UexkyResolver struct{}

// MainTags ...
func (ur *UexkyResolver) MainTags(ctx context.Context) ([]string, error) {
	return mgmt.Config.MainTags, nil
}

// Mutation:

// Auth ...
func (r *Resolver) Auth(ctx context.Context, args struct{ Email string }) (bool, error) {
	_, ok := ctx.Value(model.ContextLoggedInAccount).(bson.ObjectId)
	if ok {
		return false, nil
	}

	if !isValidateEmail(args.Email) {
		return false, errors.New("Invalid Email Address")
	}
	authURL := authEmail(args.Email)
	if err := sendAuthMail(authURL, args.Email); err != nil {
		return false, err
	}
	return true, nil
}

// AddName ...
func (r *Resolver) AddName(ctx context.Context, args struct{ Name string }) (*AccountResolver, error) {
	account, err := model.GetAccount(ctx)
	if err != nil {
		return nil, err
	}
	if err := account.AddName(ctx, args.Name); err != nil {
		return nil, err
	}
	return &AccountResolver{Account: account}, nil
}

// SyncTags ...
func (r *Resolver) SyncTags(
	ctx context.Context, args struct{ Tags []*string },
) (*AccountResolver, error) {
	account, err := model.GetAccount(ctx)
	if err != nil {
		return nil, err
	}
	tags := []string{}
	for _, t := range args.Tags {
		if t != nil {
			tags = append(tags, *t)
		}
	}
	if err := account.SyncTags(ctx, tags); err != nil {
		return nil, err
	}
	return &AccountResolver{Account: account}, nil
}

// PubThread ...
func (r *Resolver) PubThread(
	ctx context.Context,
	args struct{ Thread *model.ThreadInput },
) (
	*ThreadResolver, error,
) {
	thread, err := model.NewThread(ctx, args.Thread)
	if err != nil {
		return nil, err
	}
	return &ThreadResolver{Thread: thread}, nil
}

// PubPost ...
func (r *Resolver) PubPost(
	ctx context.Context,
	args struct{ Post *model.PostInput },
) (
	*PostResolver, error,
) {
	post, err := model.NewPost(ctx, args.Post)
	if err != nil {
		return nil, err
	}
	return &PostResolver{Post: post}, nil
}
