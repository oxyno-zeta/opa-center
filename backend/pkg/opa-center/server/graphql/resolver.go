package graphql

import "github.com/oxyno-zeta/opa-center/pkg/opa-center/business"

//go:generate gqlgen generate

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	BusiServices *business.Services
}
