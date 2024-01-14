package repository

import (
	"fmt"
	"os"
	"os/exec"
)

// FsckRepository runs the git file system check
func FsckRepository(target string, silent bool) error {
	cmd := exec.Command("git", "fsck", "--full", "--unreachable", "--strict")
	cmd.Dir = target
	if !silent {
		cmd.Stdout = os.Stdout
	}
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to fsck: %w", err)
	}

	return nil
}

// GCRepository runs the git garbage collection
func GCRepository(target string, silent bool) error {
	cmd := exec.Command("git", "gc", "--auto")
	cmd.Dir = target
	if !silent {
		cmd.Stdout = os.Stdout
	}
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to gc: %w", err)
	}

	return nil
}

// PruneRepository prunes all unreachable objects from the repository
func PruneRepository(target string, silent bool) error {
	cmd := exec.Command("git", "prune", "--expire", "now")
	cmd.Dir = target
	if !silent {
		cmd.Stdout = os.Stdout
	}
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to prune: %w", err)
	}

	return nil
}
