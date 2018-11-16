package view

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Router ...
func Router() http.Handler {
	handler := httprouter.New()
	handler.POST("/graphql/", withUexky(withAuthAndFlow(GraphQLHandle())))
	handler.GET("/auth/", withUexky(withAuthAndFlow(AuthHandle)))
	return handler
}
