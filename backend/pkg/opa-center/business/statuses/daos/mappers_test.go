//+build unit

package daos

import (
	"testing"
	"time"

	daomodels "github.com/oxyno-zeta/opa-center/pkg/opa-center/business/statuses/daos/models"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/business/statuses/models"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/database"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
)

func Test_toDao(t *testing.T) {
	now := time.Now()

	type args struct {
		ins *models.Status
	}
	tests := []struct {
		name string
		args args
		want *daomodels.Status
	}{
		{
			name: "empty input",
			args: args{
				ins: &models.Status{},
			},
			want: &daomodels.Status{
				OriginalMessage: datatypes.JSON{},
			},
		},
		{
			name: "mapper",
			args: args{
				ins: &models.Status{
					CreatedAt: now,
					UpdatedAt: now,
					ID:        "fake id",
					OriginalMessage: `
{"key1":"val1","key2":["val21","val22"]}
					`,
					PartitionID: "fake pid",
				},
			},
			want: &daomodels.Status{
				Base: database.Base{
					CreatedAt: now,
					UpdatedAt: now,
					ID:        "fake id",
				},
				PartitionID: "fake pid",
				OriginalMessage: datatypes.JSON([]byte(`
{"key1":"val1","key2":["val21","val22"]}
					`)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toDao(tt.args.ins)
			if got != nil {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_fromDao(t *testing.T) {
	now := time.Now()

	type args struct {
		ins *daomodels.Status
	}
	tests := []struct {
		name    string
		args    args
		want    *models.Status
		wantErr bool
	}{
		{
			name: "mapper",
			args: args{
				ins: &daomodels.Status{
					Base: database.Base{
						ID:        "fake id",
						CreatedAt: now,
						UpdatedAt: now,
					},
					OriginalMessage: datatypes.JSON([]byte(`{"key1":"val1"}`)),
					PartitionID:     "fake pid",
				},
			},
			want: &models.Status{
				CreatedAt:       now,
				UpdatedAt:       now,
				ID:              "fake id",
				OriginalMessage: `{"key1":"val1"}`,
				PartitionID:     "fake pid",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fromDao(tt.args.ins)
			if (err != nil) != tt.wantErr {
				t.Errorf("fromDao() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
