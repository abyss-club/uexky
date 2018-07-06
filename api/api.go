package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/julienschmidt/httprouter"
	"gitlab.com/abyss.club/uexky/mgmt"
	"gitlab.com/abyss.club/uexky/model"
	"gitlab.com/abyss.club/uexky/resolver"
)

// NewRouter make router with all apis
func NewRouter() http.Handler {
	resolver.InitRedis()
	schema := graphql.MustParseSchema(resolver.Schema, &resolver.Resolver{})
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
		log.Printf("find token cookie %v", tokenCookie)
		if err != nil { // err must be ErrNoCookie,  non-login user, do noting
			handle(w, req, p)
			return
		}

		userID, err := resolver.AuthToken(tokenCookie.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		if userID != "" {
			log.Printf("Logged user %v", userID)
			ctx := context.WithValue(req.Context(), model.ContextLoggedInUser, userID)
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
	token, err := resolver.AuthCode(code)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("验证信息错误，或已失效。 %v", err)))
		return
	}

	resolver.RedisConn.Do("DEL", code) // delete after use
	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		Domain:   mgmt.Config.Domain.WEB,
		MaxAge:   86400,
		HttpOnly: true,
	}
	if mgmt.Config.Proto == "https" {
		cookie.Secure = true
	}
	http.SetCookie(w, cookie)
	w.Header().Set("Location", mgmt.WebURLPrefix())
	w.Header().Set("Cache-Control", "no-cache, no-store")
	w.WriteHeader(http.StatusFound)
}
