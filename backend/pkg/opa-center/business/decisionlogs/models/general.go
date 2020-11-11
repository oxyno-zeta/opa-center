package models

import "github.com/oxyno-zeta/opa-center/pkg/opa-center/database/common"

type SortOrder struct {
	CreatedAt   *common.SortOrderEnum `dbfield:"created_at"`
	UpdatedAt   *common.SortOrderEnum `dbfield:"updated_at"`
	Path        *common.SortOrderEnum `dbfield:"path"`
	RequestedBy *common.SortOrderEnum `dbfield:"requested_by"`
	Timestamp   *common.SortOrderEnum `dbfield:"timestamp"`
}

type Filter struct {
	CreatedAt   *common.DateFilter    `dbfield:"created_at"`
	UpdatedAt   *common.DateFilter    `dbfield:"updated_at"`
	DecisionID  *common.GenericFilter `dbfield:"decision_id"`
	Path        *common.GenericFilter `dbfield:"path"`
	RequestedBy *common.GenericFilter `dbfield:"requested_by"`
	Timestamp   *common.DateFilter    `dbfield:"timestamp"`
	PartitionID *common.GenericFilter `dbfield:"partition_id"`
}

type Projection struct {
	ID              bool `dbfield:"id" graphqlfield:"id"`
	CreatedAt       bool `dbfield:"created_at" graphqlfield:"createdAt"`
	UpdatedAt       bool `dbfield:"updated_at" graphqlfield:"updatedAt"`
	DecisionID      bool `dbfield:"decision_id" graphqlfield:"decisionId"`
	Path            bool `dbfield:"path" graphqlfield:"path"`
	RequestedBy     bool `dbfield:"requested_by" graphqlfield:"requestedBy"`
	Timestamp       bool `dbfield:"timestamp" graphqlfield:"timestamp"`
	OriginalMessage bool `dbfield:"original_message" graphqlfield:"originalMessage"`
	PartitionID     bool `dbfield:"partition_id" graphqlfield:"partition"`
}
