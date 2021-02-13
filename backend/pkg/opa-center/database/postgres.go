package database

import (
	"database/sql"
	"time"

	"github.com/oxyno-zeta/opa-center/pkg/opa-center/config"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/log"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/metrics"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"gorm.io/driver/postgres"
)

type postresdb struct {
	logger         log.Logger
	db             *gorm.DB
	cfgManager     config.Manager
	connectionName string
	metricsCl      metrics.Client
}

func (ctx *postresdb) GetSQLDB() (*sql.DB, error) {
	return ctx.db.DB()
}

// GetGormDB will return a gorm database object.
func (ctx *postresdb) GetGormDB() *gorm.DB {
	return ctx.db
}

// Connect will connect to database engine.
func (ctx *postresdb) Connect() error {
	// Get configuration
	cfg := ctx.cfgManager.GetConfig()

	var sqlConnectionMaxLifetimeDuration time.Duration
	// Try to parse sql connection max lifetime duration
	if cfg.Database.SQLConnectionMaxLifetimeDuration != "" {
		var err error
		sqlConnectionMaxLifetimeDuration, err = time.ParseDuration(cfg.Database.SQLConnectionMaxLifetimeDuration)
		// Check error
		if err != nil {
			return err
		}
	}

	// Create gorm configuration
	gcfg := &gorm.Config{
		// Insert now function to be sure that automatic dates are in UTC
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		// Add logger
		Logger: ctx.logger.GetGormLogger(),
		// Disable foreign key constraint when migrating
		DisableForeignKeyConstraintWhenMigrating: cfg.Database.DisableForeignKeyWhenMigrating,
		// Allow global update
		AllowGlobalUpdate: cfg.Database.AllowGlobalUpdate,
		// Prepare statement for caching
		PrepareStmt: cfg.Database.PrepareStatement,
	}

	ctx.logger.Debug("Trying to connect to database engine of type PostgreSQL")
	// Connect to database
	dbResult, err := gorm.Open(postgres.Open(cfg.Database.ConnectionURL.Value), gcfg)
	// Check if error exists
	if err != nil {
		return errors.WithStack(err)
	}

	// Get prometheus gorm middleware
	md := ctx.metricsCl.DatabaseMiddleware(ctx.connectionName)
	// Apply middleware
	err = dbResult.Use(md)
	// Check if error exists
	if err != nil {
		return errors.WithStack(err)
	}

	// Trying to ping database
	sqlDB, err := dbResult.DB()
	// Check error
	if err != nil {
		return errors.WithStack(err)
	}

	// Ping db
	err = sqlDB.Ping()
	// Check error
	if err != nil {
		return errors.WithStack(err)
	}

	// Check if max idle connections exists
	if cfg.Database.SQLMaxIdleConnections != 0 {
		// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
		sqlDB.SetMaxIdleConns(cfg.Database.SQLMaxIdleConnections)
	}

	// Check if max opened connections exists
	if cfg.Database.SQLMaxOpenConnections != 0 {
		// SetMaxOpenConns sets the maximum number of open connections to the database.
		sqlDB.SetMaxOpenConns(cfg.Database.SQLMaxOpenConnections)
	}

	// Check if connection max lifetime exists
	if cfg.Database.SQLConnectionMaxLifetimeDuration != "" {
		// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
		sqlDB.SetConnMaxLifetime(sqlConnectionMaxLifetimeDuration)
	}

	// Save gorm db object
	ctx.db = dbResult

	ctx.logger.Info("Successfully connected to database engine of type PostgreSQL")

	// Return
	return nil
}

// Close will close connection to database.
func (ctx *postresdb) Close() error {
	ctx.logger.Info("Closing database connection")
	// Get sql database
	sqlDB, err := ctx.db.DB()
	// Check error
	if err != nil {
		return errors.WithStack(err)
	}

	return sqlDB.Close()
}

// Ping will ping database engine in order to test connection to engine.
func (ctx *postresdb) Ping() error {
	// Get sql database
	sqlDB, err := ctx.db.DB()
	// Check error
	if err != nil {
		return errors.WithStack(err)
	}
	// Ping database to test connection
	pingErr := sqlDB.Ping()
	// Check error
	if pingErr != nil {
		return errors.WithStack(pingErr)
	}

	return nil
}

func (ctx *postresdb) Reconnect() error {
	// Get old sql db object
	sqlDB, err := ctx.db.DB()
	// Check error
	if err != nil {
		return errors.WithStack(err)
	}
	// Defer closing old database connection
	defer sqlDB.Close()
	// Connect to new database
	err = ctx.Connect()
	if err != nil {
		return errors.WithStack(err)
	}
	// Wait for 1 sec before closing connection to old db
	// Here, we suppose that waiting 1 second is enough for reload
	time.Sleep(time.Second)

	return nil
}
