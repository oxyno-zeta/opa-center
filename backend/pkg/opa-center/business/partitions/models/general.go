package models

import "github.com/oxyno-zeta/opa-center/pkg/opa-center/database/common"

type SortOrder struct {
	CreatedAt            *common.SortOrderEnum `dbfield:"created_at"`
	UpdatedAt            *common.SortOrderEnum `dbfield:"updated_at"`
	Name                 *common.SortOrderEnum `dbfield:"name"`
	StatusDataRetention  *common.SortOrderEnum `dbfield:"status_data_retention"`
	DecisionLogRetention *common.SortOrderEnum `dbfield:"decision_log_retention"`
}

type Filter struct {
	AND                  []*Filter
	OR                   []*Filter
	CreatedAt            *common.DateFilter    `dbfield:"created_at"`
	UpdatedAt            *common.DateFilter    `dbfield:"updated_at"`
	Name                 *common.GenericFilter `dbfield:"name"`
	StatusDataRetention  *common.GenericFilter `dbfield:"status_data_retention"`
	DecisionLogRetention *common.GenericFilter `dbfield:"decision_log_retention"`
}

type Projection struct {
	ID                   bool `dbfield:"id" graphqlfield:"id"`
	CreatedAt            bool `dbfield:"created_at" graphqlfield:"createdAt"`
	UpdatedAt            bool `dbfield:"updated_at" graphqlfield:"updatedAt"`
	Name                 bool `dbfield:"name" graphqlfield:"name"`
	StatusDataRetention  bool `dbfield:"status_data_retention" graphqlfield:"statusDataRetention"`
	DecisionLogRetention bool `dbfield:"decision_log_retention" graphqlfield:"decisionLogRetention"`
}

type CreateInput struct {
	Name                 string `validate:"required,max=255"`
	StatusDataRetention  string `validate:"omitempty,max=255"`
	DecisionLogRetention string `validate:"omitempty,max=255"`
}
