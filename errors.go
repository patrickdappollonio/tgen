package main

import "fmt"

type requiredError struct {
	msg string
}

func (r requiredError) Error() string {
	return r.msg
}

type templateFuncError struct {
	line     string
	original error
}

func (t templateFuncError) Error() string {
	if t.line == "" {
		return t.original.Error()
	}

	return fmt.Sprintf("evaluating %s: %s", t.line, t.original)
}

func (t templateFuncError) Unwrap() error {
	return t.original
}

type notFoundErr struct{ name string }

func (e *notFoundErr) Error() string {
	return "strict mode on: environment variable not found: $" + e.name
}

type missingKeyErr struct{ name string }

func (e *missingKeyErr) Error() string {
	return "strict mode on: missing value in values file: " + e.name
}

type conflictingArgsError struct{ F1, F2 string }

func (e *conflictingArgsError) Error() string {
	return fmt.Sprintf("defined both --%s and --%s, only one must be used", e.F1, e.F2)
}
