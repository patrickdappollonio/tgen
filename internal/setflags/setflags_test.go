package setflags

import (
	"testing"
)

func TestParseSetValues(t *testing.T) {
	tests := []struct {
		name      string
		setValues []string
		expected  map[string]any
		wantErr   bool
	}{
		{
			name:      "simple key-value",
			setValues: []string{"key=value"},
			expected: map[string]any{
				"key": "value",
			},
		},
		{
			name:      "nested key-value",
			setValues: []string{"key.subkey=value"},
			expected: map[string]any{
				"key": map[string]any{
					"subkey": "value",
				},
			},
		},
		{
			name:      "boolean values",
			setValues: []string{"debug=true", "enabled=false"},
			expected: map[string]any{
				"debug":   true,
				"enabled": false,
			},
		},
		{
			name:      "YAML boolean values",
			setValues: []string{"debug=yes", "enabled=no", "active=on", "inactive=off"},
			expected: map[string]any{
				"debug":    true,
				"enabled":  false,
				"active":   true,
				"inactive": false,
			},
		},
		{
			name:      "integer values",
			setValues: []string{"port=8080", "replicas=3"},
			expected: map[string]any{
				"port":     8080,
				"replicas": 3,
			},
		},
		{
			name:      "float values",
			setValues: []string{"ratio=0.5", "version=1.2"},
			expected: map[string]any{
				"ratio":   0.5,
				"version": 1.2,
			},
		},
		{
			name:      "comma-separated values",
			setValues: []string{"a=b,c=d"},
			expected: map[string]any{
				"a": "b",
				"c": "d",
			},
		},
		{
			name:      "array syntax",
			setValues: []string{"tags={frontend,backend,database}"},
			expected: map[string]any{
				"tags": []any{"frontend", "backend", "database"},
			},
		},
		{
			name:      "array indexing",
			setValues: []string{"servers[0].port=80"},
			expected: map[string]any{
				"servers": []any{
					map[string]any{
						"port": 80,
					},
				},
			},
		},
		{
			name:      "null value",
			setValues: []string{"nullValue=null"},
			expected: map[string]any{
				"nullValue": nil,
			},
		},
		{
			name:      "empty array",
			setValues: []string{"emptyArray=[]"},
			expected: map[string]any{
				"emptyArray": []any{},
			},
		},
		{
			name:      "escaped comma",
			setValues: []string{"name=value1\\,value2"},
			expected: map[string]any{
				"name": "value1,value2",
			},
		},
		{
			name:      "invalid format - no equals",
			setValues: []string{"keynoequals"},
			wantErr:   true,
		},
		{
			name:      "invalid format - empty key",
			setValues: []string{"=value"},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseSetValues(tt.setValues)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseSetValues() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("ParseSetValues() unexpected error: %v", err)
				return
			}

			if !mapsEqual(result, tt.expected) {
				t.Errorf("ParseSetValues() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestParseSetStringValues(t *testing.T) {
	tests := []struct {
		name      string
		setValues []string
		expected  map[string]any
		wantErr   bool
	}{
		{
			name:      "simple key-value as string",
			setValues: []string{"key=value"},
			expected: map[string]any{
				"key": "value",
			},
		},
		{
			name:      "boolean value as string",
			setValues: []string{"debug=true"},
			expected: map[string]any{
				"debug": "true",
			},
		},
		{
			name:      "number value as string",
			setValues: []string{"port=8080"},
			expected: map[string]any{
				"port": "8080",
			},
		},
		{
			name:      "YAML boolean values as strings",
			setValues: []string{"enabled=yes", "disabled=no"},
			expected: map[string]any{
				"enabled":  "yes",
				"disabled": "no",
			},
		},
		{
			name:      "array syntax as strings",
			setValues: []string{"tags={frontend,backend,database}"},
			expected: map[string]any{
				"tags": []any{"frontend", "backend", "database"},
			},
		},
		{
			name:      "mixed types forced as strings",
			setValues: []string{"values={hello,123,true,no}"},
			expected: map[string]any{
				"values": []any{"hello", "123", "true", "no"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseSetStringValues(tt.setValues)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseSetStringValues() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("ParseSetStringValues() unexpected error: %v", err)
				return
			}

			if !mapsEqual(result, tt.expected) {
				t.Errorf("ParseSetStringValues() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestParseValueWithTypeInference(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected any
	}{
		{
			name:     "boolean true",
			value:    "true",
			expected: true,
		},
		{
			name:     "boolean false",
			value:    "false",
			expected: false,
		},
		{
			name:     "YAML boolean yes",
			value:    "yes",
			expected: true,
		},
		{
			name:     "YAML boolean no",
			value:    "no",
			expected: false,
		},
		{
			name:     "integer",
			value:    "123",
			expected: 123,
		},
		{
			name:     "float",
			value:    "3.14",
			expected: 3.14,
		},
		{
			name:     "float with trailing zero",
			value:    "1.0",
			expected: 1.0,
		},
		{
			name:     "string",
			value:    "hello",
			expected: "hello",
		},
		{
			name:     "country code no (should be boolean)",
			value:    "no",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseValueWithTypeInference(tt.value)

			if !valuesEqual(result, tt.expected) {
				t.Errorf("parseValueWithTypeInference() = %v (type %T), want %v (type %T)", result, result, tt.expected, tt.expected)
			}
		})
	}
}

// Helper functions for testing
func mapsEqual(a, b map[string]any) bool {
	if len(a) != len(b) {
		return false
	}

	for k, v := range a {
		if bv, ok := b[k]; !ok || !valuesEqual(v, bv) {
			return false
		}
	}

	return true
}

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
