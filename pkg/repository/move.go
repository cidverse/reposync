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
	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}

	// check permissions
	if err := updatePermissions(filepath.Dir(source)); err != nil {
		return fmt.Errorf("failed to update source permissions: %w", err)
	}
	if err := updatePermissions(filepath.Dir(target)); err != nil {
		return fmt.Errorf("failed to update target permissions: %w", err)
	}

	// move directory
	if err := os.Rename(source, target); err != nil {
		return fmt.Errorf("failed to move directory: %w", err)
	}

	return nil
}

func updatePermissions(path string) error {
	// get path information
	fileInfo, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to get path information: %w", err)
	}
	if !fileInfo.IsDir() {
		return fmt.Errorf("path is not a directory")
	}

	// update permissions if necessary
	if fileInfo.Mode().Perm() != 0755 {
		chmodErr := os.Chmod(path, 0755)
		if chmodErr != nil {
			return fmt.Errorf("failed to update permissions: %w", chmodErr)
		}
	}

	return nil
}
