package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

func LoadConfig(file string) (*RepoSyncConfig, error) {
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
