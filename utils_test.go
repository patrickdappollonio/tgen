package main

import "testing"

func Test_parseEnvLine(t *testing.T) {
	tests := []struct {
		name      string
		line      string
		wantKey   string
		wantValue string
		wantErr   bool
	}{
		{
			name:      "empty",
			line:      "",
			wantKey:   "",
			wantValue: "",
		},
		{
			name:      "comment",
			line:      "# comment",
			wantKey:   "",
			wantValue: "",
		},
		{
			name:      "comment with space",
			line:      " # comment",
			wantKey:   "",
			wantValue: "",
		},
		{
			name:      "key lowercase and value",
			line:      "key=value",
			wantKey:   "KEY",
			wantValue: "value",
		},
		{
			name:      "key uppercase and value",
			line:      "KEY=value",
			wantKey:   "KEY",
			wantValue: "value",
		},
		{
			name:      "multi-equals",
			line:      "KEY=value=1",
			wantKey:   "KEY",
			wantValue: "value=1",
		},
		{
			name:    "no key",
			line:    "=value1",
			wantErr: true,
		},
		{
			name:    "no separator",
			line:    "KEYvalue1",
			wantErr: true,
		},
		{
			name:      "no value",
			line:      "KEY=",
			wantKey:   "KEY",
			wantValue: "",
		},
		{
			name:      "no value with space",
			line:      "KEY= ",
			wantKey:   "KEY",
			wantValue: "",
		},
		{
			name:      "quoted value",
			line:      "KEY=\"value1\"",
			wantKey:   "KEY",
			wantValue: "value1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKey, gotValue, err := parseEnvLine(tt.line)

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotKey != tt.wantKey {
				t.Errorf("got key = %q, want %q", gotKey, tt.wantKey)
			}

			if gotValue != tt.wantValue {
				t.Errorf("got value = %q, want %q", gotValue, tt.wantValue)
			}
		})
	}
}
