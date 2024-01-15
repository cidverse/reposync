package util

import (
	"os"
)

func ResolvePath(path string) string {
	if len(path) >= 2 && path[:2] == "~/" {
		path = os.Getenv("HOME") + path[1:]
	}
	path = os.ExpandEnv(path)

	return path
}
