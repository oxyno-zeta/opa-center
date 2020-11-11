// +build unit

package utils

import (
	"encoding/base64"
	"reflect"
	"testing"

	"github.com/oxyno-zeta/opa-center/pkg/opa-center/database/pagination"
)

func Test_ToIDRelay(t *testing.T) {
	type args struct {
		prefix string
		id     string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty",
			args: args{
				id:     "",
				prefix: "",
			},
			want: base64.StdEncoding.EncodeToString([]byte(":")),
		},
		{
			name: "valid",
			args: args{id: "id", prefix: "prefix"},
			want: base64.StdEncoding.EncodeToString([]byte("prefix:id")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToIDRelay(tt.args.prefix, tt.args.id); got != tt.want {
				t.Errorf("ToIDRelay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPageInput(t *testing.T) {
	toStarString := func(s string) *string { return &s }
	toStarInt := func(i int) *int { return &i }
	type args struct {
		after  *string
		before *string
		first  *int
		last   *int
	}
	tests := []struct {
		name        string
		args        args
		want        *pagination.PageInput
		wantErr     bool
		errorString string
	}{
		{
			name:    "empty",
			args:    args{},
			wantErr: false,
			want: &pagination.PageInput{
				Limit: 10,
				Skip:  0,
			},
		},
		{
			name: "after and before error",
			args: args{
				before: toStarString("fake"),
				after:  toStarString("fake"),
			},
			wantErr:     true,
			errorString: "after and before can't be present together at the same time",
		},
		{
			name: "first and last error",
			args: args{
				first: toStarInt(1),
				last:  toStarInt(1),
			},
			wantErr:     true,
			errorString: "first and last can't be present together at the same time",
		},
		{
			name: "before and last error",
			args: args{
				before: toStarString("fake"),
			},
			wantErr:     true,
			errorString: "before must be used with last element",
		},
		{
			name: "before and last error 2",
			args: args{
				last: toStarInt(1),
			},
			wantErr:     true,
			errorString: "last must be used with before element",
		},
		{
			name: "first and after error",
			args: args{
				after: toStarString("fake"),
			},
			wantErr:     true,
			errorString: "first must be used with after element",
		},
		{
			name: "last > 0 error",
			args: args{
				before: toStarString("fake"),
				last:   toStarInt(-1),
			},
			wantErr:     true,
			errorString: "last must be > 0",
		},
		{
			name: "last > 0 error 2",
			args: args{
				before: toStarString("fake"),
				last:   toStarInt(0),
			},
			wantErr:     true,
			errorString: "last must be > 0",
		},
		{
			name: "first > 0 error",
			args: args{
				after: toStarString("fake"),
				first: toStarInt(-1),
			},
			wantErr:     true,
			errorString: "first must be > 0",
		},
		{
			name: "first > 0 error 2",
			args: args{
				after: toStarString("fake"),
				first: toStarInt(0),
			},
			wantErr:     true,
			errorString: "first must be > 0",
		},
		{
			name: "before case with negative substract",
			args: args{
				before: toStarString(base64.StdEncoding.EncodeToString([]byte("paginate:1"))),
				last:   toStarInt(2),
			},
			want: &pagination.PageInput{
				Skip:  0,
				Limit: 2,
			},
		},
		{
			name: "before case",
			args: args{
				before: toStarString(base64.StdEncoding.EncodeToString([]byte("paginate:10"))),
				last:   toStarInt(2),
			},
			want: &pagination.PageInput{
				Skip:  7,
				Limit: 2,
			},
		},
		{
			name: "before case with too big limit",
			args: args{
				before: toStarString(base64.StdEncoding.EncodeToString([]byte("paginate:10"))),
				last:   toStarInt(200),
			},
			wantErr:     true,
			errorString: "first or last is too big, maximum is 50",
		},
		{
			name: "after case",
			args: args{
				after: toStarString(base64.StdEncoding.EncodeToString([]byte("paginate:10"))),
				first: toStarInt(2),
			},
			want: &pagination.PageInput{
				Skip:  10,
				Limit: 2,
			},
		},
		{
			name: "after case 2",
			args: args{
				after: toStarString(base64.StdEncoding.EncodeToString([]byte("paginate:1"))),
				first: toStarInt(2),
			},
			want: &pagination.PageInput{
				Skip:  1,
				Limit: 2,
			},
		},
		{
			name: "after case with too big limit",
			args: args{
				after: toStarString(base64.StdEncoding.EncodeToString([]byte("paginate:10"))),
				first: toStarInt(200),
			},
			wantErr:     true,
			errorString: "first or last is too big, maximum is 50",
		},
		{
			name: "only first",
			args: args{
				first: toStarInt(2),
			},
			want: &pagination.PageInput{
				Skip:  0,
				Limit: 2,
			},
		},
		{
			name: "after is an empty string and first is present",
			args: args{
				after: toStarString(""),
				first: toStarInt(2),
			},
			want: &pagination.PageInput{
				Skip:  0,
				Limit: 2,
			},
		},
		{
			name: "before is an empty string",
			args: args{
				before: toStarString(""),
				last:   toStarInt(2),
			},
			wantErr:     true,
			errorString: "last must be used with before element",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetPageInput(tt.args.after, tt.args.before, tt.args.first, tt.args.last)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPageInput() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && err.Error() != tt.errorString {
				t.Errorf("GetPageInput() error = %v, wantErr %v", err, tt.errorString)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPageInput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parsePaginateCursor(t *testing.T) {
	type args struct {
		cursorB64 string
	}
	tests := []struct {
		name        string
		args        args
		want        int
		wantErr     bool
		errorString string
	}{
		{
			name: "invalid base64",
			args: args{
				cursorB64: "fake:fake",
			},
			wantErr:     true,
			errorString: "illegal base64 data at input byte 4",
		},
		{
			name: "invalid format",
			args: args{
				cursorB64: "fake",
			},
			wantErr:     true,
			errorString: "format error on relay token",
		},
		{
			name: "invalid prefix",
			args: args{
				cursorB64: base64.StdEncoding.EncodeToString([]byte("fake:1")),
			},
			wantErr:     true,
			errorString: "invalid relay prefix",
		},
		{
			name: "invalid number",
			args: args{
				cursorB64: base64.StdEncoding.EncodeToString([]byte("paginate:fake")),
			},
			wantErr:     true,
			errorString: "strconv.Atoi: parsing \"fake\": invalid syntax",
		},
		{
			name: "cursor must be positive",
			args: args{
				cursorB64: base64.StdEncoding.EncodeToString([]byte("paginate:-1")),
			},
			wantErr:     true,
			errorString: "cursor pagination must be positive",
		},
		{
			name: "valid",
			args: args{
				cursorB64: base64.StdEncoding.EncodeToString([]byte("paginate:1")),
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePaginateCursor(tt.args.cursorB64)
			if (err != nil) != tt.wantErr {
				t.Errorf("parsePaginateCursor() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && err.Error() != tt.errorString {
				t.Errorf("parsePaginateCursor() error = %v, wantErr %v", err, tt.errorString)
				return
			}
			if got != tt.want {
				t.Errorf("parsePaginateCursor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_GetPaginateCursor(t *testing.T) {
	type args struct {
		index int
		skip  int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty",
			args: args{
				index: 0,
				skip:  0,
			},
			want: base64.StdEncoding.EncodeToString([]byte("paginate:1")),
		},
		{
			name: "skip",
			args: args{
				index: 0,
				skip:  5,
			},
			want: base64.StdEncoding.EncodeToString([]byte("paginate:6")),
		},
		{
			name: "index",
			args: args{
				index: 5,
				skip:  0,
			},
			want: base64.StdEncoding.EncodeToString([]byte("paginate:6")),
		},
		{
			name: "index and skip",
			args: args{
				index: 5,
				skip:  5,
			},
			want: base64.StdEncoding.EncodeToString([]byte("paginate:11")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetPaginateCursor(tt.args.index, tt.args.skip); got != tt.want {
				t.Errorf("GetPaginateCursor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_GetPageInfo(t *testing.T) {
	toStarString := func(s string) *string { return &s }
	type args struct {
		startCursor string
		endCursor   string
		p           *pagination.PageOutput
	}
	tests := []struct {
		name string
		args args
		want *PageInfo
	}{
		{
			name: "empty",
			args: args{},
			want: &PageInfo{},
		},
		{
			name: "start cursor exists",
			args: args{
				startCursor: "start",
			},
			want: &PageInfo{
				StartCursor: toStarString("start"),
			},
		},
		{
			name: "end cursor exists",
			args: args{
				endCursor: "end",
			},
			want: &PageInfo{
				EndCursor: toStarString("end"),
			},
		},
		{
			name: "paginator exists",
			args: args{
				p: &pagination.PageOutput{
					HasNext:     true,
					HasPrevious: true,
				},
			},
			want: &PageInfo{
				HasNextPage:     true,
				HasPreviousPage: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetPageInfo(tt.args.startCursor, tt.args.endCursor, tt.args.p); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPageInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
