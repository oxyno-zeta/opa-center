package daos

import (
	daomodels "github.com/oxyno-zeta/opa-center/pkg/opa-center/business/statuses/daos/models"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/business/statuses/models"
	"gorm.io/datatypes"
)

func toDao(ins *models.Status) *daomodels.Status {
	val := &daomodels.Status{
		OriginalMessage: datatypes.JSON([]byte(ins.OriginalMessage)),
		PartitionID:     ins.PartitionID,
	}
	// Add other data
	val.ID = ins.ID
	val.CreatedAt = ins.CreatedAt
	val.UpdatedAt = ins.UpdatedAt

	return val
}

func fromDao(ins *daomodels.Status) (*models.Status, error) {
	// Transform in string
	bb, err := ins.OriginalMessage.MarshalJSON()
	// Check error
	if err != nil {
		return nil, err
	}

	val := &models.Status{
		ID:              ins.ID,
		CreatedAt:       ins.CreatedAt,
		UpdatedAt:       ins.UpdatedAt,
		OriginalMessage: string(bb),
		PartitionID:     ins.PartitionID,
	}

	return val, nil
}
