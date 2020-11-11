//+build unit

package common

import (
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestManageProjection(t *testing.T) {
	type Person struct {
		Name string
	}
	type Proj1 struct {
		Field1 bool
	}
	type Proj2 struct {
		Field1 bool `dbfield:""`
	}
	type Proj3 struct {
		Field1 bool `dbfield:"-"`
	}
	type Proj4 struct {
		Field1 string `dbfield:"field1"`
	}
	type Proj5 struct {
		Field1 bool `dbfield:"field1"`
		Field2 bool `dbfield:"field2"`
		Field3 bool `dbfield:"field3"`
	}
	type args struct {
		projection interface{}
	}
	tests := []struct {
		name                      string
		args                      args
		expectedIntermediateQuery string
		expectedArgs              []driver.Value
		wantErr                   bool
		errorString               string
	}{
		{
			name:        "wrong input",
			args:        args{projection: false},
			wantErr:     true,
			errorString: "projection must be an object",
		},
		{
			name: "nil sort object",
			args: args{
				projection: nil,
			},
			expectedIntermediateQuery: "*",
		},
		{
			name: "no tag",
			args: args{
				projection: &Proj1{Field1: true},
			},
			expectedIntermediateQuery: "*",
		},
		{
			name: "field ignored: empty",
			args: args{
				projection: &Proj2{Field1: true},
			},
			expectedIntermediateQuery: "*",
		},
		{
			name: "field ignored: -",
			args: args{
				projection: &Proj3{Field1: true},
			},
			expectedIntermediateQuery: "*",
		},
		{
			name: "not a boolean with tag",
			args: args{
				projection: &Proj4{Field1: "true"},
			},
			wantErr:     true,
			errorString: "field Field1 with projection tag must be a boolean",
		},
		{
			name: "multiple field in projection",
			args: args{
				projection: &Proj5{
					Field1: true,
					Field2: false,
					Field3: true,
				},
			},
			expectedIntermediateQuery: "field1,field3",
		},
		{
			name: "multiple field in projection with all fields set",
			args: args{
				projection: &Proj5{
					Field1: true,
					Field2: true,
					Field3: true,
				},
			},
			expectedIntermediateQuery: "field1,field2,field3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Error(err)
				return
			}
			defer sqlDB.Close()

			db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{Logger: logger.Discard})
			if err != nil {
				t.Error(err)
				return
			}

			got, err := ManageProjection(tt.args.projection, db)
			if (err != nil) != tt.wantErr {
				t.Errorf("ManageProjection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.errorString {
				t.Errorf("ManageProjection() error = %v, wantErr %v", err, tt.errorString)
				return
			}
			if err != nil {
				return
			}

			// Create expected query
			expectedQuery := `SELECT ` + tt.expectedIntermediateQuery + ` FROM "people" `
			expectedQuery += `ORDER BY "people"."name" LIMIT 1`

			mock.ExpectQuery(expectedQuery).
				WithArgs(tt.expectedArgs...).
				WillReturnRows(
					sqlmock.NewRows([]string{"name"}).AddRow("fake"),
				)

			// Run fake find to force query to be run
			res := got.First(&Person{})
			// Test error
			if res.Error != nil {
				t.Error(res.Error)
			}
		})
	}
}
