package tfuncs

// requireField returns an error if the given value is nil or an empty string.
func requiredField(warn string, val interface{}) (interface{}, error) {
	if val == nil {
		return val, ErrRequired(warn)
	}

	if s, ok := val.(string); ok && s == "" {
		return val, ErrRequired(warn)
	}

	return val, nil
}

type ErrRequired string

func (r ErrRequired) Error() string {
	return string(r)
}
