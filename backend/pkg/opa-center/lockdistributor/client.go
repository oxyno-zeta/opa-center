package lockdistributor

import (
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/config"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/database"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/log"
)

//go:generate mockgen -destination=./mocks/mock_Service.go -package=mocks github.com/oxyno-zeta/opa-center/pkg/opa-center/lockdistributor Service
type Service interface {
	// Get a lock object (semaphore on string) that can be acquired and release
	GetLock(name string) Lock
	// Initialize service
	Initialize(logger log.Logger) error
}

//go:generate mockgen -destination=./mocks/mock_Lock.go -package=mocks github.com/oxyno-zeta/opa-center/pkg/opa-center/lockdistributor Lock
type Lock interface {
	// Acquire lock
	Acquire() error
	// Release lock
	Release() error
	// Check if a lock with this name is already taken
	IsAlreadyTaken() (bool, error)
	// Check if the lock is released or lost because of missing heartbeat
	IsReleased() (bool, error)
}

func NewService(cfgManager config.Manager, db database.DB) Service {
	return &service{
		cfgManager: cfgManager,
		db:         db,
	}
}
