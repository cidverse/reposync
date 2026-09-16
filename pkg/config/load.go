package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cidverse/reposync/pkg/util"
	"gopkg.in/yaml.v3"
)

var fileLocations = []string{
	"$XDG_CONFIG_HOME/reposync/config.yaml",
	"$HOME/.config/reposync/config.yaml",
	"$HOME/.reposync.yaml",
}

func readConfigFile(file string) (RepoSyncConfig, error) {
	var cfg RepoSyncConfig

	fileContent, fileReadErr := os.ReadFile(file)
	if fileReadErr != nil {
		return cfg, fileReadErr
	}

	yamlErr := yaml.Unmarshal(fileContent, &cfg)
	if yamlErr != nil {
		return cfg, yamlErr
	}

	return cfg, nil
}

func applyDefaults(cfg *RepoSyncConfig) {
	for i := range cfg.Servers {
		if cfg.Servers[i].Mirror.DefaultAction == "" {
			cfg.Servers[i].Mirror.DefaultAction = RuleActionInclude
		}
		if cfg.Servers[i].Mirror.Rules == nil {
			cfg.Servers[i].Mirror.Rules = []MirrorRule{}
		}
	}
}

func loadConfig(file string) (*RepoSyncConfig, error) {
	cfg, err := readConfigFile(file)
	if err != nil {
		return nil, err
	}

	// merge includes (resolved relative to the main config file)
	baseDir := filepath.Dir(file)
	for _, include := range cfg.Includes {
		include = util.ResolveFilePath(include, baseDir)

		// skip if the referenced file does not exist
		if _, statErr := os.Stat(include); os.IsNotExist(statErr) {
			continue
		}

		override, readErr := readConfigFile(include)
		if readErr != nil {
			return nil, readErr
		}

		// includes are only supported at the top level
		if len(override.Includes) > 0 {
			return nil, fmt.Errorf("includes are only supported in the main config file, found nested includes in %s", include)
		}

		// version must match across the main config and all included files
		if cfg.Version != 0 && override.Version != 0 && cfg.Version != override.Version {
			return nil, fmt.Errorf("config version mismatch: main config uses version %d, include %s uses version %d", cfg.Version, include, override.Version)
		}

		mergeConfig(&cfg, override)
	}

	applyDefaults(&cfg)

	return &cfg, nil
}

// mergeConfig merges an override config into the base config.
// Lists are appended, scalars/maps are overridden by the override file.
func mergeConfig(base *RepoSyncConfig, override RepoSyncConfig) {
	base.Servers = append(base.Servers, override.Servers...)
	base.Sources = append(base.Sources, override.Sources...)

	// version must be identical across files (enforced in loadConfig), keep the base value
	if override.Hooks.ProjectAdded != nil {
		base.Hooks.ProjectAdded = override.Hooks.ProjectAdded
	}

	for name, bundle := range override.Bundle {
		if base.Bundle == nil {
			base.Bundle = map[string]RepoBundle{}
		}
		base.Bundle[name] = bundle
	}
}

func Load() (*RepoSyncConfig, error) {
	// allow overriding config file location via environment variable
	if env := os.Getenv("REPOSYNC_CONFIG"); env != "" {
		file := util.ResolveFilePath(env, ".")
		return loadConfig(file)
	}

	// check default locations
	file, fileErr := findFirstExistingConfigFile(fileLocations)
	if fileErr != nil {
		return nil, fmt.Errorf("no config file found, tried: %s", strings.Join(fileLocations, ", "))
	}

	return loadConfig(file)
}
