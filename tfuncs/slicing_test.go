package tfuncs

import (
	"reflect"
	"testing"
)

func Test_after(t *testing.T) {
	type args struct {
		after  int
		values interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "after 1 [1 2 3]",
			args: args{
				after:  1,
				values: []int{1, 2, 3},
			},
			want: []int{2, 3},
		},
		{
			name: "after 2 [1 2 3]",
			args: args{
				after:  2,
				values: []interface{}{1, 2, 3},
			},
			want: []interface{}{3},
		},
		{
			name: "after 3 [1 2 3]",
			args: args{
				after:  3,
				values: []interface{}{1, 2, 3},
			},
			want: []interface{}{},
		},
		{
			name: "after -1 [1 2 3]",
			args: args{
				after:  -1,
				values: []interface{}{1, 2, 3},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := after(tt.args.after, tt.args.values)

			if (err != nil) != tt.wantErr {
				t.Errorf("after() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("after() = %v (%T), want %v (%T)", got, got, tt.want, tt.want)
			}
		})
	}
}
