package main

import (
	"bytes"
	"encoding/gob"
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

func copyMap(m map[string]any) (map[string]any, error) {
	gob.Register(map[string]any{})

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)
	err := enc.Encode(m)
	if err != nil {
		return nil, err
	}
	var copy map[string]any
	err = dec.Decode(&copy)
	if err != nil {
		return nil, err
	}
	return copy, nil
}
