package models

import (
	"time"

	"github.com/oxyno-zeta/opa-center/pkg/opa-center/database"
	"gorm.io/datatypes"
)

type DecisionLog struct {
	database.Base
	DecisionID      string `gorm:"unique_index"`
	Path            string
	RequestedBy     string
	Timestamp       time.Time
	OriginalMessage datatypes.JSON
	PartitionID     string `gorm:"index"`
}
