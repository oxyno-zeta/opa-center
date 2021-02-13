package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/gorm"

	gqlprometheus "github.com/99designs/gqlgen-contrib/prometheus"
	gqlgraphql "github.com/99designs/gqlgen/graphql"
	gormprometheus "gorm.io/plugin/prometheus"
)

type prometheusMetrics struct {
	reqCnt         *prometheus.CounterVec
	resSz          *prometheus.SummaryVec
	reqDur         *prometheus.SummaryVec
	reqSz          *prometheus.SummaryVec
	up             prometheus.Gauge
	gormPrometheus map[string]gorm.Plugin
}

func (ctx *prometheusMetrics) GraphqlMiddleware() gqlgraphql.HandlerExtension {
	return gqlprometheus.Tracer{}
}

func (ctx *prometheusMetrics) PrometheusHTTPHandler() http.Handler {
	return promhttp.Handler()
}

// The gorm prometheus plugin cannot be instantiated twice because there is a loop inside that cannot be modified or stopped.
// This loop get all data from database and the loop cannot be modified in terms of the duration.
// Labels and all other options cannot be modified.
func (ctx *prometheusMetrics) DatabaseMiddleware(connectionName string) gorm.Plugin {
	// Check if gorm prometheus doesn't already exist
	if ctx.gormPrometheus[connectionName] == nil {
		// Create middleware
		md := gormprometheus.New(gormprometheus.Config{
			RefreshInterval: defaultPrometheusGormRefreshMetricsSecond, // refresh metrics interval (default 15 seconds)
		})
		// Apply labels
		md.Labels = map[string]string{
			"connection_name": connectionName,
		}
		// Save it
		ctx.gormPrometheus[connectionName] = md
	}

	return ctx.gormPrometheus[connectionName]
}

// Instrument will instrument gin routes.
func (ctx *prometheusMetrics) Instrument(serverName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		reqSz := computeApproximateRequestSize(c.Request)

		c.Next()

		status := strconv.Itoa(c.Writer.Status())
		elapsed := float64(time.Since(start)) / float64(time.Second)
		resSz := float64(c.Writer.Size())

		ctx.reqDur.WithLabelValues(serverName).Observe(elapsed)
		ctx.reqCnt.WithLabelValues(serverName, status, c.Request.Method, c.Request.Host, c.Request.URL.Path).Inc()
		ctx.reqSz.WithLabelValues(serverName).Observe(float64(reqSz))
		ctx.resSz.WithLabelValues(serverName).Observe(resSz)
	}
}

// From https://github.com/DanielHeckrath/gin-prometheus/blob/master/gin_prometheus.go
func computeApproximateRequestSize(r *http.Request) int {
	s := 0
	if r.URL != nil {
		s = len(r.URL.Path)
	}

	s += len(r.Method)
	s += len(r.Proto)

	for name, values := range r.Header {
		s += len(name)
		for _, value := range values {
			s += len(value)
		}
	}

	s += len(r.Host)

	// N.B. r.Form and r.MultipartForm are assumed to be included in r.URL.

	if r.ContentLength != -1 {
		s += int(r.ContentLength)
	}

	return s
}

func (ctx *prometheusMetrics) register() {
	ctx.reqCnt = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "How many HTTP requests processed, partitioned by status code and HTTP method.",
		},
		[]string{"server_name", "status_code", "method", "host", "path"},
	)
	prometheus.MustRegister(ctx.reqCnt)

	ctx.reqDur = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "http_request_duration_seconds",
			Help: "The HTTP request latencies in seconds.",
		},
		[]string{"server_name"},
	)
	prometheus.MustRegister(ctx.reqDur)

	ctx.reqSz = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "http_request_size_bytes",
			Help: "The HTTP request sizes in bytes.",
		},
		[]string{"server_name"},
	)
	prometheus.MustRegister(ctx.reqSz)

	ctx.resSz = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "http_response_size_bytes",
			Help: "The HTTP response sizes in bytes.",
		},
		[]string{"server_name"},
	)
	prometheus.MustRegister(ctx.resSz)

	ctx.up = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "up",
			Help: "1 = up, 0 = down",
		},
	)
	ctx.up.Set(1)
	prometheus.MustRegister(ctx.up)

	// Register gqlgen graphql prometheus metrics
	gqlprometheus.Register()
}
