package daos

import (
	daomodels "github.com/oxyno-zeta/opa-center/pkg/opa-center/business/decisionlogs/daos/models"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/business/decisionlogs/models"
	"gorm.io/datatypes"
)

func toDao(ins *models.DecisionLog) *daomodels.DecisionLog {
	val := &daomodels.DecisionLog{
		DecisionID:      ins.DecisionID,
		Path:            ins.Path,
		RequestedBy:     ins.RequestedBy,
		Timestamp:       ins.Timestamp,
		OriginalMessage: datatypes.JSON([]byte(ins.OriginalMessage)),
		PartitionID:     ins.PartitionID,
	}
	// Add other data
	val.ID = ins.ID
	val.CreatedAt = ins.CreatedAt
	val.UpdatedAt = ins.UpdatedAt

	return val
}

func fromDao(ins *daomodels.DecisionLog) (*models.DecisionLog, error) {
	// Transform in string
	bb, err := ins.OriginalMessage.MarshalJSON()
	// Check error
	if err != nil {
		return nil, err
	}

	val := &models.DecisionLog{
		ID:              ins.ID,
		CreatedAt:       ins.CreatedAt,
		UpdatedAt:       ins.UpdatedAt,
		DecisionID:      ins.DecisionID,
		Path:            ins.Path,
		RequestedBy:     ins.RequestedBy,
		Timestamp:       ins.Timestamp,
		OriginalMessage: string(bb),
		PartitionID:     ins.PartitionID,
	}

	return val, nil
}
