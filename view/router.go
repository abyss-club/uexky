package view

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gitlab.com/abyss.club/uexky/resolver"
)

// Router ...
func Router() http.Handler {
	handler := httprouter.New()

	resolver.Init() // TODO: put this to resolver
	handler.POST("/graphql/", withUexky(withAuthAndFlow(GraphQLHandle())))

	handler.GET("/auth/", withUexky(withAuthAndFlow(AuthHandle)))

	return handler
}
