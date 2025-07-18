package tfuncs

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// setupTestDir creates a temporary directory with test files and subdirectories
func setupTestDir(t *testing.T) string {
	t.Helper()

	testDir := t.TempDir()

	// Create test files
	files := []string{
		"file1.txt",
		"file2.txt",
		"subdir/subfile1.txt",
		"subdir/subfile2.txt",
		"subdir/nested/deepfile.txt",
	}

	for _, file := range files {
		fullPath := filepath.Join(testDir, file)
		err := os.MkdirAll(filepath.Dir(fullPath), 0o755)
		if err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}

		err = os.WriteFile(fullPath, []byte("test content"), 0o644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	return testDir
}

func Test_ReadDir(t *testing.T) {
	testDir := setupTestDir(t)

	tests := []struct {
		name    string
		path    string
		want    []string
		wantErr bool
	}{
		{
			name: "read test directory",
			path: testDir,
			want: []string{"file1.txt", "file2.txt", "subdir/"},
		},
		{
			name: "read subdirectory",
			path: filepath.Join(testDir, "subdir"),
			want: []string{"nested/", "subfile1.txt", "subfile2.txt"},
		},
		{
			name:    "read non-existent directory",
			path:    filepath.Join(testDir, "nonexistent"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadDir(tt.path)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readlocaldir(t *testing.T) {
	testDir := setupTestDir(t)

	// Change to test directory for relative path testing
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	err = os.Chdir(testDir)
	if err != nil {
		t.Fatalf("Failed to change to test directory: %v", err)
	}

	tests := []struct {
		name    string
		path    string
		want    []string
		wantErr bool
	}{
		{
			name: "read current directory",
			path: ".",
			want: []string{"file1.txt", "file2.txt", "subdir/"},
		},
		{
			name: "read subdirectory",
			path: "subdir",
			want: []string{"nested/", "subfile1.txt", "subfile2.txt"},
		},
		{
			name:    "absolute path should fail",
			path:    "/tmp",
			wantErr: true,
		},
		{
			name:    "path outside working directory should fail",
			path:    "../..",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readlocaldir(tt.path)

			if (err != nil) != tt.wantErr {
				t.Errorf("readlocaldir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readlocaldir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ReadDirRecursive(t *testing.T) {
	testDir := setupTestDir(t)

	tests := []struct {
		name    string
		path    string
		want    []string
		wantErr bool
	}{
		{
			name: "read test directory recursively",
			path: testDir,
			want: []string{
				"file1.txt",
				"file2.txt",
				"subdir/",
				"subdir/nested/",
				"subdir/nested/deepfile.txt",
				"subdir/subfile1.txt",
				"subdir/subfile2.txt",
			},
		},
		{
			name: "read subdirectory recursively",
			path: filepath.Join(testDir, "subdir"),
			want: []string{
				"nested/",
				"nested/deepfile.txt",
				"subfile1.txt",
				"subfile2.txt",
			},
		},
		{
			name:    "read non-existent directory",
			path:    filepath.Join(testDir, "nonexistent"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadDirRecursive(tt.path)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadDirRecursive() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadDirRecursive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readlocaldirrecursive(t *testing.T) {
	testDir := setupTestDir(t)

	// Change to test directory for relative path testing
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	err = os.Chdir(testDir)
	if err != nil {
		t.Fatalf("Failed to change to test directory: %v", err)
	}

	tests := []struct {
		name    string
		path    string
		want    []string
		wantErr bool
	}{
		{
			name: "read current directory recursively",
			path: ".",
			want: []string{
				"file1.txt",
				"file2.txt",
				"subdir/",
				"subdir/nested/",
				"subdir/nested/deepfile.txt",
				"subdir/subfile1.txt",
				"subdir/subfile2.txt",
			},
		},
		{
			name: "read subdirectory recursively",
			path: "subdir",
			want: []string{
				"nested/",
				"nested/deepfile.txt",
				"subfile1.txt",
				"subfile2.txt",
			},
		},
		{
			name:    "absolute path should fail",
			path:    "/tmp",
			wantErr: true,
		},
		{
			name:    "path outside working directory should fail",
			path:    "../..",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readlocaldirrecursive(tt.path)

			if (err != nil) != tt.wantErr {
				t.Errorf("readlocaldirrecursive() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readlocaldirrecursive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ReadDir_DirectoryTrailingSlash(t *testing.T) {
	testDir := setupTestDir(t)

	// Test that directories have trailing slash
	got, err := ReadDir(testDir)
	if err != nil {
		t.Fatalf("ReadDir() error = %v", err)
	}

	// Check that subdir has trailing slash
	found := false
	for _, entry := range got {
		if entry == "subdir/" {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("ReadDir() should return 'subdir/' with trailing slash, got %v", got)
	}

	// Check that files don't have trailing slash
	for _, entry := range got {
		if strings.HasSuffix(entry, ".txt/") {
			t.Errorf("ReadDir() should not add trailing slash to files, got %v", entry)
		}
	}
}

func Test_ReadDirRecursive_DirectoryTrailingSlash(t *testing.T) {
	testDir := setupTestDir(t)

	// Test that directories have trailing slash in recursive mode
	got, err := ReadDirRecursive(testDir)
	if err != nil {
		t.Fatalf("ReadDirRecursive() error = %v", err)
	}

	// Check that directories have trailing slash
	expectedDirs := []string{"subdir/", "subdir/nested/"}
	for _, expectedDir := range expectedDirs {
		found := false
		for _, entry := range got {
			if entry == expectedDir {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("ReadDirRecursive() should return '%s' with trailing slash, got %v", expectedDir, got)
		}
	}

	// Check that files don't have trailing slash
	for _, entry := range got {
		if strings.HasSuffix(entry, ".txt/") {
			t.Errorf("ReadDirRecursive() should not add trailing slash to files, got %v", entry)
		}
	}
}
