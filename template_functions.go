package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"text/template"
	"time"
	"unsafe"

	"github.com/Masterminds/sprig/v3"
	"sigs.k8s.io/yaml"
)

func getTemplateFunctions(virtualKV map[string]string, strict bool) template.FuncMap {
	return template.FuncMap{
		"raw":      raw,
		"required": requiredField,

		// Go built-ins
		"lowercase": sprig.FuncMap()["lower"],
		"tolower":   sprig.FuncMap()["lower"],
		"uppercase": sprig.FuncMap()["upper"],
		"toupper":   sprig.FuncMap()["upper"],
		"sprintf":   fmt.Sprintf,
		"printf":    fmt.Sprintf,
		"println":   fmt.Sprintln,

		// Environment functions
		"env":        envstrict(virtualKV, strict),
		"envdefault": envdefault(virtualKV),

		// Locally defined functions
		"rndstring":     rndgen,
		"base64encode":  sprig.FuncMap()["b64enc"],
		"base64decode":  sprig.FuncMap()["b64dec"],
		"readfile":      readfile,
		"readlocalfile": readlocalfile,
		"linebyline":    linebyline,
		"lbl":           linebyline,
		"after":         after,
		"skip":          after,

		"asMap":  asMap,
		"toYAML": toYAML,
	}
}

// requireField returns an error if the given value is nil or an empty string.
func requiredField(warn string, val interface{}) (interface{}, error) {
	if val == nil {
		return val, &requiredError{msg: warn}
	}

	if s, ok := val.(string); ok && s == "" {
		return val, &requiredError{msg: warn}
	}

	return val, nil
}

// toYAML takes an interface, marshals it to yaml, and returns a string. It will
// always return a string, even on marshal error (empty string).
//
// This is designed to be called from a template.
func toYAML(v interface{}) string {
	data, err := yaml.Marshal(v)
	if err != nil {
		// Swallow errors inside of a template.
		return ""
	}
	return strings.TrimSuffix(string(data), "\n")
}

func asMap(m any) map[string]any {
	if m == nil {
		return nil
	}

	newmap, ok := m.(map[string]any)
	if !ok {
		return nil
	}

	return newmap
}

func raw(s string) string {
	return s
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
		return "", &notFoundErr{name: k}
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

// after slices an array to only the items after the Nth item.
func after(index any, seq any) (any, error) {
	if index == nil || seq == nil {
		return nil, errors.New("both limit and seq must be provided")
	}

	indexv, err := tointE(index)
	if err != nil {
		return nil, err
	}

	if indexv < 0 {
		s := fmt.Sprintf("%d", indexv)
		return nil, errors.New("sequence bounds out of range [" + s + ":]")
	}

	seqv := reflect.ValueOf(seq)
	seqv, isNil := indirectValue(seqv)
	if isNil {
		return nil, errors.New("can't iterate over a nil value")
	}

	switch seqv.Kind() {
	case reflect.Array, reflect.Slice, reflect.String:
		// okay
	default:
		return nil, errors.New("can't iterate over " + reflect.ValueOf(seq).Type().String())
	}

	if indexv >= seqv.Len() {
		return seqv.Slice(0, 0).Interface(), nil
	}

	return seqv.Slice(indexv, seqv.Len()).Interface(), nil
}

func toInt(v interface{}) (int, bool) {
	switch v := v.(type) {
	case int:
		return v, true
	case time.Weekday:
		return int(v), true
	case time.Month:
		return int(v), true
	default:
		return 0, false
	}
}

func trimZeroDecimal(s string) string {
	var foundZero bool
	for i := len(s); i > 0; i-- {
		switch s[i-1] {
		case '.':
			if foundZero {
				return s[:i-1]
			}
		case '0':
			foundZero = true
		default:
			return s
		}
	}
	return s
}

func indirect(a interface{}) interface{} {
	if a == nil {
		return nil
	}
	if t := reflect.TypeOf(a); t.Kind() != reflect.Ptr {
		return a
	}
	v := reflect.ValueOf(a)
	for v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}

// indirect is borrowed from the Go stdlib: 'text/template/exec.go'
func indirectValue(v reflect.Value) (rv reflect.Value, isNil bool) {
	for ; v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface; v = v.Elem() {
		if v.IsNil() {
			return v, true
		}
		if v.Kind() == reflect.Interface && v.NumMethod() > 0 {
			break
		}
	}
	return v, false
}

func tointE(i interface{}) (int, error) {
	i = indirect(i)

	intv, ok := toInt(i)
	if ok {
		return intv, nil
	}

	switch s := i.(type) {
	case int64:
		return int(s), nil
	case int32:
		return int(s), nil
	case int16:
		return int(s), nil
	case int8:
		return int(s), nil
	case uint:
		return int(s), nil
	case uint64:
		return int(s), nil
	case uint32:
		return int(s), nil
	case uint16:
		return int(s), nil
	case uint8:
		return int(s), nil
	case float64:
		return int(s), nil
	case float32:
		return int(s), nil
	case string:
		v, err := strconv.ParseInt(trimZeroDecimal(s), 0, 0)
		if err == nil {
			return int(v), nil
		}
		return 0, fmt.Errorf("unable to cast %#v of type %T to int64", i, i)
	case json.Number:
		return tointE(string(s))
	case bool:
		if s {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to int", i, i)
	}
}
