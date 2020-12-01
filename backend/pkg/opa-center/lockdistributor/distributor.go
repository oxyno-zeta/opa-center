package lockdistributor

import (
	"strings"
	"time"

	"cirello.io/pglock"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/config"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/database"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/log"
)

type service struct {
	cfgManager config.Manager
	db         database.DB
	cl         *pglock.Client
}

func (s *service) Initialize(logger log.Logger) error {
	// Get configuration
	cfg := s.cfgManager.GetConfig()

	// Parse durations
	ld, err := time.ParseDuration(cfg.LockDistributor.LeaseDuration)
	if err != nil {
		return err
	}

	hf, err := time.ParseDuration(cfg.LockDistributor.HeartbeatFrequency)
	if err != nil {
		return err
	}

	// Get sql database
	sqlDB, err := s.db.GetSQLDB()
	// Check error
	if err != nil {
		return err
	}

	// Log
	logger.Debug("Trying to create lock distributor client")

	// Create pglock client
	c, err := pglock.UnsafeNew(
		sqlDB,
		pglock.WithLeaseDuration(ld),
		pglock.WithHeartbeatFrequency(hf),
		pglock.WithCustomTable(cfg.LockDistributor.TableName),
		pglock.WithLogger(logger.GetLockDistributorLogger()),
	)
	// Check error
	if err != nil {
		return err
	}

	// Create lock table
	err = c.CreateTable()
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		return err
	}

	// Save client
	s.cl = c

	// Log
	logger.Info("Successfully created lock distributor client")

	return nil
}

func (s *service) GetLock(name string) Lock {
	return &lock{
		name: name,
		s:    s,
	}
}
