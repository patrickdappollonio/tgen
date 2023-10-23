package tfuncs

import (
	"reflect"
	"testing"
)

func Test_requiredField(t *testing.T) {
	tests := []struct {
		name    string
		warn    string
		val     interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name:    "required field nullable",
			warn:    "required field must be set",
			val:     nil,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "required field not nullable",
			warn:    "required field must be set",
			val:     "",
			want:    "",
			wantErr: true,
		},
		{
			name: "required field set",
			val:  "hello",
			want: "hello",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := requiredField(tt.warn, tt.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("requiredField() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("requiredField() = %v, want %v", got, tt.want)
			}
		})
	}
}
