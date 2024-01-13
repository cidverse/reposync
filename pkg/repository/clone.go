package repository

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func Exists(target string) bool {
	gitDir := filepath.Join(target, ".git")
	_, err := os.Stat(gitDir)
	return err == nil || !os.IsNotExist(err)
}

func CloneRepository(target string, remote string, silent bool) error {
	// create parent dir
	if err := os.MkdirAll(filepath.Dir(target), 0755|os.ModeDir); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}

	// clone using git cli
	cmd := exec.Command("git", "clone", remote, target)
	if !silent {
		cmd.Stdout = os.Stdout
	}
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	return nil
}

func FetchRepository(target string, silent bool) error {
	// fetch using git cli
	cmd := exec.Command("git", "fetch")
	cmd.Dir = target
	if !silent {
		cmd.Stdout = os.Stdout
	}
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to fetch: %w", err)
	}

	return nil
}

func PullRepository(target string, silent bool) error {
	// pull using git cli
	cmd := exec.Command("git", "pull")
	cmd.Dir = target
	if !silent {
		cmd.Stdout = os.Stdout
	}
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to pull: %w", err)
	}

	return nil
}

func UpdateRemote(target string, remote string, silent bool) error {
	// update remote using git cli
	cmd := exec.Command("git", "remote", "set-url", "origin", remote)
	cmd.Dir = target
	if !silent {
		cmd.Stdout = os.Stdout
	}
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to update remote: %w", err)
	}

	return nil
}
