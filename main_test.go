package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestVersion(t *testing.T) {
	cases := []struct {
		sent     string
		expected string
	}{
		{"1.2.0", "v1.2.0"},
		{"2.2", "v2.2"},
		{"3", "v3"},
		{"development", "development"},
		{"devel", "devel"},
	}

	for _, v := range cases {
		t.Run(v.expected, func(tt *testing.T) {
			if expected := getVersion(v.sent); v.expected != expected {
				tt.Fatalf("expecting %q, but got %q for string %q", v.expected, expected, v.sent)
			}
		})
	}
}

func TestDelimiter(t *testing.T) {
	cases := []struct {
		sent  string
		left  string
		right string
		fails bool
	}{
		{"<<>>", "<<", ">>", false},
		{"{{}}", "{{", "}}", false},
		{"{}", "{", "}", false},
		{"ab", "a", "b", false},
		{"abc", "", "", true},
	}

	for _, v := range cases {
		t.Run(v.sent, func(tt *testing.T) {
			l, r, err := getDelimiter(v.sent)

			if v.fails && err == nil {
				tt.Fatalf("expecting function to fail, but got no error")
			} else if !v.fails && err != nil {
				tt.Fatalf("not expecting to fail, but got %q", err.Error())
			}

			if v.left != l {
				tt.Fatalf("expecting left side to be %q, got %q", v.left, l)
			}

			if v.right != r {
				tt.Fatalf("expecting right side to be %q, got %q", v.right, r)
			}
		})
	}
}

func TestParseLine(t *testing.T) {
	cases := []struct {
		sent  string
		left  string
		right string
	}{
		{"a=b", "A", "b"},
		{"c", "", ""},
		{"c=", "C", ""},
	}

	for _, v := range cases {
		t.Run(v.sent, func(tt *testing.T) {
			l, r := parseLine(v.sent)

			if v.left != l {
				tt.Fatalf("expecting left side to be %q, got %q", v.left, l)
			}

			if v.right != r {
				tt.Fatalf("expecting right side to be %q, got %q", v.right, r)
			}
		})
	}
}

func TestLoadVirtualEnv(t *testing.T) {
	contents := `# This is a comment
my_const_user=demo
MY_UPPER=abc`

	f, err := ioutil.TempFile(os.TempDir(), "testing_")
	if err != nil {
		t.Fatalf("not expecting an error creating temp file, got %s", err.Error())
	}
	defer os.Remove(f.Name())

	f.WriteString(contents)
	f.Close()

	if err := loadVirtualEnv(f.Name()); err != nil {
		t.Fatalf("not expecting an error loading virtualenv, got %s", err.Error())
	}

	upper, found := virtualenv["MY_CONST_USER"]
	if !found {
		t.Fatalf("expecting to find key MY_CONST_USER, but wasn't found in virtualenv")
	}

	if upper != "demo" {
		t.Fatalf("expecting MY_CONST_USER to say \"demo\" but got %q", upper)
	}

	lower, found := virtualenv["MY_UPPER"]
	if !found {
		t.Fatalf("expecting to find key MY_UPPER, but wasn't found in virtualenv")
	}

	if lower != "abc" {
		t.Fatalf("expecting MY_UPPER to say \"abc\" but got %q", lower)
	}

	if err := loadVirtualEnv(""); err != nil {
		t.Fatalf("not expecting loadVirtualEnv to return an error when calling with empty, but got %s", err.Error())
	}

	err = loadVirtualEnv("/this/doesn't/exist")
	if err == nil {
		t.Fatalf("expecting loadVirtualEnv to fail, but got no error")
	}
}
