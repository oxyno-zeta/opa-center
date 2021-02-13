package tracing

import (
	"context"

	gqlgraphql "github.com/99designs/gqlgen/graphql"
	"github.com/gin-gonic/gin"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/config"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/log"
)

// Service Tracing service.
//go:generate mockgen -destination=./mocks/mock_Service.go -package=mocks github.com/oxyno-zeta/opa-center/pkg/opa-center/tracing Service
type Service interface {
	// Reload service
	Reload() error
	// Get opentracing tracer
	GetTracer() opentracing.Tracer
	// Http Gin HttpMiddleware to add trace per request
	HTTPMiddleware(getRequestID func(ctx context.Context) string) gin.HandlerFunc
	// Graphql Middleware
	GraphqlMiddleware() gqlgraphql.HandlerExtension
}

// Trace structure.
//go:generate mockgen -destination=./mocks/mock_Trace.go -package=mocks github.com/oxyno-zeta/opa-center/pkg/opa-center/tracing Trace
type Trace interface {
	// Add tag to trace
	SetTag(key string, value interface{})
	// Get a child trace
	GetChildTrace(operationName string) Trace
	// End the trace
	Finish()
	// Get the trace ID
	GetTraceID() string
}

func New(cfgManager config.Manager, logger log.Logger) (Service, error) {
	return newService(cfgManager, logger)
}
