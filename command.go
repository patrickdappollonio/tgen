package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"text/template"
)

func command(w io.Writer, c conf) error {
	var b *bytes.Buffer

	if c.templateFile != "" {
		if c.rawTemplate != "" {
			return &conflictingArgsError{"file", "raw"}
		}

		if c.stdin {
			return &conflictingArgsError{"file", "stdin"}
		}

		bt, err := loadFile(c.templateFile)
		if err != nil {
			return err
		}

		b = bt
	}

	if c.rawTemplate != "" {
		if c.templateFile != "" {
			return &conflictingArgsError{"raw", "file"}
		}

		if c.stdin {
			return &conflictingArgsError{"raw", "stdin"}
		}

		b = bytes.NewBufferString(c.rawTemplate)
	}

	if c.stdin {
		if c.templateFile != "" {
			return &conflictingArgsError{"stdin", "file"}
		}

		if c.rawTemplate != "" {
			return &conflictingArgsError{"stdin", "raw"}
		}

		bt, err := loadFile(os.Stdin.Name())
		if err != nil {
			return err
		}

		b = bt
	}

	if b == nil {
		return errors.New("needs to specify either a template file (using --file) or a raw template (using --raw or --stdin)")
	}

	envVars, err := loadVirtualEnv(c.environmentFile)
	if err != nil {
		return err
	}

	c.t = template.New(appName).Funcs(getTemplateFunctions(envVars, c.strictMode))

	if c.customDelimiters != "" {
		l, r, err := getDelimiter(c.customDelimiters)
		if err != nil {
			return err
		}

		c.t = c.t.Delims(l, r)
	}

	return executeTemplate(c.t, c.templateFile, w, envVars, b)
}
