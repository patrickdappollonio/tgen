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
		"env": func(k string) (string, error) {
			return envfunc(k, strict)
		},
		"raw": func(s string) string {
			return s
		},

		"sprintf": func(s string, args ...interface{}) string {
			return fmt.Sprintf(s, args...)
		},

		"envdefault": func(k, defval string) (string, error) {
			if s, _ := envfunc(k, false); s != "" {
				return s, nil
			}

			return defval, nil
		},

		"rndstring": rndgen,
		"lowercase": strings.ToLower,
		"uppercase": strings.ToUpper,
		"title":     strings.Title,
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
