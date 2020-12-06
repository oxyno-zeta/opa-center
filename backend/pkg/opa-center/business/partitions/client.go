package partitions

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/authx/authorization"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/business/partitions/daos"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/business/partitions/models"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/config"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/database"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/database/pagination"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/log"
)

var errInvalidNameTemplate = "name must match regex %s"

type Service interface {
	// Initialize service
	Initialize() error
	// Reload service
	Reload() error
	// Add services
	AddServices(decisionLogsSvc, statusesSvc RetentionService)
	// Migrate database
	MigrateDB(systemLogger log.Logger) error
	// Get data paginated
	GetAllPaginated(
		ctx context.Context,
		page *pagination.PageInput,
		sort *models.SortOrder,
		filter *models.Filter,
		projection *models.Projection,
	) ([]*models.Partition, *pagination.PageOutput, error)
	// Create partition
	Create(ctx context.Context, inp *models.CreateInput) (*models.Partition, error)
	// Update partition
	Update(ctx context.Context, inp *models.UpdateInput) (*models.Partition, error)
	// Find by id used internally only
	UnsecureFindByID(id string) (*models.Partition, error)
	// Find by id
	FindByID(ctx context.Context, id string, projection *models.Projection) (*models.Partition, error)
	// Generate OPA configuration
	GenerateOPAConfiguration(ctx context.Context, id string) (string, error)
}

type RetentionService interface {
	ManageRetention(logger log.Logger, retentionDuration time.Duration, partitionID string) error
}

func NewService(db database.DB, authorizationSvc authorization.Service, cfgManager config.Manager, logger log.Logger) (Service, error) {
	// Create dao
	dao := daos.NewDao(db)
	// Create template
	opaCfgTemplate, err := loadOpaCfgTemplate()
	// Check error
	if err != nil {
		return nil, err
	}

	return &service{
		dao:              dao,
		validator:        validator.New(),
		authorizationSvc: authorizationSvc,
		cfgManager:       cfgManager,
		opaCfgTemplate:   opaCfgTemplate,
		logger:           logger,
	}, nil
}
