package models

import "github.com/oxyno-zeta/opa-center/pkg/opa-center/database/common"

type SortOrder struct {
	CreatedAt *common.SortOrderEnum `dbfield:"created_at"`
	UpdatedAt *common.SortOrderEnum `dbfield:"updated_at"`
}

type Filter struct {
	AND         []*Filter
	OR          []*Filter
	CreatedAt   *common.DateFilter    `dbfield:"created_at"`
	UpdatedAt   *common.DateFilter    `dbfield:"updated_at"`
	PartitionID *common.GenericFilter `dbfield:"partition_id"`
}

type Projection struct {
	ID              bool `dbfield:"id" graphqlfield:"id"`
	CreatedAt       bool `dbfield:"created_at" graphqlfield:"createdAt"`
	UpdatedAt       bool `dbfield:"updated_at" graphqlfield:"updatedAt"`
	OriginalMessage bool `dbfield:"original_message" graphqlfield:"originalMessage"`
	PartitionID     bool `dbfield:"partition_id" graphqlfield:"partition"`
}
