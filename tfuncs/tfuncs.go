package tfuncs

import (
	"fmt"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

func GetFunctions(virtualKV map[string]string, strict bool) template.FuncMap {
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
		"rndstring":             rndgen,
		"base64encode":          sprig.FuncMap()["b64enc"],
		"base64decode":          sprig.FuncMap()["b64dec"],
		"readfile":              ReadFile,
		"readlocalfile":         readLocalFile,
		"readdir":               readDir,
		"readlocaldir":          readLocalDir,
		"readdirrecursive":      readDirRecursive,
		"readlocaldirrecursive": readLocalDirRecursive,
		"linebyline":            linebyline,
		"lbl":                   linebyline,
		"after":                 after,
		"skip":                  after,

		"asMap":   asMap,
		"toYAML":  toYAML,
		"rnditem": rnditem[any],
	}
}
