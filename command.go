package main

import (
	"io"
	"os"
)

func command(w io.Writer, c conf) error {
	// You can't pass "--file" and "--execute" together
	if c.templateFilePath != "" && c.stdinTemplateFile != "" {
		return &conflictingArgsError{"file", "execute"}
	}

	tg := &tgen{Strict: c.strictMode}

	// Read template from "-x" or "--execute" flag
	if c.stdinTemplateFile != "" {
		tg.setTemplate(os.Stdin.Name(), c.stdinTemplateFile)
	}

	// Read template file (either from "--file" or stdin)
	if pathToOpen := c.templateFilePath; pathToOpen != "" {
		var err error
		switch pathToOpen {
		case "-":
			err = tg.loadTemplateFile(os.Stdin.Name(), os.Stdin)
		default:
			err = tg.loadTemplatePath(pathToOpen)
		}

		if err != nil {
			return err
		}
	}

	// Set delimiters
	if c.customDelimiters != "" {
		if err := tg.setDelimiters(c.customDelimiters); err != nil {
			return err
		}
	}

	// Load environment variable file
	if c.environmentFile != "" {
		if err := tg.loadEnvValues(c.environmentFile); err != nil {
			return err
		}
	}

	// Load yaml values file
	if c.valuesFile != "" {
		if err := tg.loadYAMLValues(c.valuesFile); err != nil {
			return err
		}
	}

	// Parse and merge set values
	if len(c.setValues) > 0 {
		if err := tg.mergeSetValues(c.setValues); err != nil {
			return err
		}
	}

	// Parse and merge set-string values
	if len(c.setStringValues) > 0 {
		if err := tg.mergeSetStringValues(c.setStringValues); err != nil {
			return err
		}
	}

	// Render code
	return tg.render(w)
}
