package util

import (
	"os"
	"path/filepath"
)

func ResolvePath(path string) string {
	if len(path) >= 2 && path[:2] == "~/" {
		path = os.Getenv("HOME") + path[1:]
	}
	path = os.ExpandEnv(path)

	return path
}

// ResolveFilePath resolves a path for use with config files: expands environment
// variables, expands a leading "~/", joins relative paths against baseDir and
// returns an absolute path.
func ResolveFilePath(path string, baseDir string) string {
	path = os.ExpandEnv(path)
	if len(path) >= 2 && path[:2] == "~/" {
		path = os.Getenv("HOME") + path[1:]
	}
	if !filepath.IsAbs(path) {
		path = filepath.Join(baseDir, path)
	}
	if abs, err := filepath.Abs(path); err == nil {
		path = abs
	}

	return path
}
