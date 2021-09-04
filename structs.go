package main

import (
	"fmt"
	"text/template"
)

type conf struct {
	environmentFile  string
	templateFile     string
	rawTemplate      string
	stdin            bool
	customDelimiters string
	strictMode       bool

	t *template.Template
}

type enotfounderr struct{ name string }

func (e *enotfounderr) Error() string {
	return "strict mode on: environment variable not found: $" + e.name
}

type conflictingArgsError struct {
	F1, F2 string
}

func (e *conflictingArgsError) Error() string {
	return fmt.Sprintf("defined both --%s and --%s, only one must be used", e.F1, e.F2)
}
