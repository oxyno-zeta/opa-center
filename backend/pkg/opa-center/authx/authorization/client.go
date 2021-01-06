package authorization

import (
	"context"

	"github.com/oxyno-zeta/opa-center/pkg/opa-center/config"
)

//go:generate mockgen -destination=./mocks/mock_Service.go -package=mocks github.com/oxyno-zeta/opa-center/pkg/opa-center/authx/authorization Service
type Service interface {
	// Check if it is authorized
	IsAuthorized(ctx context.Context, action, resource string) (bool, error)
	// Check authorized and fail if not authorized
	CheckAuthorized(ctx context.Context, action, resource string) error
}

func NewService(cfgManager config.Manager) Service {
	return &service{
		cfgManager: cfgManager,
	}
}
