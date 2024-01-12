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

	// default values
	for _, s := range cfg.Servers {
		if s.Mirror.DefaultAction == "" {
			s.Mirror.DefaultAction = RuleActionInclude
		}
		if s.Mirror.Rules == nil {
			s.Mirror.Rules = []MirrorRule{}
		}
	}

	return &cfg, nil
}

func Load() (*RepoSyncConfig, error) {
	// allow overriding config file location via environment variable
	if os.Getenv("REPOSYNC_CONFIG") != "" {
		return loadConfig(os.Getenv("REPOSYNC_CONFIG"))
	}

	// check default locations
	file, fileErr := findFirstExistingConfigFile(fileLocations)
	if fileErr != nil {
		return nil, fmt.Errorf("no config file found, tried: %s", strings.Join(fileLocations, ", "))
	}

	return loadConfig(file)
}
