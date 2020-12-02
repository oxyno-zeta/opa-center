package models

import (
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/database"
)

type Partition struct {
	database.Base
	Name                 string `gorm:"unique_index"`
	StatusDataRetention  string
	DecisionLogRetention string
}
