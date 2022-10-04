package main

import (
	"reflect"
	"testing"
)

func Test_seq(t *testing.T) {
	tests := []struct {
		name    string
		args    []int
		want    []int
		wantErr bool
	}{
		{
			name: "seq 10",
			args: []int{10},
			want: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
		{
			name: "seq 1 5",
			args: []int{1, 5},
			want: []int{1, 2, 3, 4, 5},
		},
		{
			name: "seq 5 -1 1",
			args: []int{5, -1, 1},
			want: []int{5, 4, 3, 2, 1},
		},
		{
			name: "seq 1 2 5",
			args: []int{1, 2, 5},
			want: []int{1, 3, 5},
		},
		{
			name: "seq 5 -2 1",
			args: []int{5, -2, 1},
			want: []int{5, 3, 1},
		},
		{
			name:    "seq 1 0 5",
			args:    []int{1, 0, 5},
			wantErr: true, // step cannot be zero
		},
		{
			name:    "seq 1 2 5 7",
			args:    []int{1, 2, 5, 7},
			wantErr: true, // incorrect number of arguments
		},
		{
			name:    "seq 100 2 1",
			args:    []int{100, 2, 1},
			wantErr: true, // increment must be < 0
		},
		{
			name:    "seq 1 -2 100",
			args:    []int{1, -2, 100},
			wantErr: true, // increment must be > 0
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := seq(tt.args...)

			if (err != nil) != tt.wantErr {
				t.Errorf("seq() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("seq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_slice(t *testing.T) {
	tests := []struct {
		name string
		args []interface{}
		want []interface{}
	}{
		{
			name: "slice 1 2 3",
			args: []interface{}{1, 2, 3},
			want: []interface{}{1, 2, 3},
		},
		{
			name: "slice 1 2 3 a b c",
			args: []interface{}{1, 2, 3, "a", "b", "c"},
			want: []interface{}{1, 2, 3, "a", "b", "c"},
		},
		{
			name: "slice 1 2 3 a b c false true interface{} struct{} nil",
			args: []interface{}{1, 2, 3, "a", "b", "c", 4, 5, 6, false, true, interface{}(nil), struct{}{}, nil},
			want: []interface{}{1, 2, 3, "a", "b", "c", 4, 5, 6, false, true, interface{}(nil), struct{}{}, nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := slice(tt.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("slice() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
