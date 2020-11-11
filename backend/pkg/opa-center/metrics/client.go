package metrics

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Client Client metrics interface
//go:generate mockgen -destination=./mocks/mock_Client.go -package=mocks github.com/oxyno-zeta/opa-center/pkg/opa-center/metrics Client
type Client interface {
	// Instrument web server
	Instrument(serverName string) gin.HandlerFunc
	// Get prometheus handler for http expose
	GetPrometheusHTTPHandler() http.Handler
}

// NewMetricsClient will generate a new Client.
func NewMetricsClient() Client {
	ctx := &prometheusMetrics{}
	ctx.register()

	return ctx
}
