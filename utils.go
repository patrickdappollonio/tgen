package main

import (
	"fmt"
	"strings"
)

func parseEnvLine(line string) (string, string, error) {
	line = strings.TrimSpace(line)

	if line == "" {
		return "", "", nil
	}

	if strings.HasPrefix(line, "#") {
		return "", "", nil
	}

	key, value, found := strings.Cut(line, "=")
	if !found {
		return "", "", fmt.Errorf("invalid environment line: key=value separator not found: %q", line)
	}

	key = strings.ToUpper(strings.TrimSpace(key))

	if key == "" {
		return "", "", fmt.Errorf("key empty for environment line: %q", line)
	}

	value = strings.TrimSpace(value)

	if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
		value = strings.TrimPrefix(value, "\"")
		value = strings.TrimSuffix(value, "\"")
	}

	return key, value, nil
}

func copyMap(m map[string]any) map[string]any {
	cp := make(map[string]interface{})
	for k, v := range m {
		vm, ok := v.(map[string]interface{})
		if ok {
			cp[k] = copyMap(vm)
		} else {
			cp[k] = v
		}
	}

	return cp
}

// mergeMap deeply merges two maps, with values from the second map taking precedence
func mergeMap(dest, src map[string]any) map[string]any {
	if dest == nil {
		dest = make(map[string]any)
	}

	result := copyMap(dest)

	for k, v := range src {
		if srcMap, ok := v.(map[string]any); ok {
			if destMap, exists := result[k]; exists {
				if destMapTyped, ok := destMap.(map[string]any); ok {
					result[k] = mergeMap(destMapTyped, srcMap)
				} else {
					// If destination is not a map, replace it
					result[k] = copyMap(srcMap)
				}
			} else {
				result[k] = copyMap(srcMap)
			}
		} else {
			// For non-map values, the source value takes precedence
			result[k] = v
		}
	}

	return result
}
