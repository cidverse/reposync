package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

var fileLocations = []string{
	"$XDG_CONFIG_HOME/reposync/config.yaml",
	"$HOME/.config/reposync/config.yaml",
	"$HOME/.reposync.yaml",
}

func findFirstExistingConfigFile() (string, error) {
	for _, file := range fileLocations {
		// resolve env vars
		file = os.ExpandEnv(file)

		// skip if path contains unresolved env vars
		if strings.Contains(file, "$") {
			continue
		}

		// check if file exists
		if _, err := os.Stat(file); err == nil {
			return file, nil
		}
	}

	return "", os.ErrNotExist
}

func loadConfig(file string) (*RepoSyncConfig, error) {
	cfg := RepoSyncConfig{}

	fileContent, fileReadErr := os.ReadFile(file)
	if fileReadErr != nil {
		return nil, fileReadErr
	}

	yamlErr := yaml.Unmarshal(fileContent, &cfg)
	if yamlErr != nil {
		return nil, yamlErr
	}

	return &cfg, nil
}

func Load() (*RepoSyncConfig, error) {
	file, fileErr := findFirstExistingConfigFile()
	if fileErr != nil {
		return nil, fmt.Errorf("no config file found, tried: %s", strings.Join(fileLocations, ", "))
	}

	return loadConfig(file)
}
