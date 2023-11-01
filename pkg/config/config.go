package config

type RepoSyncConfig struct {
	Servers []Server              `yaml:"servers"`
	Sources []RepoSource          `yaml:"sources"`
	Bundle  map[string]RepoBundle `yaml:"bundle"`
}

type Server struct {
	Server    string       `yaml:"url"`
	Type      string       `yaml:"type"`
	TargetDir string       `yaml:"target"`
	Auth      RepoSyncAuth `yaml:"auth"`
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
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type RepoBundleOptions struct {
	SourcePrefix string   `yaml:"source-prefix"`
	TargetPrefix string   `yaml:"target-prefix"`
	Extensions   []string `yaml:"extensions"`
}
