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

func TestMergeMap(t *testing.T) {
	tests := []struct {
		name     string
		dest     map[string]any
		src      map[string]any
		expected map[string]any
	}{
		{
			name: "merge simple maps",
			dest: map[string]any{
				"key1": "value1",
			},
			src: map[string]any{
				"key2": "value2",
			},
			expected: map[string]any{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			name: "override existing key",
			dest: map[string]any{
				"key": "oldvalue",
			},
			src: map[string]any{
				"key": "newvalue",
			},
			expected: map[string]any{
				"key": "newvalue",
			},
		},
		{
			name: "merge nested maps",
			dest: map[string]any{
				"config": map[string]any{
					"existing": "value",
				},
			},
			src: map[string]any{
				"config": map[string]any{
					"new": "value",
				},
			},
			expected: map[string]any{
				"config": map[string]any{
					"existing": "value",
					"new":      "value",
				},
			},
		},
		{
			name: "override nested value",
			dest: map[string]any{
				"config": map[string]any{
					"key": "oldvalue",
				},
			},
			src: map[string]any{
				"config": map[string]any{
					"key": "newvalue",
				},
			},
			expected: map[string]any{
				"config": map[string]any{
					"key": "newvalue",
				},
			},
		},
		{
			name: "replace non-map with map",
			dest: map[string]any{
				"key": "stringvalue",
			},
			src: map[string]any{
				"key": map[string]any{
					"nested": "value",
				},
			},
			expected: map[string]any{
				"key": map[string]any{
					"nested": "value",
				},
			},
		},
		{
			name: "nil dest map",
			dest: nil,
			src: map[string]any{
				"key": "value",
			},
			expected: map[string]any{
				"key": "value",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mergeMap(tt.dest, tt.src)

			if !mapsEqual(result, tt.expected) {
				t.Errorf("mergeMap() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// mapsEqual compares two maps for equality, handling nested maps
func mapsEqual(a, b map[string]any) bool {
	if len(a) != len(b) {
		return false
	}

	for k, v := range a {
		bv, exists := b[k]
		if !exists {
			return false
		}

		if !valuesEqual(v, bv) {
			return false
		}
	}

	return true
}

// valuesEqual compares two values for equality, handling nested maps and slices
func valuesEqual(a, b any) bool {
	switch av := a.(type) {
	case map[string]any:
		if bv, ok := b.(map[string]any); ok {
			return mapsEqual(av, bv)
		}
		return false
	case []any:
		if bv, ok := b.([]any); ok {
			return slicesEqual(av, bv)
		}
		return false
	case nil:
		return b == nil
	default:
		return a == b
	}
}

// slicesEqual compares two slices for equality
func slicesEqual(a, b []any) bool {
	if len(a) != len(b) {
		return false
	}

	for i, av := range a {
		if !valuesEqual(av, b[i]) {
			return false
		}
	}

	return true
}
