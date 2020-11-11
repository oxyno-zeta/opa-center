package models

import (
	"time"
)

type DecisionLog struct {
	ID              string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DecisionID      string    `validate:"required,max=255"`
	Path            string    `validate:"required,max=255"`
	RequestedBy     string    `validate:"required,max=255"`
	Timestamp       time.Time `validate:"required"`
	OriginalMessage string    `validate:"required"`
	PartitionID     string
}
