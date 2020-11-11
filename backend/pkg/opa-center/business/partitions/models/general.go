package models

import "github.com/oxyno-zeta/opa-center/pkg/opa-center/database/common"

type SortOrder struct {
	CreatedAt *common.SortOrderEnum `dbfield:"created_at"`
	UpdatedAt *common.SortOrderEnum `dbfield:"updated_at"`
	Name      *common.SortOrderEnum `dbfield:"name"`
}

type Filter struct {
	CreatedAt *common.DateFilter    `dbfield:"created_at"`
	UpdatedAt *common.DateFilter    `dbfield:"updated_at"`
	Name      *common.GenericFilter `dbfield:"name"`
}

type Projection struct {
	ID        bool `dbfield:"id" graphqlfield:"id"`
	CreatedAt bool `dbfield:"created_at" graphqlfield:"createdAt"`
	UpdatedAt bool `dbfield:"updated_at" graphqlfield:"updatedAt"`
	Name      bool `dbfield:"name" graphqlfield:"name"`
}

type CreateInput struct {
	Name string `validate:"required,max=255"`
}
