package metrics

import (
	"net/http"

	gqlgraphql "github.com/99designs/gqlgen/graphql"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Avoid adding a big number because getting metrics get a lock on gorm.
const defaultPrometheusGormRefreshMetricsSecond = 15

// Client Client metrics interface.
//go:generate mockgen -destination=./mocks/mock_Client.go -package=mocks github.com/oxyno-zeta/opa-center/pkg/opa-center/metrics Client
type Client interface {
	// Instrument web server.
	Instrument(serverName string) gin.HandlerFunc
	// Get prometheus handler for http expose.
	PrometheusHTTPHandler() http.Handler
	// Get database middleware.
	DatabaseMiddleware(connectionName string) gorm.Plugin
	// Get graphql middleware.
	GraphqlMiddleware() gqlgraphql.HandlerExtension
}

// NewMetricsClient will generate a new Client.
func NewMetricsClient() Client {
	ctx := &prometheusMetrics{
		gormPrometheus: map[string]gorm.Plugin{},
	}
	// Register
	ctx.register()

	return ctx
}
