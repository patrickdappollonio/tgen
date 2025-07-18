package main

import (
	"reflect"
	"testing"

	"github.com/patrickdappollonio/tgen/internal/setflags"
)

func TestMergeSetValues(t *testing.T) {
	tg := &tgen{}

	tests := []struct {
		name         string
		existingYAML map[string]any
		setValues    []string
		expected     map[string]any
		wantErr      bool
	}{
		{
			name:         "merge with empty existing values",
			existingYAML: nil,
			setValues:    []string{"key=value"},
			expected: map[string]any{
				"key": "value",
				"Values": map[string]any{
					"key": "value",
				},
			},
		},
		{
			name: "merge with existing values",
			existingYAML: map[string]any{
				"existing": "value",
				"Values": map[string]any{
					"existing": "value",
				},
			},
			setValues: []string{"new=value"},
			expected: map[string]any{
				"existing": "value",
				"new":      "value",
				"Values": map[string]any{
					"existing": "value",
					"new":      "value",
				},
			},
		},
		{
			name: "override existing values",
			existingYAML: map[string]any{
				"key": "oldvalue",
				"Values": map[string]any{
					"key": "oldvalue",
				},
			},
			setValues: []string{"key=newvalue"},
			expected: map[string]any{
				"key": "newvalue",
				"Values": map[string]any{
					"key": "newvalue",
				},
			},
		},
		{
			name: "merge nested structures",
			existingYAML: map[string]any{
				"app": map[string]any{
					"name": "myapp",
				},
				"Values": map[string]any{
					"app": map[string]any{
						"name": "myapp",
					},
				},
			},
			setValues: []string{"app.version=1.0"},
			expected: map[string]any{
				"app": map[string]any{
					"name":    "myapp",
					"version": 1.0, // Should be float
				},
				"Values": map[string]any{
					"app": map[string]any{
						"name":    "myapp",
						"version": 1.0, // Should be float
					},
				},
			},
		},
		{
			name: "YAML boolean edge case - country code as boolean",
			existingYAML: map[string]any{
				"user": map[string]any{
					"name": "John",
				},
				"Values": map[string]any{
					"user": map[string]any{
						"name": "John",
					},
				},
			},
			setValues: []string{"user.country=no"},
			expected: map[string]any{
				"user": map[string]any{
					"name":    "John",
					"country": false, // "no" should be parsed as boolean
				},
				"Values": map[string]any{
					"user": map[string]any{
						"name":    "John",
						"country": false,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tg.yamlValues = tt.existingYAML
			err := tg.mergeSetValues(tt.setValues)

			if tt.wantErr {
				if err == nil {
					t.Errorf("mergeSetValues() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("mergeSetValues() unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(tg.yamlValues, tt.expected) {
				t.Errorf("mergeSetValues() = %v, want %v", tg.yamlValues, tt.expected)
			}
		})
	}
}

func TestMergeSetStringValues(t *testing.T) {
	tg := &tgen{}

	tests := []struct {
		name         string
		existingYAML map[string]any
		setValues    []string
		expected     map[string]any
		wantErr      bool
	}{
		{
			name:         "merge with empty existing values",
			existingYAML: nil,
			setValues:    []string{"key=value"},
			expected: map[string]any{
				"key": "value",
				"Values": map[string]any{
					"key": "value",
				},
			},
		},
		{
			name: "override existing typed values with strings",
			existingYAML: map[string]any{
				"port":  8080,
				"debug": true,
				"Values": map[string]any{
					"port":  8080,
					"debug": true,
				},
			},
			setValues: []string{"port=9000", "debug=false"},
			expected: map[string]any{
				"port":  "9000",  // Should be string
				"debug": "false", // Should be string
				"Values": map[string]any{
					"port":  "9000",
					"debug": "false",
				},
			},
		},
		{
			name: "force boolean as string",
			existingYAML: map[string]any{
				"enabled": true,
				"Values": map[string]any{
					"enabled": true,
				},
			},
			setValues: []string{"enabled=true"},
			expected: map[string]any{
				"enabled": "true", // Should be string
				"Values": map[string]any{
					"enabled": "true",
				},
			},
		},
		{
			name: "YAML boolean edge case - country code as string",
			existingYAML: map[string]any{
				"user": map[string]any{
					"name": "John",
				},
				"Values": map[string]any{
					"user": map[string]any{
						"name": "John",
					},
				},
			},
			setValues: []string{"user.country=no"},
			expected: map[string]any{
				"user": map[string]any{
					"name":    "John",
					"country": "no", // Should be string with --set-string
				},
				"Values": map[string]any{
					"user": map[string]any{
						"name":    "John",
						"country": "no",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tg.yamlValues = tt.existingYAML
			err := tg.mergeSetStringValues(tt.setValues)

			if tt.wantErr {
				if err == nil {
					t.Errorf("mergeSetStringValues() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("mergeSetStringValues() unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(tg.yamlValues, tt.expected) {
				t.Errorf("mergeSetStringValues() = %v, want %v", tg.yamlValues, tt.expected)
			}
		})
	}
}

// Test that the tgen methods properly delegate to the internal package
func TestTgenParseSetValues(t *testing.T) {
	tg := &tgen{}
	result, err := tg.parseSetValues([]string{"key=value"})
	if err != nil {
		t.Errorf("parseSetValues() unexpected error: %v", err)
		return
	}

	expected := map[string]any{"key": "value"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("parseSetValues() = %v, want %v", result, expected)
	}
}

func TestTgenParseSetStringValues(t *testing.T) {
	tg := &tgen{}
	result, err := tg.parseSetStringValues([]string{"key=123"})
	if err != nil {
		t.Errorf("parseSetStringValues() unexpected error: %v", err)
		return
	}

	expected := map[string]any{"key": "123"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("parseSetStringValues() = %v, want %v", result, expected)
	}
}

// Test advanced functionality through the internal package
func TestAdvancedSetFunctionality(t *testing.T) {
	// Test comma-separated values
	result, err := setflags.ParseSetValues([]string{"a=b,c=d"})
	if err != nil {
		t.Errorf("ParseSetValues() unexpected error: %v", err)
		return
	}

	expected := map[string]any{"a": "b", "c": "d"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ParseSetValues() = %v, want %v", result, expected)
	}

	// Test array syntax
	result, err = setflags.ParseSetValues([]string{"tags={web,api}"})
	if err != nil {
		t.Errorf("ParseSetValues() unexpected error: %v", err)
		return
	}

	expected = map[string]any{"tags": []any{"web", "api"}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ParseSetValues() = %v, want %v", result, expected)
	}
}
