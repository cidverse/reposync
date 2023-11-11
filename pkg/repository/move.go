package repository

import (
	"fmt"
	"os"
	"path/filepath"
)

// MoveRepository moves a vcs repository
func MoveRepository(source string, target string) error {
	// check if source dir exists
	if _, err := os.Stat(source); os.IsNotExist(err) {
		return fmt.Errorf("source directory does not exist")
	}

	// check if target dir already exists
	if _, err := os.Stat(target); err == nil {
		return fmt.Errorf("target directory already exists")
	}

	// create parent dir
	if err := os.MkdirAll(filepath.Dir(target), 0644); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}

	// move directory
	if err := os.Rename(source, target); err != nil {
		return fmt.Errorf("failed to move directory: %w", err)
	}

	return nil
}
