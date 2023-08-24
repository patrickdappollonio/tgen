package main

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"sigs.k8s.io/yaml"
)

func Test_tgen_setTemplate(t *testing.T) {
	type args struct {
		name    string
		content string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test",
			args: args{
				name:    "test",
				content: "test test test test test test test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &tgen{}
			tr.setTemplate(tt.args.name, tt.args.content)

			if tr.templateFileName != tt.args.name {
				t.Errorf("templateFileName = %v, want %v", tr.templateFileName, tt.args.name)
			}

			if tr.templateFileContent != tt.args.content {
				t.Errorf("templateFileContent = %v, want %v", tr.templateFileContent, tt.args.content)
			}
		})
	}
}

func writeTemporaryTestFile(t *testing.T, content string) *os.File {
	f, err := os.CreateTemp("", "tgen_test_*")
	if err != nil {
		t.Fatalf("unable to create temporary file: %s", err.Error())
	}

	if err := os.WriteFile(f.Name(), []byte(content), 0644); err != nil {
		t.Fatalf("unable to write to temporary file: %s", err.Error())
	}

	return f
}

func Test_tgen_loadTemplatePath(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{
			name:    "basic",
			content: "test test test test test test test",
		},
		{
			name:    "empty",
			content: "",
		},
		{
			name:    "not found",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var location string

			switch tt.name {
			case "not found":
				location = "/var/log/notfound"

			default:
				f := writeTemporaryTestFile(t, tt.content)
				defer os.Remove(f.Name())
				location = f.Name()
			}

			tg := &tgen{}
			err := tg.loadTemplatePath(location)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if !tt.wantErr {
				if err != nil {
					t.Fatalf("unexpected error: %s", err.Error())
				}

				if tg.templateFileName != location {
					t.Errorf("templateFileName = %v, want %v", tg.templateFileName, location)
				}

				if tg.templateFileContent != tt.content {
					t.Errorf("templateFileContent = %v, want %v", tg.templateFileContent, tt.content)
				}
			}
		})
	}
}

func Test_tgen_loadTemplateFile(t *testing.T) {
	tests := []struct {
		name          string
		content       string
		overwriteName string
		wantErr       bool
	}{
		{
			name:    "basic",
			content: "test test test test test test test",
		},
		{
			name:    "nil pointer",
			wantErr: true,
		},
		{
			name:          "file name overwritten",
			overwriteName: "test",
			content:       "foo bar baz qux quux corge grault garply waldo fred plugh xyzzy thud",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var f *os.File

			switch tt.name {
			case "nil pointer":
				f = nil

			default:
				f = writeTemporaryTestFile(t, tt.content)
				defer os.Remove(f.Name())
			}

			var fileSize int64

			if f != nil {
				fi, err := f.Stat()
				if err != nil {
					t.Fatalf("unable to stat file: %s", err.Error())
				}
				fileSize = fi.Size()
			}

			tr := &tgen{}
			err := tr.loadTemplateFile(tt.overwriteName, f)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
			}

			if !tt.wantErr {
				if err != nil {
					t.Fatalf("unexpected error: %s", err.Error())
				}

				if tt.overwriteName != "" && tr.templateFileName != tt.overwriteName {
					t.Errorf("templateFileName = %q, want %q", tr.templateFileName, tt.overwriteName)
				}

				if a := int64(len(tr.templateFileContent)); a != fileSize {
					t.Errorf("template file size = %v, want %v", a, fileSize)
				}

				if tr.templateFileContent != tt.content {
					t.Errorf("templateFileContent = %q, want %q", tr.templateFileContent, tt.content)
				}
			}
		})
	}
}

func Test_tgen_loadYAMLValues(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{
			name:    "basic",
			content: "foo: bar",
		},
		{
			name:    "empty",
			content: "",
		},
		{
			name:    "invalid",
			content: "[foo90802394uuj^%#",
			wantErr: true,
		},
		{
			name:    "not found",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var path string

			switch tt.name {
			case "not found":
				path = "/var/log/notfound"

			default:
				f := writeTemporaryTestFile(t, tt.content)
				defer os.Remove(f.Name())
				path = f.Name()
			}

			tr := &tgen{}
			err := tr.loadYAMLValues(path)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
			}

			if !tt.wantErr {
				if err != nil {
					t.Fatalf("unexpected error: %s", err.Error())
				}

				want := map[string]any{}
				if err := yaml.Unmarshal([]byte(tt.content), &want); err != nil {
					t.Fatalf("unable to unmarshal test content: %s", err.Error())
				}

				want["Values"] = copyMap(want)
				if !reflect.DeepEqual(tr.yamlValues, want) {
					t.Errorf("values = %v, want %v", tr.yamlValues, want)
				}
			}
		})
	}
}

func mapToKV(m map[string]string) string {
	var b bytes.Buffer

	for k, v := range m {
		b.WriteString(fmt.Sprintf("%s=%s", k, v) + "\n")
	}

	return b.String()
}

func Test_tgen_loadEnvValues(t *testing.T) {
	tests := []struct {
		name    string
		content map[string]string
		wantErr bool
	}{
		{
			name: "basic",
			content: map[string]string{
				"foo":  "bar",
				"baz":  "qux",
				"NAME": "VALUE",
			},
		},
		{
			name:    "empty",
			content: map[string]string{},
		},
		{
			name: "more than one equal",
			content: map[string]string{
				"foo": "===bar",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var path string

			switch tt.name {
			case "not found":
				path = "/var/log/notfound"

			default:
				f := writeTemporaryTestFile(t, mapToKV(tt.content))
				defer os.Remove(f.Name())
				path = f.Name()
			}

			tr := &tgen{}
			err := tr.loadEnvValues(path)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil\ncontents: %#v", tr.envValues)
				}
			}

			if !tt.wantErr {
				if err != nil {
					t.Fatalf("unexpected error: %s", err.Error())
				}

				for kin, vin := range tt.content {
					if v, ok := tr.envValues[strings.ToUpper(kin)]; !ok {
						t.Fatalf("missing key %q", kin)
					} else if v != vin {
						t.Fatalf("value for key %q = %q, want %q", kin, v, vin)
					}
				}
			}
		})
	}
}

func Test_tgen_setDelimiters(t *testing.T) {
	tests := []struct {
		name       string
		delimiters string
		wantLeft   string
		wantRight  string
		wantErr    bool
	}{
		{
			name:       "basic",
			delimiters: "{{}}",
			wantLeft:   "{{",
			wantRight:  "}}",
		},
		{
			name:       "alt",
			delimiters: `<%%>`,
			wantLeft:   "<%",
			wantRight:  "%>",
		},
		{
			name:       "empty",
			delimiters: "",
			wantErr:    true,
		},
		{
			name:       "invalid",
			delimiters: "foo",
			wantErr:    true,
		},
		{
			name:       "invalid",
			delimiters: "foo}}",
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &tgen{}

			if err := tr.setDelimiters(tt.delimiters); (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}

			if tr.preDelimiter != tt.wantLeft {
				t.Errorf("preDelimiter = %q, want %q", tr.preDelimiter, tt.wantLeft)
			}

			if tr.postDelimiter != tt.wantRight {
				t.Errorf("postDelimiter = %q, want %q", tr.postDelimiter, tt.wantRight)
			}
		})
	}
}

func Test_tgen_render(t *testing.T) {
	type fields struct {
		Strict              bool
		templateFileName    string
		templateFileContent string
		yamlValues          map[string]any
		envValues           map[string]string
		preDelimiter        string
		postDelimiter       string
	}

	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "basic",
			fields: fields{
				templateFileName:    "test",
				templateFileContent: `The {{ .Values.foo }} says {{ env "bar" }}`,
				yamlValues: map[string]any{
					"foo": "cow",
				},
				envValues: map[string]string{
					"BAR": "moo",
				},
			},
			want: `The cow says moo`,
		},
		{
			name: "custom delimiters",
			fields: fields{
				templateFileName:    "test",
				templateFileContent: `The <% .Values.foo %> says <% env "bar" %>`,
				yamlValues: map[string]any{
					"foo": "cow",
				},
				envValues: map[string]string{
					"BAR": "moo",
				},
				preDelimiter:  "<%",
				postDelimiter: "%>",
			},
			want: `The cow says moo`,
		},
		{
			name: "missing env no strict",
			fields: fields{
				templateFileName:    "test",
				templateFileContent: `The {{ .Values.foo }} says {{ env "bar" }}`,
				yamlValues: map[string]any{
					"foo": "cow",
				},
				envValues: map[string]string{},
			},
			want: `The cow says `,
		},
		{
			name: "missing env strict",
			fields: fields{
				templateFileName:    "test",
				templateFileContent: `The {{ .Values.foo }} says {{ env "bar" }}`,
				yamlValues: map[string]any{
					"foo": "cow",
				},
				envValues: map[string]string{},
				Strict:    true,
			},
			wantErr: true,
		},
		{
			name: "missing yaml no strict",
			fields: fields{
				templateFileName:    "test",
				templateFileContent: `The {{ .Values.foo }} says {{ env "bar" }}`,
				yamlValues:          map[string]any{},
				envValues: map[string]string{
					"BAR": "moo",
				},
			},
			want: `The  says moo`,
		},
		{
			name: "missing yaml strict",
			fields: fields{
				templateFileName:    "test",
				templateFileContent: `The {{ .Values.foo }} says {{ env "bar" }}`,
				yamlValues:          map[string]any{},
				envValues: map[string]string{
					"BAR": "moo",
				},
				Strict: true,
			},
			wantErr: true,
		},
		{
			name: "missing both no strict",
			fields: fields{
				templateFileName:    "test",
				templateFileContent: `The {{ .Values.foo }} says {{ env "bar" }}`,
				yamlValues:          map[string]any{},
				envValues:           map[string]string{},
			},
			want: `The  says `,
		},
		{
			name: "missing both strict",
			fields: fields{
				templateFileName:    "test",
				templateFileContent: `The {{ .Values.foo }} says {{ env "bar" }}`,
				yamlValues:          map[string]any{},
				envValues:           map[string]string{},
				Strict:              true,
			},
			wantErr: true,
		},
		{
			name: "invalid template",
			fields: fields{
				templateFileName:    "test",
				templateFileContent: `The {{ .Values.foo }} says {{ env "bar" }`,
				yamlValues: map[string]any{
					"foo": "cow",
				},
				envValues: map[string]string{
					"BAR": "moo",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			preDelim, postDelim := tt.fields.preDelimiter, tt.fields.postDelimiter
			if preDelim == "" {
				preDelim = "{{"
			}

			if postDelim == "" {
				postDelim = "}}"
			}

			newvals := copyMap(tt.fields.yamlValues)
			newvals["Values"] = copyMap(newvals)

			tr := &tgen{
				Strict:              tt.fields.Strict,
				templateFileName:    tt.fields.templateFileName,
				templateFileContent: tt.fields.templateFileContent,
				yamlValues:          newvals,
				envValues:           tt.fields.envValues,
				preDelimiter:        preDelim,
				postDelimiter:       postDelim,
			}

			var w bytes.Buffer

			if err := tr.render(&w); (err != nil) != tt.wantErr {
				t.Errorf("render() not expecting an error, got: %q", err.Error())
				return
			}

			if gotW := w.String(); gotW != tt.want {
				t.Errorf("render() = %q, want %q", gotW, tt.want)
			}
		})
	}
}
