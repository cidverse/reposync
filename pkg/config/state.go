package config

import (
	"encoding/json"
	"os"
	"time"
)

type SyncState struct {
	Repositories map[string]RepositoryState `json:"repositories"`
}

type RepositoryState struct {
	ID        string    `json:"id"`
	Namespace string    `json:"namespace"`
	Name      string    `json:"name"`
	Remote    string    `json:"remote"`
	Directory string    `json:"directory"`
	LastSync  time.Time `json:"last_sync"`
}

func StateFile() string {
	if os.Getenv("XDG_CONFIG_HOME") != "" {
		return os.ExpandEnv("$XDG_CONFIG_HOME/reposync/state.json")
	}
	if os.Getenv("HOME") != "" {
		return os.ExpandEnv("$HOME/.config/reposync/state.json")
	}

	// fallback to current directory
	return "state.json"
}

func LoadState(file string) (*SyncState, error) {
	s := &SyncState{
		Repositories: map[string]RepositoryState{},
	}

	// if file does not exist, return empty state
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return s, nil
	}

	// read file
	data, err := os.ReadFile(file)
	if err != nil {
		return s, err
	}

	// unmarshal
	if err := json.Unmarshal(data, &s); err != nil {
		return s, err
	}

	return s, nil
}

func SaveState(file string, state *SyncState) error {
	// ensure path exists
	err := createParentDir(file)
	if err != nil {
		return err
	}

	// write to file
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(file, data, 0644); err != nil {
		return err
	}

	return nil
}
