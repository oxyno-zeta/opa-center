package models

import (
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/database"
	"gorm.io/datatypes"
)

type Status struct {
	database.Base
	OriginalMessage datatypes.JSON
	PartitionID     string `gorm:"index"`
}
