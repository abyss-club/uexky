package graph

import (
	"gitlab.com/abyss.club/uexky/auth"
	"gitlab.com/abyss.club/uexky/uexky"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Auth  *auth.Service
	Uexky *uexky.Service
}
