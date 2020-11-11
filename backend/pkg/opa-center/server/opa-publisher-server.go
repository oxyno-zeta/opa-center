package server

import (
	"net/http"
	"strconv"

	helmet "github.com/danielkov/gin-helmet"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/authx/authentication"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/business"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/config"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/log"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/metrics"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/server/middlewares"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/server/rest"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/tracing"
)

type OPAPublisherServer struct {
	logger            log.Logger
	cfgManager        config.Manager
	metricsCl         metrics.Client
	tracingSvc        tracing.Service
	busiServices      *business.Services
	authenticationSvc authentication.Client
	server            *http.Server
}

func NewOPAPublisherServer(
	logger log.Logger, cfgManager config.Manager, metricsCl metrics.Client,
	tracingSvc tracing.Service, busiServices *business.Services,
	authenticationSvc authentication.Client,
) *OPAPublisherServer {
	return &OPAPublisherServer{
		logger:            logger,
		cfgManager:        cfgManager,
		metricsCl:         metricsCl,
		tracingSvc:        tracingSvc,
		busiServices:      busiServices,
		authenticationSvc: authenticationSvc,
	}
}

func (svr *OPAPublisherServer) GenerateServer() error {
	// Get configuration
	cfg := svr.cfgManager.GetConfig()
	// Generate router
	r := svr.generateRouter()

	// Create server
	addr := cfg.OPAPublisherServer.ListenAddr + ":" + strconv.Itoa(cfg.OPAPublisherServer.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// Prepare for configuration onChange
	svr.cfgManager.AddOnChangeHook(func() {
		// Generate router
		r := svr.generateRouter()
		// Change server handler
		server.Handler = r
		svr.logger.Info("OPA Publisher Server handler reloaded")
	})

	// Store server
	svr.server = server

	return nil
}

func (svr *OPAPublisherServer) generateRouter() http.Handler {
	// Get configuration
	// cfg := svr.cfgManager.GetConfig()
	// Set release mod
	gin.SetMode(gin.ReleaseMode)
	// Create router
	router := gin.New()
	// Manage no route
	router.NoRoute(func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "404 not found"})
	})
	// Add middlewares
	router.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))
	router.Use(gin.Recovery())
	router.Use(helmet.Default())
	router.Use(middlewares.RequestID(svr.logger))
	router.Use(svr.tracingSvc.Middleware(middlewares.GetRequestIDFromContext))
	router.Use(log.Middleware(svr.logger, middlewares.GetRequestIDFromGin, tracing.GetSpanIDFromContext))
	router.Use(svr.metricsCl.Instrument("opa-publisher"))

	// Add authentication
	// if cfg.OPAPublisherAuthentication != nil && cfg.OPAPublisherAuthentication.BasicAuthAccounts != nil {
	// 	// Create gin auth basic
	// 	abasic := gin.Accounts{}
	// 	// Loop over users
	// 	for k, v := range cfg.OPAPublisherAuthentication.BasicAuthAccounts {
	// 		// Add user
	// 		abasic[k] = v.Value
	// 	}
	// 	// Add it to router
	// 	router.Use(gin.BasicAuth(abasic))
	// }

	// Add REST endpoints
	rest.AddDecisionLogsEndpoints(router, svr.busiServices)

	return router
}

func (svr *OPAPublisherServer) Listen() error {
	svr.logger.Infof("OPA Publisher Server listening on %s", svr.server.Addr)
	err := svr.server.ListenAndServe()

	return err
}
