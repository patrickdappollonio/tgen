package setflags

import (
	"fmt"
	"strconv"
	"strings"
)

// KeyValuePair represents a parsed key-value pair
type KeyValuePair struct {
	Key   string
	Value string
}

// KeyPart represents a part of a key path (either a map key or array index)
type KeyPart struct {
	Key     string
	Index   int
	IsArray bool
}

// ParseSetValues parses Helm-style --set values into a nested map structure with type inference
func ParseSetValues(setValues []string) (map[string]any, error) {
	result := make(map[string]any)

	for _, setValue := range setValues {
		if setValue == "" {
			continue
		}

		// Parse comma-separated values within a single --set flag
		pairs, err := parseCommaSeparatedPairs(setValue)
		if err != nil {
			return nil, err
		}

		for _, pair := range pairs {
			key, value := pair.Key, pair.Value

			// Parse the value with type inference (including special values like null, arrays)
			parsedValue, err := parseValueWithAdvancedSyntax(value, true)
			if err != nil {
				return nil, fmt.Errorf("invalid value for key %s: %w", key, err)
			}

			// Set the value in the nested structure
			if err := setNestedValue(result, key, parsedValue); err != nil {
				return nil, err
			}
		}
	}

	return result, nil
}

// ParseSetStringValues parses Helm-style --set-string values into a nested map structure (all values as strings)
func ParseSetStringValues(setValues []string) (map[string]any, error) {
	result := make(map[string]any)

	for _, setValue := range setValues {
		if setValue == "" {
			continue
		}

		// Parse comma-separated values within a single --set-string flag
		pairs, err := parseCommaSeparatedPairs(setValue)
		if err != nil {
			return nil, err
		}

		for _, pair := range pairs {
			key, value := pair.Key, pair.Value

			// Parse the value but force everything to strings (no type inference)
			parsedValue, err := parseValueWithAdvancedSyntax(value, false)
			if err != nil {
				return nil, fmt.Errorf("invalid value for key %s: %w", key, err)
			}

			// Set the value in the nested structure
			if err := setNestedValue(result, key, parsedValue); err != nil {
				return nil, err
			}
		}
	}

	return result, nil
}

// parseCommaSeparatedPairs parses comma-separated key=value pairs with proper escaping
func parseCommaSeparatedPairs(input string) ([]KeyValuePair, error) {
	var pairs []KeyValuePair
	var current strings.Builder
	var inQuotes bool
	var inBraces bool
	var braceLevel int
	var escaped bool

	for _, r := range input {
		if escaped {
			current.WriteRune(r)
			escaped = false
			continue
		}

		if r == '\\' {
			escaped = true
			current.WriteRune(r)
			continue
		}

		if r == '"' {
			inQuotes = !inQuotes
			current.WriteRune(r)
			continue
		}

		if r == '{' && !inQuotes {
			braceLevel++
			inBraces = true
			current.WriteRune(r)
			continue
		}

		if r == '}' && !inQuotes {
			braceLevel--
			if braceLevel == 0 {
				inBraces = false
			}
			current.WriteRune(r)
			continue
		}

		if r == ',' && !inQuotes && !inBraces {
			// Found a separator, process the current pair
			pairStr := strings.TrimSpace(current.String())
			if pairStr != "" {
				pair, err := parseSinglePair(pairStr)
				if err != nil {
					return nil, err
				}
				pairs = append(pairs, pair)
			}
			current.Reset()
			continue
		}

		current.WriteRune(r)
	}

	// Process the last pair
	pairStr := strings.TrimSpace(current.String())
	if pairStr != "" {
		pair, err := parseSinglePair(pairStr)
		if err != nil {
			return nil, err
		}
		pairs = append(pairs, pair)
	}

	return pairs, nil
}

// parseSinglePair parses a single key=value pair
func parseSinglePair(pair string) (KeyValuePair, error) {
	parts := strings.SplitN(pair, "=", 2)
	if len(parts) != 2 {
		return KeyValuePair{}, fmt.Errorf("invalid pair format: %s (expected key=value)", pair)
	}

	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	if key == "" {
		return KeyValuePair{}, fmt.Errorf("empty key in pair: %s", pair)
	}

	return KeyValuePair{Key: key, Value: value}, nil
}

// parseValueWithAdvancedSyntax parses a value with advanced Helm syntax
func parseValueWithAdvancedSyntax(value string, typeInference bool) (any, error) {
	// Handle null values
	if value == "null" {
		return nil, nil
	}

	// Handle empty arrays
	if value == "[]" {
		return []any{}, nil
	}

	// Handle array syntax {a,b,c}
	if strings.HasPrefix(value, "{") && strings.HasSuffix(value, "}") {
		return parseArrayValue(value, typeInference)
	}

	// Handle escaped values
	unescapedValue := unescapeValue(value)

	// Apply type inference if enabled
	if typeInference {
		return parseValueWithTypeInference(unescapedValue), nil
	}

	return unescapedValue, nil
}

// parseArrayValue parses array syntax like {a,b,c}
func parseArrayValue(value string, typeInference bool) ([]any, error) {
	// Remove { and }
	content := strings.TrimSpace(value[1 : len(value)-1])
	if content == "" {
		return []any{}, nil
	}

	// Split by commas, respecting escaping
	items, err := splitArrayItems(content)
	if err != nil {
		return nil, err
	}

	var result []any
	for _, item := range items {
		item = strings.TrimSpace(item)
		if typeInference {
			result = append(result, parseValueWithTypeInference(item))
		} else {
			result = append(result, unescapeValue(item))
		}
	}

	return result, nil
}

// splitArrayItems splits array items by comma, respecting escaping
func splitArrayItems(content string) ([]string, error) {
	var items []string
	var current strings.Builder
	var escaped bool

	for _, r := range content {
		if escaped {
			current.WriteRune(r)
			escaped = false
			continue
		}

		if r == '\\' {
			escaped = true
			continue
		}

		if r == ',' {
			items = append(items, current.String())
			current.Reset()
			continue
		}

		current.WriteRune(r)
	}

	items = append(items, current.String())
	return items, nil
}

// unescapeValue removes escape sequences from a value
func unescapeValue(value string) string {
	var result strings.Builder
	var escaped bool

	for _, r := range value {
		if escaped {
			result.WriteRune(r)
			escaped = false
			continue
		}

		if r == '\\' {
			escaped = true
			continue
		}

		result.WriteRune(r)
	}

	return result.String()
}

// setNestedValue sets a value in a nested map structure, supporting array indexing
func setNestedValue(result map[string]any, key string, value any) error {
	// Parse the key to handle array indexing and dots
	keyParts, err := parseKeyPath(key)
	if err != nil {
		return err
	}

	return setValueAtPath(result, keyParts, value)
}

// parseKeyPath parses a key path like "servers[0].port" or "nodeSelector.\"kubernetes\.io/role\""
func parseKeyPath(key string) ([]KeyPart, error) {
	var parts []KeyPart
	var current strings.Builder
	var inQuotes bool
	var escaped bool
	var lastWasArrayIndex bool

	for i := 0; i < len(key); i++ {
		r := rune(key[i])

		if escaped {
			current.WriteRune(r)
			escaped = false
			continue
		}

		if r == '\\' {
			escaped = true
			continue
		}

		if r == '"' {
			inQuotes = !inQuotes
			continue // Don't include quotes in the key
		}

		if r == '.' && !inQuotes {
			// End of current key part
			partStr := current.String()
			// Only check for empty key parts if we haven't just processed an array index
			if partStr == "" && !lastWasArrayIndex {
				return nil, fmt.Errorf("empty key part in key path: %s", key)
			}
			if partStr != "" {
				parts = append(parts, KeyPart{Key: partStr, Index: -1, IsArray: false})
			}
			current.Reset()
			lastWasArrayIndex = false
			continue
		}

		if r == '[' && !inQuotes {
			// Start of array index
			partStr := current.String()

			// Find the closing bracket
			j := i + 1
			for j < len(key) && key[j] != ']' {
				j++
			}
			if j >= len(key) {
				return nil, fmt.Errorf("unclosed array index in key: %s", key)
			}

			// Parse the index
			indexStr := key[i+1 : j]
			index, err := strconv.Atoi(indexStr)
			if err != nil {
				return nil, fmt.Errorf("invalid array index '%s' in key: %s", indexStr, key)
			}

			// If we have a current key part, add it with the array index
			// If we don't have a current key part, this is a consecutive array index
			if partStr != "" {
				parts = append(parts, KeyPart{Key: partStr, Index: index, IsArray: true})
			} else {
				// This is a consecutive array index like [1] after [0]
				parts = append(parts, KeyPart{Key: "", Index: index, IsArray: true})
			}

			current.Reset()
			lastWasArrayIndex = true
			i = j // Skip to after the closing bracket
			continue
		}

		current.WriteRune(r)
		lastWasArrayIndex = false
	}

	// Add the last part
	partStr := current.String()
	if partStr != "" {
		parts = append(parts, KeyPart{Key: partStr, Index: -1, IsArray: false})
	} else if len(parts) == 0 {
		return nil, fmt.Errorf("empty key path: %s", key)
	}

	return parts, nil
}

// setValueAtPath sets a value at the specified path in a nested structure
func setValueAtPath(result map[string]any, keyParts []KeyPart, value any) error {
	if len(keyParts) == 0 {
		return fmt.Errorf("empty key path")
	}

	// Navigate to the parent of the final key
	current := result
	for _, part := range keyParts[:len(keyParts)-1] {
		if part.IsArray {
			// Handle array indexing
			if _, exists := current[part.Key]; !exists {
				current[part.Key] = []any{}
			}

			slice, ok := current[part.Key].([]any)
			if !ok {
				return fmt.Errorf("conflict: key %s is not an array", part.Key)
			}

			// Extend slice if necessary
			for len(slice) <= part.Index {
				slice = append(slice, make(map[string]any))
			}
			current[part.Key] = slice

			// Move to the indexed element
			elem, ok := slice[part.Index].(map[string]any)
			if !ok {
				elem = make(map[string]any)
				slice[part.Index] = elem
			}
			current = elem
		} else {
			// Handle regular map key
			if _, exists := current[part.Key]; !exists {
				current[part.Key] = make(map[string]any)
			}

			nextMap, ok := current[part.Key].(map[string]any)
			if !ok {
				return fmt.Errorf("conflict: key %s already exists with non-map value", part.Key)
			}
			current = nextMap
		}
	}

	// Set the final value
	finalPart := keyParts[len(keyParts)-1]
	if finalPart.IsArray {
		// Handle array indexing for the final key
		if _, exists := current[finalPart.Key]; !exists {
			current[finalPart.Key] = []any{}
		}

		slice, ok := current[finalPart.Key].([]any)
		if !ok {
			return fmt.Errorf("conflict: key %s is not an array", finalPart.Key)
		}

		// Extend slice if necessary
		for len(slice) <= finalPart.Index {
			slice = append(slice, nil)
		}
		slice[finalPart.Index] = value
		current[finalPart.Key] = slice
	} else {
		// Handle regular map key
		current[finalPart.Key] = value
	}

	return nil
}

// parseValueWithTypeInference attempts to parse a string value into its natural type
func parseValueWithTypeInference(value string) any {
	// Handle booleans first (including YAML booleans)
	if isBooleanValue(value) {
		return parseBooleanValue(value)
	}

	// Try to parse as integer
	if intVal, err := strconv.Atoi(value); err == nil {
		return intVal
	}

	// Try to parse as float
	if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
		return floatVal
	}

	// Default to string
	return value
}

// isBooleanValue checks if a string represents a boolean value (including YAML booleans)
func isBooleanValue(value string) bool {
	lower := strings.ToLower(value)
	return lower == "true" || lower == "false" ||
		lower == "yes" || lower == "no" ||
		lower == "on" || lower == "off"
}

// parseBooleanValue parses a boolean value (including YAML booleans)
func parseBooleanValue(value string) bool {
	lower := strings.ToLower(value)
	return lower == "true" || lower == "yes" || lower == "on"
}
