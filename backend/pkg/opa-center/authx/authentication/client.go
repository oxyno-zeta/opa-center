package authentication

import (
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/config"
)

//go:generate mockgen -destination=./mocks/mock_Client.go -package=mocks github.com/oxyno-zeta/opa-center/pkg/opa-center/authx/authentication Client
type Client interface {
	// Middleware will redirect authentication to basic auth or OIDC depending on request path and resources declared.
	Middleware(unauthorizedPathRegexList []*regexp.Regexp) gin.HandlerFunc
	// OIDCEndpoints will set OpenID Connect endpoints for authentication and callback.
	OIDCEndpoints(router gin.IRouter) error
}

func NewService(cfgManager config.Manager) Client {
	return &service{
		cfgManager: cfgManager,
	}
}
