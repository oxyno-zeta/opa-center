package business

import (
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/authx/authorization"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/business/decisionlogs"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/business/partitions"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/business/statuses"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/config"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/database"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/log"
)

type Services struct {
	systemLogger    log.Logger
	DecisionLogsSvc decisionlogs.Service
	PartitionsSvc   partitions.Service
	StatusSvc       statuses.Service
}

func (s *Services) MigrateDB() error {
	funcs := []func(log.Logger) error{
		s.DecisionLogsSvc.MigrateDB,
		s.PartitionsSvc.MigrateDB,
		s.StatusSvc.MigrateDB,
	}

	// Loop over all migrations
	for i := 0; i < len(funcs); i++ {
		f := funcs[i]
		// Execute
		err := f(s.systemLogger)
		// Check error
		if err != nil {
			return err
		}
	}

	// Default case
	return nil
}

func NewServices(systemLogger log.Logger, db database.DB, authSvc authorization.Service, cfgManager config.Manager) (*Services, error) {
	// Create partitions service
	pSvc, err := partitions.NewService(db, authSvc, cfgManager)
	// Check error
	if err != nil {
		return nil, err
	}
	// Create decision logs service
	dlSvc := decisionlogs.NewService(db, authSvc, pSvc)
	// Create status service
	stSvc := statuses.NewService(db, authSvc, pSvc)

	return &Services{
		systemLogger:    systemLogger,
		DecisionLogsSvc: dlSvc,
		PartitionsSvc:   pSvc,
		StatusSvc:       stSvc,
	}, nil
}
