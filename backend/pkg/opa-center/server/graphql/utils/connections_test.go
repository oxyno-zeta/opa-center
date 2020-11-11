// +build unit

package utils

import (
	"testing"

	"github.com/oxyno-zeta/opa-center/pkg/opa-center/database/pagination"
	"github.com/stretchr/testify/assert"
)

func TestMapConnection(t *testing.T) {
	starString := func(s string) *string { return &s }
	type Person struct{ Name string }
	type PersonEdge struct {
		Cursor string
		Node   *Person
	}
	type WrongEdge1 struct {
	}
	type WrongEdge2 struct {
		Cursor int
		Node   *Person
	}
	type WrongEdge3 struct {
		Cursor string
		Node   string
	}
	type WrongEdge4 struct {
		Cursor string
		Node   Person
	}
	type WrongEdge5 struct {
		Cursor string
	}
	type PersonConnection struct {
		Edges    []*PersonEdge
		PageInfo *PageInfo
	}
	type WrongConnection1 struct {
	}
	type WrongConnection2 struct {
		Edges    []string
		PageInfo *PageInfo
	}
	type WrongConnection11 struct {
		Edges    []*string
		PageInfo *PageInfo
	}
	type WrongConnection12 struct {
		Edges    string
		PageInfo *PageInfo
	}
	type WrongConnection3 struct {
		Edges    []PersonEdge
		PageInfo *PageInfo
	}
	type WrongConnection4 struct {
		Edges    []*PersonEdge
		PageInfo string
	}
	type WrongConnection5 struct {
		Edges    []*PersonEdge
		PageInfo PageInfo
	}
	type WrongConnection6 struct {
		Edges []*PersonEdge
	}
	type WrongConnection7 struct {
		Edges    []*WrongEdge1
		PageInfo *PageInfo
	}
	type WrongConnection8 struct {
		Edges    []*WrongEdge2
		PageInfo *PageInfo
	}
	type WrongConnection9 struct {
		Edges    []*WrongEdge3
		PageInfo *PageInfo
	}
	type WrongConnection10 struct {
		Edges    []*WrongEdge4
		PageInfo *PageInfo
	}
	type WrongConnection13 struct {
		Edges    []*WrongEdge5
		PageInfo *PageInfo
	}
	type args struct {
		result  interface{}
		list    interface{}
		pageOut *pagination.PageOutput
	}
	tests := []struct {
		name           string
		args           args
		wantErr        bool
		errorString    string
		expectedResult interface{}
		testResult     bool
	}{
		{
			name: "nil connection result",
			args: args{
				result:  nil,
				list:    []*Person{{Name: "fake1"}, {Name: "fake2"}},
				pageOut: &pagination.PageOutput{},
			},
			wantErr:     true,
			errorString: "connection result argument mustn't be nil",
		},
		{
			name: "nil input list",
			args: args{
				result:  &PersonConnection{},
				list:    nil,
				pageOut: &pagination.PageOutput{},
			},
			wantErr:     true,
			errorString: "list argument mustn't be nil",
		},
		{
			name: "wrong type as connection result",
			args: args{
				result:  "fake",
				list:    []*Person{{Name: "fake1"}, {Name: "fake2"}},
				pageOut: &pagination.PageOutput{},
			},
			wantErr:     true,
			errorString: "connection result argument must be a pointer to a connection object",
		},
		{
			name: "wrong type as input list",
			args: args{
				result:  &PersonConnection{},
				list:    "fake",
				pageOut: &pagination.PageOutput{},
			},
			wantErr:     true,
			errorString: "list argument must be a slice",
		},
		{
			name: "no field Edges in connection",
			args: args{
				result:  &WrongConnection1{},
				list:    []*Person{{Name: "fake1"}, {Name: "fake2"}},
				pageOut: &pagination.PageOutput{},
			},
			wantErr:     true,
			errorString: "field Edges not found in connection object",
		},
		{
			name: "Edges with a slice of wrong slice type ([]string)",
			args: args{
				result:  &WrongConnection2{},
				list:    []*Person{{Name: "fake1"}, {Name: "fake2"}},
				pageOut: &pagination.PageOutput{},
			},
			wantErr:     true,
			errorString: "field Edges must be a slice of struct pointers ([]*Edge)",
		},
		{
			name: "Edges with a slice of wrong slice type ([]*string)",
			args: args{
				result:  &WrongConnection11{},
				list:    []*Person{{Name: "fake1"}, {Name: "fake2"}},
				pageOut: &pagination.PageOutput{},
			},
			wantErr:     true,
			errorString: "field Edges must be a slice of struct pointers ([]*Edge)",
		},
		{
			name: "Edges with a slice of wrong slice type ([]PersonEdge)",
			args: args{
				result:  &WrongConnection3{},
				list:    []*Person{{Name: "fake1"}, {Name: "fake2"}},
				pageOut: &pagination.PageOutput{},
			},
			wantErr:     true,
			errorString: "field Edges must be a slice of struct pointers ([]*Edge)",
		},
		{
			name: "Edges with a slice of wrong slice type (string)",
			args: args{
				result:  &WrongConnection12{},
				list:    []*Person{{Name: "fake1"}, {Name: "fake2"}},
				pageOut: &pagination.PageOutput{},
			},
			wantErr:     true,
			errorString: "field Edges must be a slice",
		},
		{
			name: "PageInfo with wrong slice type (string)",
			args: args{
				result:  &WrongConnection4{},
				list:    []*Person{{Name: "fake1"}, {Name: "fake2"}},
				pageOut: &pagination.PageOutput{},
			},
			wantErr:     true,
			errorString: "field PageInfo isn't with the type *PageInfo",
		},
		{
			name: "PageInfo with wrong slice type (PageInfo)",
			args: args{
				result:  &WrongConnection5{},
				list:    []*Person{{Name: "fake1"}, {Name: "fake2"}},
				pageOut: &pagination.PageOutput{},
			},
			wantErr:     true,
			errorString: "field PageInfo isn't with the type *PageInfo",
		},
		{
			name: "PageInfo field not found",
			args: args{
				result:  &WrongConnection6{},
				list:    []*Person{{Name: "fake1"}, {Name: "fake2"}},
				pageOut: &pagination.PageOutput{},
			},
			wantErr:     true,
			errorString: "field PageInfo not found in connection object",
		},
		{
			name: "Edge without any needed key",
			args: args{
				result:  &WrongConnection7{},
				list:    []*Person{{Name: "fake1"}, {Name: "fake2"}},
				pageOut: &pagination.PageOutput{},
			},
			wantErr:     true,
			errorString: "field Cursor not found in Edge object",
		},
		{
			name: "Edge with cursor not a string",
			args: args{
				result:  &WrongConnection8{},
				list:    []*Person{{Name: "fake1"}, {Name: "fake2"}},
				pageOut: &pagination.PageOutput{},
			},
			wantErr:     true,
			errorString: "field Cursor from Edge object must be a string",
		},
		{
			name: "Edge with Node field not a valid structure",
			args: args{
				result:  &WrongConnection9{},
				list:    []*Person{{Name: "fake1"}, {Name: "fake2"}},
				pageOut: &pagination.PageOutput{},
			},
			wantErr:     true,
			errorString: "field Node must have the same type of items in the list argument",
		},
		{
			name: "Edge with Node field not a valid structure",
			args: args{
				result:  &WrongConnection10{},
				list:    []*Person{{Name: "fake1"}, {Name: "fake2"}},
				pageOut: &pagination.PageOutput{},
			},
			wantErr:     true,
			errorString: "field Node must have the same type of items in the list argument",
		},
		{
			name: "Edge without Node",
			args: args{
				result:  &WrongConnection13{},
				list:    []*Person{{Name: "fake1"}, {Name: "fake2"}},
				pageOut: &pagination.PageOutput{},
			},
			wantErr:     true,
			errorString: "field Node not found in Edge object",
		},
		{
			name: "empty list array",
			args: args{
				result:  &PersonConnection{},
				list:    []*Person{},
				pageOut: &pagination.PageOutput{},
			},
			expectedResult: &PersonConnection{
				Edges: nil,
				PageInfo: &PageInfo{
					HasNextPage:     false,
					HasPreviousPage: false,
					EndCursor:       nil,
					StartCursor:     nil,
				},
			},
			testResult: true,
		},
		{
			name: "2 elements in list array without any previous and next page",
			args: args{
				result:  &PersonConnection{},
				list:    []*Person{{Name: "fake1"}, {Name: "fake2"}},
				pageOut: &pagination.PageOutput{},
			},
			expectedResult: &PersonConnection{
				Edges: []*PersonEdge{
					{
						Cursor: "cGFnaW5hdGU6MQ==",
						Node:   &Person{Name: "fake1"},
					},
					{
						Cursor: "cGFnaW5hdGU6Mg==",
						Node:   &Person{Name: "fake2"},
					},
				},
				PageInfo: &PageInfo{
					HasNextPage:     false,
					HasPreviousPage: false,
					StartCursor:     starString("cGFnaW5hdGU6MQ=="),
					EndCursor:       starString("cGFnaW5hdGU6Mg=="),
				},
			},
			testResult: true,
		},
		{
			name: "2 elements in list array with previous and next page",
			args: args{
				result:  &PersonConnection{},
				list:    []*Person{{Name: "fake1"}, {Name: "fake2"}},
				pageOut: &pagination.PageOutput{HasNext: true, HasPrevious: true},
			},
			expectedResult: &PersonConnection{
				Edges: []*PersonEdge{
					{
						Cursor: "cGFnaW5hdGU6MQ==",
						Node:   &Person{Name: "fake1"},
					},
					{
						Cursor: "cGFnaW5hdGU6Mg==",
						Node:   &Person{Name: "fake2"},
					},
				},
				PageInfo: &PageInfo{
					HasNextPage:     true,
					HasPreviousPage: true,
					StartCursor:     starString("cGFnaW5hdGU6MQ=="),
					EndCursor:       starString("cGFnaW5hdGU6Mg=="),
				},
			},
			testResult: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := MapConnection(tt.args.result, tt.args.list, tt.args.pageOut)
			if (err != nil) != tt.wantErr {
				t.Errorf("MapConnection() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && err.Error() != tt.errorString {
				t.Errorf("MapConnection() error = %v, wantErr %v", err, tt.errorString)
				return
			}
			if tt.testResult {
				assert.Equal(t, tt.expectedResult, tt.args.result)
			}
		})
	}
}
