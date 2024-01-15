package util

import (
	"os"
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
