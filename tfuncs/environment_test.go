package tfuncs

import "testing"

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
