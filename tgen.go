package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

type tgen struct {
	Strict bool

	templateFileName    string
	templateFileContent string
	yamlValues          map[string]any
	envValues           map[string]string

	preDelimiter, postDelimiter string
}

func (t *tgen) setTemplate(name, content string) {
	t.templateFileName = name
	t.templateFileContent = content
}

func (t *tgen) loadTemplatePath(templatepath string) error {
	if templatepath == "" {
		return fmt.Errorf("template path is empty")
	}

	bf, err := readfile(templatepath)
	if err != nil {
		return err
	}

	t.templateFileName = templatepath
	t.templateFileContent = bf
	return nil
}

func (t *tgen) loadTemplateFile(overwriteName string, f *os.File) error {
	if f == nil {
		return fmt.Errorf("template file is empty")
	}

	name := f.Name()
	if overwriteName != "" {
		name = overwriteName
	}

	var buf bytes.Buffer
	read, err := io.Copy(&buf, f)

	if err != nil {
		return fmt.Errorf("unable to read template file %q: %w", name, err)
	}

	if read == 0 {
		return fmt.Errorf("template file %q is empty", name)
	}

	t.templateFileName = name
	t.templateFileContent = buf.String()
	return nil
}

func (t *tgen) loadYAMLValues(yamlpath string) error {
	if yamlpath == "" {
		return fmt.Errorf("yaml values file path is empty")
	}

	valuesfile := map[string]any{}

	bf, err := readfile(yamlpath)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal([]byte(bf), &valuesfile); err != nil {
		return fmt.Errorf("unable to parse values file %q: %s", yamlpath, err.Error())
	}

	copied, err := copyMap(valuesfile)
	if err != nil {
		return fmt.Errorf("unable to duplicate values file: %s", err.Error())
	}

	valuesfile["Values"] = copied
	t.yamlValues = valuesfile
	return nil
}

func (t *tgen) loadEnvValues(envpath string) error {
	envVars := make(map[string]string)

	if envpath == "" {
		return nil
	}

	data, err := readfile(envpath)
	if err != nil {
		return err
	}

	sc := bufio.NewScanner(bytes.NewBufferString(data))
	for sc.Scan() {
		key, value, err := parseEnvLine(sc.Text())
		if err != nil {
			return err
		}

		if key != "" && value != "" {
			envVars[key] = value
		}
	}

	t.envValues = envVars
	return nil
}

func (t *tgen) setDelimiters(delimiters string) error {
	size := len(delimiters)

	if size < 2 || size%2 != 0 {
		return fmt.Errorf("delimiter size needs to be multiple of two and have 2 or more characters")
	}

	div := size / 2
	t.preDelimiter = delimiters[:div]
	t.postDelimiter = delimiters[div:]
	return nil
}

func (t *tgen) render(w io.Writer) error {
	baseTemplate := template.New(t.templateFileName).Funcs(getTemplateFunctions(t.envValues, t.Strict))

	if t.Strict {
		baseTemplate = baseTemplate.Option("missingkey=error")
	} else {
		baseTemplate = baseTemplate.Option("missingkey=zero")
	}

	if t.preDelimiter != "" && t.postDelimiter != "" {
		baseTemplate = baseTemplate.Delims(t.preDelimiter, t.postDelimiter)
	}

	var temp bytes.Buffer

	parsed, err := baseTemplate.Parse(t.templateFileContent)
	if err != nil {
		return fmt.Errorf("unable to parse template file %q: %s", t.templateFileName, err.Error())
	}

	if err := parsed.Execute(&temp, t.yamlValues); err != nil {
		return t.replaceTemplateRenderError(err)
	}

	if t.Strict {
		_, err = fmt.Fprint(w, temp.String())
		return err
	}

	// Due to an unfortunate agreement and lack of behaviour change in the Go standard
	// library, I'm forced to trim the <no value> string from the output directly.
	// See helm's engine implementation of this
	// https://github.com/helm/helm/blob/7ed9d16dc764a5b94b378a7e217865efaa0d9ac8/pkg/engine/engine.go#L267
	// and the original issue, not solved but closed as wontfix:
	// https://github.com/golang/go/issues/24963
	str := strings.ReplaceAll(temp.String(), "<no value>", "")
	_, err = fmt.Fprint(w, str)
	return err
}

func (t *tgen) replaceTemplateRenderError(err error) error {
	if err == nil {
		return nil
	}

	if _, ok := err.(template.ExecError); ok {
		if strings.Contains(err.Error(), "environment variable not found") {
			return &enotfounderr{name: err.Error()[strings.LastIndex(err.Error(), ": $")+3:]}
		}

		if strings.Contains(err.Error(), "map has no entry for key") {
			return &emissingkeyerr{name: err.Error()[strings.LastIndex(err.Error(), ":")+2:]}
		}
	}

	return err
}
