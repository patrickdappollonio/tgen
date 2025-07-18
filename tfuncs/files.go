package tfuncs

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ReadFile reads the contents of the file at the specified path and returns it as a string.
// The path can be either relative or absolute. This function can read any file that the
// process has access to, including system files outside the current working directory.
func ReadFile(path string) (string, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(contents), nil
}

// readLocalFile reads the contents of a file and returns it as a string, but only allows
// relative paths within the current working directory and its subdirectories.
// This function provides security by preventing access to files outside the current
// working directory through path traversal attacks.
//
// Returns an error if:
//   - The path is absolute
//   - The resolved path is outside the current working directory
//   - The path points to a directory instead of a file
//   - The file cannot be read
func readLocalFile(path string) (string, error) {
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

// readDir reads the contents of a directory and returns a sorted slice of entry names.
// Directories are returned with a trailing "/" to distinguish them from files.
// The path can be either relative or absolute. This function can read any directory
// that the process has access to, including system directories.
//
// The returned slice is sorted alphabetically for consistent output.
// Symbolic links are not followed.
func readDir(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() {
			name += "/"
		}
		result = append(result, name)
	}

	sort.Strings(result)
	return result, nil
}

// readLocalDir reads the contents of a directory and returns a sorted slice of entry names,
// but only allows relative paths within the current working directory and its subdirectories.
// Directories are returned with a trailing "/" to distinguish them from files.
// This function provides security by preventing access to directories outside the current
// working directory through path traversal attacks.
//
// The returned slice is sorted alphabetically for consistent output.
// Symbolic links are not followed.
//
// Returns an error if:
//   - The path is absolute
//   - The resolved path is outside the current working directory
//   - The directory cannot be read
func readLocalDir(path string) ([]string, error) {
	if filepath.IsAbs(path) {
		return nil, fmt.Errorf("unable to open local directory %q: path is absolute, only relative paths are allowed on \"readlocaldir\"", path)
	}

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	cleanpath := filepath.Join(wd, path)

	if !strings.HasPrefix(cleanpath, wd) {
		return nil, fmt.Errorf("unable to open local directory %q: directory is not under current working directory", cleanpath)
	}

	return readDir(cleanpath)
}

// readDirRecursive reads the contents of a directory recursively and returns a sorted slice
// of all file and directory paths relative to the root directory. Directories are returned
// with a trailing "/" to distinguish them from files. The path can be either relative or
// absolute. This function can read any directory that the process has access to, including
// system directories.
//
// The returned paths use forward slashes for consistency across platforms and are sorted
// alphabetically for consistent output. The root directory itself is not included in the results.
// Symbolic links are not followed.
//
// For example, if the directory structure is:
//
//	testdata/
//	├── file1.txt
//	├── file2.txt
//	└── subdir/
//	    └── subfile.txt
//
// The function would return: ["file1.txt", "file2.txt", "subdir/", "subdir/subfile.txt"]
func readDirRecursive(path string) ([]string, error) {
	var result []string

	err := filepath.WalkDir(path, func(walkPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory itself
		if walkPath == path {
			return nil
		}

		// Get relative path from the root
		relPath, err := filepath.Rel(path, walkPath)
		if err != nil {
			return err
		}

		// Convert to forward slashes for consistency
		relPath = filepath.ToSlash(relPath)

		// Add trailing slash for directories
		if d.IsDir() {
			relPath += "/"
		}

		result = append(result, relPath)
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Strings(result)
	return result, nil
}

// readLocalDirRecursive reads the contents of a directory recursively and returns a sorted
// slice of all file and directory paths relative to the root directory, but only allows
// relative paths within the current working directory and its subdirectories.
// Directories are returned with a trailing "/" to distinguish them from files.
// This function provides security by preventing access to directories outside the current
// working directory through path traversal attacks.
//
// The returned paths use forward slashes for consistency across platforms and are sorted
// alphabetically for consistent output. The root directory itself is not included in the results.
// Symbolic links are not followed.
//
// Returns an error if:
//   - The path is absolute
//   - The resolved path is outside the current working directory
//   - The directory cannot be read
//
// For example, if the directory structure is:
//
//	testdata/
//	├── file1.txt
//	├── file2.txt
//	└── subdir/
//	    └── subfile.txt
//
// The function would return: ["file1.txt", "file2.txt", "subdir/", "subdir/subfile.txt"]
func readLocalDirRecursive(path string) ([]string, error) {
	if filepath.IsAbs(path) {
		return nil, fmt.Errorf("unable to open local directory %q: path is absolute, only relative paths are allowed on \"readlocaldirrecursive\"", path)
	}

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	cleanpath := filepath.Join(wd, path)

	if !strings.HasPrefix(cleanpath, wd) {
		return nil, fmt.Errorf("unable to open local directory %q: directory is not under current working directory", cleanpath)
	}

	return readDirRecursive(cleanpath)
}
