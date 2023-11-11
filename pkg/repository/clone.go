package repository

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func CloneRepository(target string, remote string) error {
	// create parent dir
	if err := os.MkdirAll(filepath.Dir(target), 0755|os.ModeDir); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}

	// clone using git cli
	cmd := exec.Command("git", "clone", remote, target)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	return nil
}
