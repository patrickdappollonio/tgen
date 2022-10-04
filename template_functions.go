package main

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
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

	"golang.org/x/text/cases"
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
		"title":      cases.Title,
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
		"seq":           seq,
		"list":          slice,
		"slice":         slice,
		"after":         after,
		"skip":          after,
		"shuffle":       shuffle,
		"first":         first,
		"last":          last,
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

func seq(values ...int) ([]int, error) {
	var start, step, end int

	switch len(values) {
	case 1:
		start = 1
		step = 1
		end = values[0]

	case 2:
		start = values[0]
		step = 1
		end = values[1]

	case 3:
		start = values[0]
		step = values[1]
		end = values[2]

	default:
		return nil, fmt.Errorf("seq: incorrect number of arguments: %d", len(values))
	}

	if step == 0 {
		return nil, fmt.Errorf("seq: step cannot be zero")
	}

	if start < end && step < 0 {
		return nil, fmt.Errorf("seq: increment must be > 0 since %d < %d", start, end)
	}
	if start > end && step > 0 {
		return nil, fmt.Errorf("seq: increment must be > 0 since %d > %d", start, end)
	}

	size := 0
	posstep := step
	if step < 0 {
		posstep = -step
	}

	if end >= start {
		size = (((end - start) / posstep) + 1)
	} else {
		size = (((start - end) / posstep) + 1)

	}

	result := make([]int, int(size))
	value := start
	for i := 0; ; i++ {
		result[i] = value
		value += step

		if (step < 0 && value < end) || (step > 0 && value > end) {
			break
		}
	}

	return result, nil
}

func slice(values ...interface{}) []interface{} {
	return values
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

func shuffle(seq any) (any, error) {
	if seq == nil {
		return nil, errors.New("seq must be provided")
	}

	seqv := reflect.ValueOf(seq)
	seqv, isNil := indirectValue(seqv)
	if isNil {
		return nil, errors.New("can't iterate over a nil value")
	}

	if seqv.Len() == 0 {
		return nil, errors.New("can't shuffle an empty sequence")
	}

	switch seqv.Kind() {
	case reflect.Array, reflect.Slice, reflect.String:
		// skip
	default:
		return nil, errors.New("can't iterate over " + reflect.ValueOf(seq).Type().String())
	}

	shuffled := reflect.MakeSlice(reflect.TypeOf(seq), seqv.Len(), seqv.Len())

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomIndices := rnd.Perm(seqv.Len())

	for index, value := range randomIndices {
		shuffled.Index(value).Set(seqv.Index(index))
	}

	return shuffled.Interface(), nil
}

func first(seq any) (any, error) {
	if seq == nil {
		return nil, errors.New("seq must be provided")
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

	if seqv.Len() == 0 {
		return nil, errors.New("can't get first item of an empty sequence")
	}

	return seqv.Index(0).Interface(), nil
}

func last(seq any) (any, error) {
	if seq == nil {
		return nil, errors.New("seq must be provided")
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

	if seqv.Len() == 0 {
		return nil, errors.New("can't get last item of an empty sequence")
	}

	return seqv.Index(seqv.Len() - 1).Interface(), nil
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
