//+build unit

package authorization

import (
	"reflect"
	"testing"
)

func Test_deleteEmpty(t *testing.T) {
	type args struct {
		s []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Clean",
			args: args{s: []string{"", "f", ""}},
			want: []string{"f"},
		},
		{
			name: "No Clean",
			args: args{s: []string{"f"}},
			want: []string{"f"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := deleteEmpty(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("deleteEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}
