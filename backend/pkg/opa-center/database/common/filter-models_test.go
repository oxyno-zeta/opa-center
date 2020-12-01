//+build unit

package common

import (
	"reflect"
	"testing"
	"time"
)

func TestDateFilter_GetGenericFilter(t *testing.T) {
	notADate := "not a date"
	dateStr := "2020-09-19T23:10:35+02:00"
	date, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		t.Error(err)
		return
	}

	date = date.UTC()

	type fields struct {
		Eq        *string
		NotEq     *string
		Gte       *string
		NotGte    *string
		Gt        *string
		NotGt     *string
		Lte       *string
		NotLte    *string
		Lt        *string
		NotLt     *string
		In        []string
		NotIn     []string
		IsNull    bool
		IsNotNull bool
	}
	tests := []struct {
		name    string
		fields  fields
		want    *GenericFilter
		wantErr bool
	}{
		{
			name:   "Eq case",
			fields: fields{Eq: &dateStr},
			want:   &GenericFilter{Eq: &date},
		},
		{
			name:    "Eq not a date",
			fields:  fields{Eq: &notADate},
			wantErr: true,
		},
		{
			name:   "NotEq case",
			fields: fields{NotEq: &dateStr},
			want:   &GenericFilter{NotEq: &date},
		},
		{
			name:    "NotEq not a date",
			fields:  fields{NotEq: &notADate},
			wantErr: true,
		},
		{
			name:   "Gte case",
			fields: fields{Gte: &dateStr},
			want:   &GenericFilter{Gte: &date},
		},
		{
			name:    "Gte not a date",
			fields:  fields{Gte: &notADate},
			wantErr: true,
		},
		{
			name:   "NotGte case",
			fields: fields{NotGte: &dateStr},
			want:   &GenericFilter{NotGte: &date},
		},
		{
			name:    "NotGte not a date",
			fields:  fields{NotGte: &notADate},
			wantErr: true,
		},
		{
			name:   "Gt case",
			fields: fields{Gt: &dateStr},
			want:   &GenericFilter{Gt: &date},
		},
		{
			name:    "Gt not a date",
			fields:  fields{Gt: &notADate},
			wantErr: true,
		},
		{
			name:   "NotGt case",
			fields: fields{NotGt: &dateStr},
			want:   &GenericFilter{NotGt: &date},
		},
		{
			name:    "NotGt not a date",
			fields:  fields{NotGt: &notADate},
			wantErr: true,
		},
		{
			name:   "Lte case",
			fields: fields{Lte: &dateStr},
			want:   &GenericFilter{Lte: &date},
		},
		{
			name:    "Lte not a date",
			fields:  fields{Lte: &notADate},
			wantErr: true,
		},
		{
			name:   "NotLte case",
			fields: fields{NotLte: &dateStr},
			want:   &GenericFilter{NotLte: &date},
		},
		{
			name:    "NotLte not a date",
			fields:  fields{NotLte: &notADate},
			wantErr: true,
		},
		{
			name:   "Lt case",
			fields: fields{Lt: &dateStr},
			want:   &GenericFilter{Lt: &date},
		},
		{
			name:    "Lt not a date",
			fields:  fields{Lt: &notADate},
			wantErr: true,
		},
		{
			name:   "NotLt case",
			fields: fields{NotLt: &dateStr},
			want:   &GenericFilter{NotLt: &date},
		},
		{
			name:    "NotLt not a date",
			fields:  fields{NotLt: &notADate},
			wantErr: true,
		},
		{
			name:   "In case",
			fields: fields{In: []string{dateStr}},
			want:   &GenericFilter{In: []*time.Time{&date}},
		},
		{
			name:    "In not a date",
			fields:  fields{In: []string{notADate}},
			wantErr: true,
		},
		{
			name:   "NotIn case",
			fields: fields{NotIn: []string{dateStr}},
			want:   &GenericFilter{NotIn: []*time.Time{&date}},
		},
		{
			name:    "NotIn not a date",
			fields:  fields{NotIn: []string{notADate}},
			wantErr: true,
		},
		{
			name:   "IsNull case",
			fields: fields{IsNull: true},
			want:   &GenericFilter{IsNull: true},
		},
		{
			name:   "IsNotNull case",
			fields: fields{IsNotNull: true},
			want:   &GenericFilter{IsNotNull: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DateFilter{
				Eq:        tt.fields.Eq,
				NotEq:     tt.fields.NotEq,
				Gte:       tt.fields.Gte,
				NotGte:    tt.fields.NotGte,
				Gt:        tt.fields.Gt,
				NotGt:     tt.fields.NotGt,
				Lte:       tt.fields.Lte,
				NotLte:    tt.fields.NotLte,
				Lt:        tt.fields.Lt,
				NotLt:     tt.fields.NotLt,
				In:        tt.fields.In,
				NotIn:     tt.fields.NotIn,
				IsNull:    tt.fields.IsNull,
				IsNotNull: tt.fields.IsNotNull,
			}
			got, err := d.GetGenericFilter()
			if (err != nil) != tt.wantErr {
				t.Errorf("DateFilter.GetGenericFilter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DateFilter.GetGenericFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}
