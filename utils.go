package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

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

func loadVirtualEnv(fp string) (map[string]string, error) {
	envVars := make(map[string]string)

	if fp == "" {
		return nil, nil
	}

	data, err := loadFile(fp)
	if err != nil {
		return nil, err
	}

	sc := bufio.NewScanner(data)
	for sc.Scan() {
		k, v := parseLine(sc.Text())
		if k == "" || v == "" {
			continue
		}

		envVars[k] = v
	}

	return envVars, nil
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

func executeTemplate(t *template.Template, templateLocation string, output io.Writer, envVars map[string]string, templateData *bytes.Buffer) error {
	tmpl, err := t.Parse(templateData.String())
	if err != nil {
		return fmt.Errorf("unable to parse template file %q: %s", templateLocation, err.Error())
	}

	var temp bytes.Buffer

	if err := tmpl.Execute(&temp, nil); err != nil {
		if _, ok := err.(template.ExecError); ok {
			if strings.Contains(err.Error(), "environment variable not found") {
				return &enotfounderr{name: err.Error()[strings.LastIndex(err.Error(), ": $")+3:]}
			}
		}

		return err
	}

	_, err = io.Copy(output, &temp)
	return err
}
