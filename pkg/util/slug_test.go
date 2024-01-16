package util

import (
	"testing"

	"github.com/gosimple/slug"
)

func TestSlugify(t *testing.T) {
	testCases := []struct {
		input          string
		namingStyle    string
		expectedOutput string
	}{
		{"Hello World", "lowercase", "hello world"},
		{"Hello World", "slug", "hello-world"},
		{"Hello World/My Project", "slug", "hello-world/my-project"},
		{"Hello World", "name", "Hello World"},

		// Test case for empty naming style, defaulting to slug
		{"Hello World", "", slug.Make("Hello World")},
	}

	for _, tc := range testCases {
		t.Run(tc.namingStyle, func(t *testing.T) {
			result := Slugify(tc.input, tc.namingStyle)

			if result != tc.expectedOutput {
				t.Errorf("Expected %s, but got %s", tc.expectedOutput, result)
			}
		})
	}
}
