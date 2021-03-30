package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"text/template"
	"time"
	"unsafe"
)

func getTemplateFunctions(strict bool) template.FuncMap {
	return template.FuncMap{
		"env": envstrict(strict),

		"raw": func(s string) string {
			return s
		},

		"envdefault": envdefault,

		"rndstring": rndgen,
		"lowercase": strings.ToLower,
		"lower":     strings.ToLower,
		"uppercase": strings.ToUpper,
		"upper":     strings.ToUpper,
		"title":     strings.Title,
		"sprintf":   fmt.Sprintf,
		"printf":    fmt.Sprintf,
		"println":   fmt.Sprintln,
		"trim":      strings.TrimSpace,

		"trimPrefix": strings.TrimPrefix,
		"trimSuffix": strings.TrimSuffix,

		"repeat": func(count int, str string) string {
			return strings.Repeat(str, count)
		},

		"nospace": func(str string) string {
			return strings.NewReplacer(" ", "").Replace(str)
		},

		"quote":  quote,
		"squote": squote,

		"indent":  indent,
		"nindent": nindent,
	}
}

func indent(spaces int, v string) string {
	pad := strings.Repeat(" ", spaces)
	return pad + strings.Replace(v, "\n", "\n"+pad, -1)
}

func nindent(spaces int, v string) string {
	return "\n" + indent(spaces, v)
}

func envdefault(k, defval string) (string, error) {
	if s, _ := envfunc(k, false); s != "" {
		return s, nil
	}

	return defval, nil
}

func quote(s string) string  { return `"` + s + `"` }
func squote(s string) string { return `'` + s + `'` }

func envstrict(strict bool) func(s string) (string, error) {
	return func(s string) (string, error) {
		return envfunc(s, strict)
	}
}

func envfunc(k string, strictMode bool) (string, error) {
	k = strings.ToUpper(k)

	if v, found := os.LookupEnv(k); found {
		return v, nil
	}

	if v, found := loadedEnvVars[k]; found {
		return v, nil
	}

	if strictMode {
		return "", &enotfounderr{name: k}
	}

	return "", nil
}

// from: https://stackoverflow.com/a/31832326

var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func rndgen(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}
