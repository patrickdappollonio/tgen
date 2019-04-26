package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"os"
	"strings"
	"time"
	"unsafe"
)

var templateFunctions = template.FuncMap{
	"env": envfunc,
	"raw": func(s string) string {
		return s
	},

	"sprintf": func(s string, args ...interface{}) string {
		return fmt.Sprintf(s, args...)
	},

	"envdefault": func(k, defval string) (string, error) {
		s, err := envfunc(k)

		if err != nil {
			if _, ok := err.(*enotfounderr); ok {
				return defval, nil
			}

			return "", err
		}

		if s != "" {
			return s, nil
		}

		return defval, nil
	},

	"rndstring": rndgen,
	"lowercase": strings.ToLower,
	"uppercase": strings.ToUpper,
	"title":     strings.Title,
}

func envfunc(k string) (string, error) {
	k = strings.ToUpper(k)

	if v, found := os.LookupEnv(k); found {
		return v, nil
	}

	if v, found := envvars[k]; found {
		return v, nil
	}

	if strict {
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
