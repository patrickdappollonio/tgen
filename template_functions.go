package main

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
	"unsafe"
)

func getTemplateFunctions(virtualKV map[string]string, strict bool) template.FuncMap {
	return template.FuncMap{
		"raw": func(s string) string {
			return s
		},

		// Go built-ins
		"lowercase":  strings.ToLower,
		"lower":      strings.ToLower,
		"uppercase":  strings.ToUpper,
		"upper":      strings.ToUpper,
		"title":      strings.Title,
		"sprintf":    fmt.Sprintf,
		"printf":     fmt.Sprintf,
		"println":    fmt.Sprintln,
		"trim":       strings.TrimSpace,
		"trimPrefix": strings.TrimPrefix,
		"trimSuffix": strings.TrimSuffix,
		"split":      strings.Split,
		"base":       filepath.Base,
		"dir":        filepath.Dir,
		"clean":      filepath.Clean,
		"ext":        filepath.Ext,
		"isAbs":      filepath.IsAbs,

		// Environment functions
		"env":        envstrict(virtualKV, strict),
		"envdefault": envdefault(virtualKV),

		// Locally defined functions
		"rndstring":     rndgen,
		"repeat":        repeat,
		"nospace":       nospace,
		"quote":         quote,
		"squote":        squote,
		"indent":        indent,
		"nindent":       nindent,
		"b64enc":        base64encode,
		"base64encode":  base64encode,
		"b64dec":        base64decode,
		"base64decode":  base64decode,
		"sha1sum":       sha1sum,
		"sha256sum":     sha256sum,
		"replace":       replace,
		"readfile":      readfile,
		"readlocalfile": readlocalfile,
		"linebyline":    linebyline,
		"lbl":           linebyline,
	}
}

func readfile(path string) (string, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(contents), nil
}

func readlocalfile(path string) (string, error) {
	if filepath.IsAbs(path) {
		return "", fmt.Errorf("unable to open local file %q: path is absolute, only relative paths are allowed on \"readlocalfile\"", path)
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	cleanpath := filepath.Join(wd, path)

	if cleanpath == wd {
		return "", &fs.PathError{Op: "read", Path: cleanpath, Err: errors.New("is a directory")}
	}

	if !strings.HasPrefix(cleanpath, wd) {
		return "", fmt.Errorf("unable to open local file %q: file is not under current working directory", cleanpath)
	}

	return readfile(cleanpath)
}

func linebyline(lines string) []string {
	return strings.Split(lines, "\n")
}

func replace(old, new, src string) string {
	return strings.Replace(src, old, new, -1)
}

func sha256sum(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}

func sha1sum(input string) string {
	hash := sha1.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}

func base64encode(v string) string {
	return base64.StdEncoding.EncodeToString([]byte(v))
}

func base64decode(v string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(v)

	if err != nil {
		return "", err
	}

	return string(data), nil
}

func nospace(str string) string {
	return strings.NewReplacer(" ", "").Replace(str)
}

func repeat(count int, str string) string {
	return strings.Repeat(str, count)
}

func indent(spaces int, v string) string {
	pad := strings.Repeat(" ", spaces)
	return pad + strings.Replace(v, "\n", "\n"+pad, -1)
}

func nindent(spaces int, v string) string {
	return "\n" + indent(spaces, v)
}

func quote(s string) string  { return `"` + s + `"` }
func squote(s string) string { return `'` + s + `'` }

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
		return "", &enotfounderr{name: k}
	}

	return "", nil
}

// from: https://stackoverflow.com/a/31832326

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func rndgen(n int) string {
	var src = rand.NewSource(time.Now().UnixNano())

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
