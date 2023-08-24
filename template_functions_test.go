package main

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

func Test_raw(t *testing.T) {
	got := raw("hello")

	if got != "hello" {
		t.Errorf("raw() = %v, want %v", got, "hello")
	}
}

func Test_linebyline(t *testing.T) {
	lines := linebyline("hello\nworld")

	if len(lines) != 2 {
		t.Errorf("linebyline() = %v, want %v", len(lines), 2)
	}

	if lines[0] != "hello" {
		t.Errorf("linebyline() = %v, want %v", lines[0], "hello")
	}

	if lines[1] != "world" {
		t.Errorf("linebyline() = %v, want %v", lines[1], "world")
	}
}

func Test_envstrict(t *testing.T) {
	tests := []struct {
		name    string
		kv      map[string]string
		strict  bool
		key     string
		want    string
		wantErr bool
	}{
		{
			name:   "strict",
			kv:     map[string]string{"FOO": "bar"},
			strict: true,
			key:    "FOO",
			want:   "bar",
		},

		{
			name:   "not strict",
			kv:     map[string]string{"FOO": "bar"},
			strict: false,
			key:    "FOO",
			want:   "bar",
		},
		{
			name:    "strict missing key",
			kv:      map[string]string{"FOO": "bar"},
			strict:  true,
			key:     "BAR",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := envstrict(tt.kv, tt.strict)(tt.key)

			if (err != nil) != tt.wantErr {
				t.Errorf("envstrict() error = %v, wantErr %v", err, tt.wantErr)
			}

			if got != tt.want {
				t.Errorf("envstrict() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_envdefault(t *testing.T) {
	tests := []struct {
		name    string
		kv      map[string]string
		key     string
		defval  string
		want    string
		wantErr bool
	}{
		{
			name: "read existent",
			kv:   map[string]string{"FOO": "bar"},
			key:  "FOO",
			want: "bar",
		},
		{
			name:   "read non-existent with defval",
			kv:     map[string]string{"FOO": "bar"},
			key:    "BAR",
			defval: "baz",
			want:   "baz",
		},
		{
			name: "read non-existent without defval",
			kv:   map[string]string{"FOO": "bar"},
			key:  "BAR",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := envdefault(tt.kv)(tt.key, tt.defval)

			if (err != nil) != tt.wantErr {
				t.Errorf("envdefault() error = %v, wantErr %v", err, tt.wantErr)
			}

			if got != tt.want {
				t.Errorf("envdefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_rndgen(t *testing.T) {
	got1 := rndgen(5)
	got2 := rndgen(5)

	if got1 == got2 {
		t.Errorf("rndgen() values must be different, got = %q, %q", got1, got2)
	}
}

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
