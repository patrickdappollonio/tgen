package tfuncs

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func ReadFile(path string) (string, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(contents), nil
}

func readlocalfile(path string) (string, error) {
	if filepath.IsAbs(path) {
		return "", fmt.Errorf("unable to open local file %q: path is absolute, only relative paths are allowed on \"readlocalfile\"", path)
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	cleanpath := filepath.Join(wd, path)

	if cleanpath == wd {
		return "", &fs.PathError{Op: "read", Path: cleanpath, Err: errors.New("is a directory")}
	}

	if !strings.HasPrefix(cleanpath, wd) {
		return "", fmt.Errorf("unable to open local file %q: file is not under current working directory", cleanpath)
	}

	return ReadFile(cleanpath)
}
