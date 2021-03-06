package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"text/template"

	"github.com/spf13/cobra"
)

const appName = "tgen"

type conf struct {
	environmentFile  string
	templateFile     string
	rawTemplate      string
	stdin            bool
	customDelimiters string
	strictMode       bool

	t *template.Template
}

var (
	version       = "development"
	loadedEnvVars = make(map[string]string)
)

func main() {
	var configs conf

	var root = &cobra.Command{
		Use:          appName,
		Short:        appName + " is a template generator with the power of Go Templates",
		Version:      version,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return command(os.Stdout, configs)
		},
	}

	root.Flags().StringVarP(&configs.environmentFile, "environment", "e", "", "an optional environment file to use (key=value formatted) to perform replacements")
	root.Flags().StringVarP(&configs.templateFile, "file", "f", "", "the template file to process (required)")
	root.Flags().StringVarP(&configs.customDelimiters, "delimiter", "d", "", `delimiter (default "{{}}")`)
	root.Flags().StringVarP(&configs.rawTemplate, "execute", "x", "", "a raw template to execute directly, without providing --file")
	root.Flags().BoolVarP(&configs.stdin, "stdin", "i", false, "a stdin input to execute directly, without providing --file or --execute")
	root.Flags().BoolVarP(&configs.strictMode, "strict", "s", false, "enables strict mode: if an environment variable in the file is defined but not set, it'll fail")

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func command(w io.Writer, c conf) error {
	c.t = template.New(appName).Funcs(getTemplateFunctions(c.strictMode))

	if c.customDelimiters != "" {
		l, r, err := getDelimiter(c.customDelimiters)
		if err != nil {
			return err
		}

		c.t = c.t.Delims(l, r)
	}

	var b *bytes.Buffer

	if c.templateFile != "" {
		if c.rawTemplate != "" {
			return errors.New("defined both --file and --raw, only one must be used")
		}

		if c.stdin {
			return errors.New("defined both --file and --stdin, only one must be used")
		}

		bt, err := loadFile(c.templateFile)
		if err != nil {
			return err
		}

		b = bt
	}

	if c.rawTemplate != "" {
		if c.templateFile != "" {
			return errors.New("defined both --raw and --file, only one must be used")
		}

		if c.stdin {
			return errors.New("defined both --raw and --stdin, only one must be used")
		}

		b = bytes.NewBufferString(c.rawTemplate)
	}

	if c.stdin {
		if c.templateFile != "" {
			return errors.New("defined both --stdin and --file, only one must be used")
		}

		if c.rawTemplate != "" {
			return errors.New("defined both --stdin and --raw, only one must be used")
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

	if err := loadVirtualEnv(c.environmentFile); err != nil {
		return err
	}

	return executeTemplate(c.t, c.templateFile, w, b)
}
