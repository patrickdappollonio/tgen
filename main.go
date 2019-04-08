package main

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"

	"io"
	"os"
	"path/filepath"
	"strings"
	texttemp "text/template"

	"github.com/spf13/cobra"
)

var (
	envfile, templatefile string

	delimiter string
	version   = "development"

	strict = false

	envvars = make(map[string]string)
)

var root = &cobra.Command{
	Use:          os.Args[0],
	Short:        os.Args[0] + " is a template generator with the power of Go Templates",
	Version:      version,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return command(os.Stdout, delimiter, templatefile, envfile)
	},
}

type enotfounderr struct{ name string }

func (e *enotfounderr) Error() string {
	return "strict mode on: environment variable not found: $" + e.name
}

var envfunc = func(k string) (string, error) {
	k = strings.ToUpper(k)

	if v, found := os.LookupEnv(k); found {
		return v, nil
	}

	if v, found := envvars[k]; found {
		return v, nil
	}

	if strict {
		return "", &enotfounderr{name: k}
	}

	return "", nil
}

var t = template.New(os.Args[0]).Funcs(template.FuncMap{
	"env": envfunc,
	"raw": func(s string) string {
		return s
	},

	"sprintf": func(s string, args ...interface{}) string {
		return fmt.Sprintf(s, args...)
	},

	"envdefault": func(k, defval string) (string, error) {
		s, err := envfunc(k)

		if err != nil {
			if _, ok := err.(*enotfounderr); ok {
				return defval, nil
			}

			return "", err
		}

		if s != "" {
			return s, nil
		}

		return defval, nil
	},
})

func init() {
	root.Flags().StringVarP(&envfile, "environment", "e", "", "an optional environment file to use (key=value formatted) to perform replacements")
	root.Flags().StringVarP(&templatefile, "file", "f", "", "the template file to process (required)")
	root.Flags().StringVarP(&delimiter, "delimiter", "d", "", `delimiter (default "{{}}")`)
	root.Flags().BoolVarP(&strict, "strict", "s", false, "enables strict mode: if an environment variable in the file is defined but not set, it'll fail")
	root.MarkFlagRequired("file")
}

func main() {
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func command(w io.Writer, delim, tmpl, env string) error {
	if delim != "" {
		l, r, err := getDelimiter(delim)
		if err != nil {
			return err
		}

		t = t.Delims(l, r)
	}

	b, err := loadFile(tmpl)
	if err != nil {
		return err
	}

	if err := loadVirtualEnv(env); err != nil {
		return err
	}

	return executeTemplate(w, b)
}

func executeTemplate(w io.Writer, b *bytes.Buffer) error {
	tmpl, err := t.Parse(b.String())
	if err != nil {
		return fmt.Errorf("unable to parse template file %q: %s", templatefile, err.Error())
	}

	var temp bytes.Buffer

	if err := tmpl.Execute(&temp, nil); err != nil {
		if _, ok := err.(texttemp.ExecError); ok {
			if strings.Contains(err.Error(), "environment variable not found") {
				return &enotfounderr{name: err.Error()[strings.LastIndex(err.Error(), ": $")+3:]}
			}
		}

		return err
	}

	if _, err := io.Copy(os.Stdout, &temp); err != nil {
		return err
	}

	return nil
}

func loadFile(fp string) (*bytes.Buffer, error) {
	tmplfile, err := filepath.Abs(fp)
	if err != nil {
		return nil, fmt.Errorf("unable to get path to file %q: %s", fp, err.Error())
	}

	f, err := os.Open(tmplfile)
	if err != nil {
		return nil, fmt.Errorf("unable to open file %q: %s", fp, err.Error())
	}

	defer f.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, f); err != nil {
		return nil, fmt.Errorf("unable to read file %q: %s", fp, err.Error())
	}

	return &buf, nil
}

func loadVirtualEnv(fp string) error {
	if fp == "" {
		return nil
	}

	data, err := loadFile(fp)
	if err != nil {
		return err
	}

	sc := bufio.NewScanner(data)
	for sc.Scan() {
		k, v := parseLine(sc.Text())
		if k == "" || v == "" {
			continue
		}

		envvars[k] = v
	}

	return nil
}

func parseLine(line string) (string, string) {
	if strings.HasPrefix(strings.TrimSpace(line), "#") {
		return "", ""
	}

	items := strings.Split(line, "=")
	if len(items) < 2 {
		return "", ""
	}

	return strings.ToUpper(items[0]), strings.Join(items[1:], "=")
}

func getDelimiter(d string) (string, string, error) {
	size := len(d)

	if size < 2 || size%2 != 0 {
		return "", "", fmt.Errorf("delimiter size needs to be multiple of two and have 2 or more characters")
	}

	div := size / 2
	return d[:div], d[div:], nil
}
