package main

import (
	"fmt"
)

type conf struct {
	environmentFile   string
	templateFilePath  string
	stdinTemplateFile string
	valuesFile        string
	strictMode        bool
	customDelimiters  string
}

type enotfounderr struct{ name string }

func (e *enotfounderr) Error() string {
	return "strict mode on: environment variable not found: $" + e.name
}

type emissingkeyerr struct{ name string }

func (e *emissingkeyerr) Error() string {
	return "strict mode on: missing value in values file: " + e.name
}

type conflictingArgsError struct{ F1, F2 string }

func (e *conflictingArgsError) Error() string {
	return fmt.Sprintf("defined both --%s and --%s, only one must be used", e.F1, e.F2)
}
