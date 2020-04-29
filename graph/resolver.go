package graph

import "gitlab.com/abyss.club/uexky/service"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Service *service.Service
}

func NewResolver(s service.Service) Resolver {
	return Resolver{Service: &s}
}
