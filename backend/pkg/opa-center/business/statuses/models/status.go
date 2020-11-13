package models

import "time"

type Status struct {
	ID              string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	OriginalMessage string `validate:"required"`
	PartitionID     string
}
