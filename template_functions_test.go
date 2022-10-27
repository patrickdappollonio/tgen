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

func Test_shuffle(t *testing.T) {
	tests := []struct {
		name    string
		seq     any
		wantErr bool
	}{
		{
			name: "shuffle 1 2 3",
			seq:  []int{1, 2, 3},
		},
		{
			name:    "nil",
			seq:     nil,
			wantErr: true,
		},
		{
			name:    "empty",
			seq:     []int{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := shuffle(tt.seq)
			if (err != nil) != tt.wantErr {
				t.Errorf("shuffle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			var l1, l2 int

			if got != nil {
				l1 = reflect.ValueOf(got).Len()
			}

			if tt.seq != nil {
				l2 = reflect.ValueOf(tt.seq).Len()
			}

			if !tt.wantErr && l1 != l2 {
				t.Errorf("shuffle() got length = %d (original: %d)", l1, l2)
			}
		})
	}
}

func Test_first(t *testing.T) {
	tests := []struct {
		name    string
		seq     any
		want    any
		wantErr bool
	}{
		{
			name: "first 1 2 3",
			seq:  []int{1, 2, 3},
			want: 1,
		},
		{
			name:    "nil",
			seq:     nil,
			wantErr: true,
		},
		{
			name:    "empty",
			seq:     []int{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := first(tt.seq)
			if (err != nil) != tt.wantErr {
				t.Errorf("first() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("first() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_last(t *testing.T) {
	tests := []struct {
		name    string
		seq     any
		want    any
		wantErr bool
	}{
		{
			name: "last 1 2 3",
			seq:  []int{1, 2, 3},
			want: 3,
		},
		{
			name:    "nil",
			seq:     nil,
			wantErr: true,
		},
		{
			name:    "empty",
			seq:     []int{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := last(tt.seq)
			if (err != nil) != tt.wantErr {
				t.Errorf("last() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("last() = %v, want %v", got, tt.want)
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

func Test_replace(t *testing.T) {
	got := replace("hello", "hello", "world")

	if got != "world" {
		t.Errorf("replace() = %v, want %v", got, "world")
	}
}

func Test_sha256sum(t *testing.T) {
	got := sha256sum("hello")

	if got != "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824" {
		t.Errorf("sha256sum() = %v, want %v", got, "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824")
	}
}

func Test_sha1sum(t *testing.T) {
	got := sha1sum("hello")

	if got != "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d" {
		t.Errorf("sha1sum() = %v, want %v", got, "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d")
	}
}

func Test_base64encode(t *testing.T) {
	got := base64encode("hello")

	if got != "aGVsbG8=" {
		t.Errorf("base64encode() = %v, want %v", got, "aGVsbG8=")
	}
}

func Test_base64decode(t *testing.T) {
	got, err := base64decode("aGVsbG8=")

	if err != nil {
		t.Errorf("base64decode() error = %v", err)
	}

	if got != "hello" {
		t.Errorf("base64decode() = %v, want %v", got, "hello")
	}

	_, err = base64decode("987sd6f98sydf^%!!@!@#")

	if err == nil {
		t.Errorf("base64decode() want error, got none")
	}
}

func Test_nospace(t *testing.T) {
	got := nospace("hello world")

	if got != "helloworld" {
		t.Errorf("nospace() = %v, want %v", got, "helloworld")
	}
}

func Test_repeat(t *testing.T) {
	got := repeat(3, "hello")

	if got != "hellohellohello" {
		t.Errorf("repeat() = %v, want %v", got, "hellohellohello")
	}
}

func Test_indent(t *testing.T) {
	got := indent(3, "hello")

	if got != "   hello" {
		t.Errorf("indent() = %v, want %v", got, "   hello")
	}
}

func Test_nindent(t *testing.T) {
	got := nindent(3, "hello")

	if got != "\n   hello" {
		t.Errorf("nindent() = %v, want %v", got, "\n   hello")
	}
}

func Test_quote(t *testing.T) {
	got := quote("hello")

	if got != "\"hello\"" {
		t.Errorf("quote() = %v, want %v", got, "\"hello\"")
	}
}

func Test_squote(t *testing.T) {
	got := squote("hello")

	if got != "'hello'" {
		t.Errorf("squote() = %v, want %v", got, "'hello'")
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
