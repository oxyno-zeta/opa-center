package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/InVisionApp/go-health/v2"
	"github.com/InVisionApp/go-health/v2/handlers"
	helmet "github.com/danielkov/gin-helmet"
	"github.com/gin-gonic/gin"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/config"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/log"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/metrics"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/server/middlewares"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/tracing"
)

type InternalServer struct {
	logger     log.Logger
	cfgManager config.Manager
	metricsCl  metrics.Client
	checkers   []*health.Config
	server     *http.Server
}

type CheckerInput struct {
	Name     string
	CheckFn  func() error
	Interval time.Duration
}

func NewInternalServer(logger log.Logger, cfgManager config.Manager, metricsCl metrics.Client) *InternalServer {
	return &InternalServer{
		logger:     logger,
		cfgManager: cfgManager,
		metricsCl:  metricsCl,
		checkers:   make([]*health.Config, 0),
	}
}

// AddChecker allow to add a health checker.
func (svr *InternalServer) AddChecker(chI *CheckerInput) {
	svr.checkers = append(svr.checkers, &health.Config{
		Name:     chI.Name,
		Checker:  &customHealthChecker{fn: chI.CheckFn},
		Fatal:    true,
		Interval: chI.Interval,
	})
}

func (svr *InternalServer) generateInternalRouter() (http.Handler, error) {
	// Get configuration
	cfg := svr.cfgManager.GetConfig()

	// Set release mod
	gin.SetMode(gin.ReleaseMode)
	// Create router
	router := gin.New()
	// Add middlewares
	router.Use(gin.Recovery())
	router.Use(helmet.Default())
	router.Use(middlewares.RequestID(svr.logger))
	router.Use(log.Middleware(svr.logger, middlewares.GetRequestIDFromGin, tracing.GetSpanIDFromContext))
	router.Use(svr.metricsCl.Instrument("internal"))
	// Add cors if configured
	err := manageCORS(router, cfg.InternalServer)
	// Check error
	if err != nil {
		return nil, err
	}

	// Create a new health instance
	h := health.New()

	// Disable logging
	h.DisableLogging()

	// Add checkers
	err = h.AddChecks(svr.checkers)
	if err != nil {
		return nil, err
	}

	// Start health checker
	err = h.Start()
	if err != nil {
		return nil, err
	}

	// Add metrics path
	router.GET("/metrics", gin.WrapH(svr.metricsCl.PrometheusHTTPHandler()))
	router.GET("/health", gin.WrapH(handlers.NewJSONHandlerFunc(h, nil)))

	return router, nil
}

func (svr *InternalServer) Listen() error {
	svr.logger.Infof("Internal server listening on %s", svr.server.Addr)
	err := svr.server.ListenAndServe()

	return err
}

func (svr *InternalServer) GenerateServer() error {
	// Get configuration
	cfg := svr.cfgManager.GetConfig()
	// Generate internal router
	r, err := svr.generateInternalRouter()
	if err != nil {
		return err
	}
	// Create server
	addr := cfg.InternalServer.ListenAddr + ":" + strconv.Itoa(cfg.InternalServer.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}
	// Store server
	svr.server = server

	return nil
}
