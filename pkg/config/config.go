package config

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/cidverse/go-rules/pkg/expr"
	"github.com/cidverse/go-vcsapp/pkg/platform/api"
	"github.com/cidverse/go-vcsapp/pkg/vcsapp"
	"github.com/rs/zerolog/log"
)

type RepoSyncConfig struct {
	Version int                   `yaml:"version"`
	Servers []Server              `yaml:"servers"`
	Sources []RepoSource          `yaml:"sources"`
	Bundle  map[string]RepoBundle `yaml:"bundle"`
}

type Server struct {
	Server    string       `yaml:"url"`
	Type      string       `yaml:"type"`
	TargetDir string       `yaml:"target"`
	Auth      RepoSyncAuth `yaml:"auth"`
	Mirror    MirrorOpts   `yaml:"mirror"`
}

type RepoSource struct {
	Url       string            `yaml:"url"`
	Ref       string            `yaml:"ref"`
	Group     []string          `yaml:"group"`
	TargetDir string            `yaml:"target"`
	Bundle    RepoBundleOptions `yaml:"bundle"`
	Auth      RepoSyncAuth      `yaml:"auth"`
}

type RepoBundle struct {
	TargetDir string       `yaml:"target"`
	Sources   []RepoSource `yaml:"sources"`
}

type RepoSyncAuth struct {
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	PasswordFile string `yaml:"password-file"`
}

type RepoBundleOptions struct {
	SourcePrefix string   `yaml:"source-prefix"`
	TargetPrefix string   `yaml:"target-prefix"`
	Extensions   []string `yaml:"extensions"`
}

type MirrorOpts struct {
	LocalDir      string       `yaml:"dir"`
	CloneMethod   CloneMethod  `yaml:"clone-method"`
	Rules         []MirrorRule `yaml:"rules"`
	DefaultAction RuleAction   `yaml:"default-action"`
}

type MirrorRule struct {
	Rule   string     `yaml:"rule"`
	Action RuleAction `yaml:"action"`
}

type CloneMethod string

const (
	CloneMethodHTTPS CloneMethod = "https"
	CloneMethodSSH   CloneMethod = "ssh"
)

type RuleAction string

const (
	RuleActionInclude RuleAction = "include"
	RuleActionExclude RuleAction = "exclude"
)

func AuthToPlatformConfig(serverType string, serverUrl string, auth RepoSyncAuth) (vcsapp.PlatformConfig, error) {
	// password file
	if auth.PasswordFile != "" {
		// resolve path
		auth.PasswordFile = strings.Replace(auth.PasswordFile, "~", os.Getenv("HOME"), 1)
		auth.PasswordFile = os.ExpandEnv(auth.PasswordFile)

		// read
		file, err := os.ReadFile(auth.PasswordFile)
		if err != nil {
			return vcsapp.PlatformConfig{}, fmt.Errorf("failed to read password file: %w", err)
		}
		auth.Password = string(file)
	}

	// password is required
	if auth.Password == "" {
		return vcsapp.PlatformConfig{}, fmt.Errorf("no password provided")
	}

	// platform config
	if serverType == "github" {
		return vcsapp.PlatformConfig{
			GitHubUsername: os.ExpandEnv(auth.Username),
			GitHubToken:    os.ExpandEnv(auth.Password),
		}, nil
	} else if serverType == "gitlab" {
		return vcsapp.PlatformConfig{
			GitLabServer:      os.ExpandEnv(serverUrl),
			GitLabAccessToken: os.ExpandEnv(auth.Password),
		}, nil
	}

	return vcsapp.PlatformConfig{}, fmt.Errorf("unsupported server type: %s", serverType)
}

func EvaluateRules(rules []MirrorRule, defaultAction RuleAction, repo api.Repository) RuleAction {
	ctx := map[string]interface{}{
		"uniqueId": fmt.Sprintf("%s/%d", repo.PlatformId, repo.Id),
		"id":       repo.Id,
		"group":    repo.Namespace,
		"path":     path.Join(repo.Namespace, repo.Name),
		"name":     repo.Name,
		"is_fork":  repo.IsFork,
	}

	for _, rule := range rules {
		match, err := expr.EvaluateRule(rule.Rule, ctx)
		if err != nil {
			log.Fatal().Err(err).Str("rule", rule.Rule).Msg("failed to evaluate rule, check your configuration file syntax")
		}
		if match {
			return rule.Action
		}
	}

	return defaultAction
}
