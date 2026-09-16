package util

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolvePath(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"~/documents/file.txt", os.Getenv("HOME") + "/documents/file.txt"},
		{"/absolute/path/file.txt", "/absolute/path/file.txt"},
		{"", ""},
		{"$HOME/documents/file.txt", os.Getenv("HOME") + "/documents/file.txt"},
		{"$HOME/$USER/documents/file.txt", os.Getenv("HOME") + "/" + os.Getenv("USER") + "/documents/file.txt"},
	}

	for _, test := range tests {
		result := ResolvePath(test.input)
		if result != test.expected {
			t.Errorf("For input %q, expected %q, but got %q", test.input, test.expected, result)
		}
	}
}

func TestResolveFilePath(t *testing.T) {
	home := os.Getenv("HOME")

	tests := []struct {
		input   string
		baseDir string
		abs     bool
	}{
		{"~/documents/file.txt", "/tmp/base", true},
		{"/absolute/path/file.txt", "/tmp/base", true},
		{"relative/file.txt", "/tmp/base", true},
		{"$HOME/documents/file.txt", "/tmp/base", true},
	}

	for _, test := range tests {
		result := ResolveFilePath(test.input, test.baseDir)
		if !filepath.IsAbs(result) {
			t.Errorf("For input %q, expected absolute path, got %q", test.input, result)
		}
		if test.input == "relative/file.txt" && !strings.HasPrefix(result, "/tmp/base/") {
			t.Errorf("For relative input, expected base dir prefix, got %q", result)
		}
		if test.input == "~/documents/file.txt" && result != filepath.Join(home, "documents/file.txt") {
			t.Errorf("For tilde input, expected %q, got %q", filepath.Join(home, "documents/file.txt"), result)
		}
	}
}
