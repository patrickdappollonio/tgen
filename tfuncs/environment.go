package tfuncs

import (
	"os"
	"strings"
)

func envstrict(kv map[string]string, strict bool) func(s string) (string, error) {
	return func(s string) (string, error) {
		return envfunc(s, kv, strict)
	}
}

func envdefault(kv map[string]string) func(k, defval string) (string, error) {
	return func(k, defval string) (string, error) {
		if s, _ := envfunc(k, kv, false); s != "" {
			return s, nil
		}

		return defval, nil
	}
}

func envfunc(k string, kv map[string]string, strictMode bool) (string, error) {
	k = strings.ToUpper(k)

	if v, found := os.LookupEnv(k); found {
		return v, nil
	}

	if v, found := kv[k]; found {
		return v, nil
	}

	if strictMode {
		return "", ErrVarNotFound(k)
	}

	return "", nil
}

type ErrVarNotFound string

func (e ErrVarNotFound) Error() string {
	return "strict mode on: environment variable not found: $" + string(e)
}
