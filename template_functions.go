package main

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
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
		"lowercase":  strings.ToLower,    // "HELLO" → "hello"
		"lower":      strings.ToLower,    // "HELLO" → "hello"
		"uppercase":  strings.ToUpper,    // "hello" → "HELLO"
		"upper":      strings.ToUpper,    // "hello" → "HELLO"
		"title":      strings.Title,      // "hello" → "Hello"
		"sprintf":    fmt.Sprintf,        // sprintf "Hello, %s" "world" → "Hello, world"
		"printf":     fmt.Sprintf,        // printf "Hello, %s" "world" → "Hello, world"
		"println":    fmt.Sprintln,       // println "Hello" "world!" → "Hello world!\n"
		"trim":       strings.TrimSpace,  // trim "   hello   " → "hello"
		"trimPrefix": strings.TrimPrefix, // trimPrefix "abcdef" "abc" → "def"
		"trimSuffix": strings.TrimSuffix, // trimSuffix "abcdef" "def" → "abc"
		"base":       filepath.Base,      // base "/foo/bar/baz" → "baz"
		"dir":        filepath.Dir,       // dir "/foo/bar/baz" → "/foo/bar"
		"clean":      filepath.Clean,     // clean "/foo/bar/../baz" → "/foo/baz"
		"ext":        filepath.Ext,       // ext "/foo.zip" → ".zip"
		"isAbs":      filepath.IsAbs,     // isAbs "foo.zip" → false

		// Environment functions
		"env":        envstrict(virtualKV, strict), // env "user" → "patrick"
		"envdefault": envdefault(virtualKV),        // env "SQL_HOST" "sql.example.com" → "sql.example.com"

		// Locally defined functions
		"rndstring":    rndgen,       // rndstring 8 → "lFEqUUOJ"
		"repeat":       repeat,       // repeat 3 "abc" → "abcabcabc"
		"nospace":      nospace,      // nospace "hello world!" → "helloworld!"
		"quote":        quote,        // quote "hey" → `"hey"`
		"squote":       squote,       // squote "hey" → "'hey'"
		"indent":       indent,       // indent 3 "abc" → "  abc"
		"nindent":      nindent,      // nindent 3 "abc" → "\n   abc"
		"b64enc":       base64encode, // b64enc "abc" → "YWJj"
		"base64encode": base64encode, // base64encode "abc" → "YWJj"
		"b64dec":       base64decode, // b64dec "YWJj" → "abc"
		"base64decode": base64decode, // base64decode "YWJj" → "abc"
		"sha1sum":      sha1sum,      // sha1sum "abc" → "a9993e364706816aba3e25717850c26c9cd0d89d"
		"sha256sum":    sha256sum,    // sha256sum "abc" → "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"
		"replace":      replace,      // replace "World" "Patrick" "Hello, World!" → "Hello, Patrick!"
		"readfile":     readfile,     // readfile "foobar.txt" → "Hello, world!"
		"linebyline":   linebyline,   // linebyline "foo\nbar" → ["foo", "bar"]
		"lbl":          linebyline,   // linebyline "foo\nbar" → ["foo", "bar"]
	}
}

func readfile(path string) (string, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(contents), nil
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
