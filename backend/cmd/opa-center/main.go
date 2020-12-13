package main

import (
	"time"

	"github.com/oxyno-zeta/opa-center/pkg/opa-center/authx/authentication"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/authx/authorization"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/business"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/config"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/database"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/log"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/metrics"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/server"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/tracing"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/version"
	"golang.org/x/sync/errgroup"
)

func main() {
	// Create new logger
	logger := log.NewLogger()

	// Create configuration manager
	cfgManager := config.NewManager(logger)

	// Load configuration
	err := cfgManager.Load()
	if err != nil {
		logger.WithError(err).Fatal(err)
	}

	// Get configuration
	cfg := cfgManager.GetConfig()
	// Configure logger
	err = logger.Configure(cfg.Log.Level, cfg.Log.Format, cfg.Log.FilePath)
	if err != nil {
		logger.WithError(err).Fatal(err)
	}

	// Watch change for logger (special case)
	cfgManager.AddOnChangeHook(func() {
		// Get configuration
		cfg := cfgManager.GetConfig()
		// Configure logger
		err = logger.Configure(cfg.Log.Level, cfg.Log.Format, cfg.Log.FilePath)
		if err != nil {
			logger.WithError(err).Error(err)
		}
	})

	// Getting version
	v := version.GetVersion()
	logger.Infof("Starting version: %s (git commit: %s) built on %s", v.Version, v.GitCommit, v.BuildDate)

	// Create metrics client
	metricsCl := metrics.NewMetricsClient()

	// Generate tracing service instance
	tracingSvc, err := tracing.New(cfgManager, logger)
	// Check error
	if err != nil {
		logger.WithError(err).Fatal(err)
	}
	// Prepare on reload hook
	cfgManager.AddOnChangeHook(func() {
		err = tracingSvc.Reload()
		if err != nil {
			logger.WithError(err).Fatal(err)
		}
	})

	// Create database service
	db := database.NewDatabase(cfgManager, logger)
	// Connect to engine
	err = db.Connect()
	if err != nil {
		logger.WithError(err).Fatal(err)
	}
	// Add configuration reload hook
	cfgManager.AddOnChangeHook(func() {
		err = db.Reconnect()
		if err != nil {
			logger.WithError(err).Fatal(err)
		}
	})

	// Create authentication service
	authoSvc := authorization.NewService(cfgManager)

	// Create business services
	busServices, err := business.NewServices(logger, db, authoSvc, cfgManager)
	// Check error
	if err != nil {
		logger.WithError(err).Fatal(err)
	}

	// Migrate database
	err = busServices.MigrateDB()
	if err != nil {
		logger.WithError(err).Fatal(err)
	}

	// Initialize services
	err = busServices.Initialize()
	if err != nil {
		logger.WithError(err).Fatal(err)
	}
	// Add configuration reload hook
	cfgManager.AddOnChangeHook(func() {
		err = busServices.Reload()
		if err != nil {
			logger.WithError(err).Fatal(err)
		}
	})

	// Create authentication service
	authenticationSvc := authentication.NewService(cfgManager)

	// Create servers
	svr := server.NewServer(logger, cfgManager, metricsCl, tracingSvc, busServices, authenticationSvc, authoSvc)
	opaSvr := server.NewOPAPublisherServer(logger, cfgManager, metricsCl, tracingSvc, busServices, authenticationSvc)
	intSvr := server.NewInternalServer(logger, cfgManager, metricsCl)

	// Generate server
	err = svr.GenerateServer()
	if err != nil {
		logger.WithError(err).Fatal(err)
	}
	// Generate internal server
	err = intSvr.GenerateServer()
	if err != nil {
		logger.WithError(err).Fatal(err)
	}
	// Generate opa server
	err = opaSvr.GenerateServer()
	if err != nil {
		logger.WithError(err).Fatal(err)
	}

	// Add checker for internal server
	intSvr.AddChecker(&server.CheckerInput{
		Name:     "database",
		CheckFn:  db.Ping,
		Interval: 2 * time.Second, //nolint:gomnd // Won't do a const for that
	})

	var g errgroup.Group

	g.Go(svr.Listen)
	g.Go(intSvr.Listen)
	g.Go(opaSvr.Listen)

	if err := g.Wait(); err != nil {
		logger.WithError(err).Fatal(err)
	}
}
