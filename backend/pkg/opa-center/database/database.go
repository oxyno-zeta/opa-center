package database

import (
	"database/sql"

	"github.com/oxyno-zeta/opa-center/pkg/opa-center/config"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/log"
	"gorm.io/gorm"
)

//go:generate mockgen -destination=./mocks/mock_DB.go -package=mocks github.com/oxyno-zeta/opa-center/pkg/opa-center/database DB
type DB interface {
	// Get Gorm db object.
	GetGormDB() *gorm.DB
	// Get SQL db object.
	GetSQLDB() (*sql.DB, error)
	// Connect to database.
	Connect() error
	// Close database connection.
	Close() error
	// Ping database.
	Ping() error
	// Reconnect to database.
	Reconnect() error
}

// NewDatabase will generate a new DB object.
func NewDatabase(cfgManager config.Manager, logger log.Logger) DB {
	return &postresdb{logger: logger, cfgManager: cfgManager}
}
