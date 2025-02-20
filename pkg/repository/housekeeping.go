package repository

import (
	"fmt"
	"os"
	"os/exec"
)

// RepackRepository repacks all objects in the repository to optimize storage
func RepackRepository(target string, silent bool) error {
	cmd := exec.Command("git", "repack", "-a", "-d", "--write-bitmap-index")
	cmd.Dir = target
	if !silent {
		cmd.Stdout = os.Stdout
	}
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to repack repository: %w", err)
	}

	return nil
}

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

// RegenerateCommitGraph regenerates the commit graph for the repository
func RegenerateCommitGraph(target string, silent bool) error {
	cmd := exec.Command("git", "commit-graph", "write", "--reachable")
	cmd.Dir = target
	if !silent {
		cmd.Stdout = os.Stdout
	}
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to regenerate commit graph: %w", err)
	}

	return nil
}
