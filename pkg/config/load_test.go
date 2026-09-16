package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempFile(t *testing.T, dir string, name string, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write %s: %v", name, err)
	}
	return path
}

func TestLoadConfigWithIncludes(t *testing.T) {
	dir := t.TempDir()

	writeTempFile(t, dir, "config.yaml", `
version: 1
servers:
  - url: https://github.com
    type: github
includes:
  - team.yaml
  - missing.yaml
`)

	// override file adds a server; missing.yaml must be skipped silently
	writeTempFile(t, dir, "team.yaml", `
servers:
  - url: https://gitlab.example.com
    type: gitlab
    mirror:
      dir: /tmp/gitlab
`)

	cfg, err := loadConfig(filepath.Join(dir, "config.yaml"))
	if err != nil {
		t.Fatalf("loadConfig failed: %v", err)
	}

	if len(cfg.Servers) != 2 {
		t.Fatalf("expected 2 servers (base + override), got %d", len(cfg.Servers))
	}
	if cfg.Servers[0].Server != "https://github.com" {
		t.Errorf("expected base server first, got %q", cfg.Servers[0].Server)
	}
	if cfg.Servers[1].Server != "https://gitlab.example.com" {
		t.Errorf("expected override server appended, got %q", cfg.Servers[1].Server)
	}
	// defaults still applied after merge
	if cfg.Servers[1].Mirror.DefaultAction != RuleActionInclude {
		t.Errorf("expected default action to be applied to override server, got %q", cfg.Servers[1].Mirror.DefaultAction)
	}
}

func TestLoadConfigIncludesMissingFileSkipped(t *testing.T) {
	dir := t.TempDir()

	writeTempFile(t, dir, "config.yaml", `
includes:
  - definitely-not-here.yaml
servers:
  - url: https://github.com
    type: github
`)

	cfg, err := loadConfig(filepath.Join(dir, "config.yaml"))
	if err != nil {
		t.Fatalf("loadConfig failed: %v", err)
	}
	if len(cfg.Servers) != 1 {
		t.Fatalf("expected only the base server, got %d", len(cfg.Servers))
	}
}

func TestLoadConfigNestedIncludesRejected(t *testing.T) {
	dir := t.TempDir()

	writeTempFile(t, dir, "config.yaml", `
includes:
  - team.yaml
`)
	writeTempFile(t, dir, "team.yaml", `
includes:
  - sub.yaml
servers:
  - url: https://github.com
    type: github
`)
	writeTempFile(t, dir, "sub.yaml", `
sources:
  - url: https://github.com/cidverse/go-rules
`)

	if _, err := loadConfig(filepath.Join(dir, "config.yaml")); err == nil {
		t.Fatal("expected error for nested includes, got nil")
	}
}

func TestLoadConfigIncludesTildeAndEnv(t *testing.T) {
	dir := t.TempDir()
	subdir := filepath.Join(dir, "sub")
	if err := os.MkdirAll(subdir, 0755); err != nil {
		t.Fatalf("failed to create subdir: %v", err)
	}

	t.Setenv("HOME", dir)
	t.Setenv("REPOSYNC_TEAM_FILE", "~/team.yaml")
	writeTempFile(t, dir, "team.yaml", `
sources:
  - url: https://github.com/cidverse/go-rules
`)

	writeTempFile(t, subdir, "config.yaml", `
includes:
  - ~/team.yaml
  - ${REPOSYNC_TEAM_FILE}
`)

	cfg, err := loadConfig(filepath.Join(subdir, "config.yaml"))
	if err != nil {
		t.Fatalf("loadConfig failed: %v", err)
	}
	if len(cfg.Sources) != 2 {
		t.Fatalf("expected 2 sources (both includes), got %d", len(cfg.Sources))
	}
}


func TestLoadConfigIncludesRelativeToMainFile(t *testing.T) {
	dir := t.TempDir()
	subdir := filepath.Join(dir, "configs")
	if err := os.MkdirAll(subdir, 0755); err != nil {
		t.Fatalf("failed to create subdir: %v", err)
	}

	writeTempFile(t, dir, "configs/config.yaml", `
includes:
  - ../shared.yaml
`)
	writeTempFile(t, dir, "shared.yaml", `
sources:
  - url: https://github.com/cidverse/go-rules
    target: /tmp/go-rules
`)

	cfg, err := loadConfig(filepath.Join(subdir, "config.yaml"))
	if err != nil {
		t.Fatalf("loadConfig failed: %v", err)
	}

	if len(cfg.Sources) != 1 {
		t.Fatalf("expected 1 source from included file, got %d", len(cfg.Sources))
	}
	if cfg.Sources[0].Url != "https://github.com/cidverse/go-rules" {
		t.Errorf("unexpected source url %q", cfg.Sources[0].Url)
	}
}

func TestMergeConfigOverridesHooksAndScalars(t *testing.T) {
	dir := t.TempDir()

	writeTempFile(t, dir, "config.yaml", `
version: 1
includes:
  - override.yaml
hooks:
  project-added:
    - echo base
`)
	writeTempFile(t, dir, "override.yaml", `
version: 1
hooks:
  project-added:
    - echo override
`)

	cfg, err := loadConfig(filepath.Join(dir, "config.yaml"))
	if err != nil {
		t.Fatalf("loadConfig failed: %v", err)
	}

	if cfg.Version != 1 {
		t.Errorf("expected version to remain 1, got %d", cfg.Version)
	}
	if len(cfg.Hooks.ProjectAdded) != 1 || cfg.Hooks.ProjectAdded[0] != "echo override" {
		t.Errorf("expected hooks to be overridden by include, got %v", cfg.Hooks.ProjectAdded)
	}
}

func TestLoadConfigVersionMismatchRejected(t *testing.T) {
	dir := t.TempDir()

	writeTempFile(t, dir, "config.yaml", `
version: 1
includes:
  - override.yaml
`)
	writeTempFile(t, dir, "override.yaml", `
version: 2
servers:
  - url: https://github.com
    type: github
`)

	if _, err := loadConfig(filepath.Join(dir, "config.yaml")); err == nil {
		t.Fatal("expected error for version mismatch, got nil")
	}
}

func TestLoadConfigVersionUnsetInIncludeAllowed(t *testing.T) {
	dir := t.TempDir()

	writeTempFile(t, dir, "config.yaml", `
version: 1
includes:
  - override.yaml
`)
	// include omits version -> treated as same version
	writeTempFile(t, dir, "override.yaml", `
sources:
  - url: https://github.com/cidverse/go-rules
`)

	cfg, err := loadConfig(filepath.Join(dir, "config.yaml"))
	if err != nil {
		t.Fatalf("loadConfig failed: %v", err)
	}
	if len(cfg.Sources) != 1 {
		t.Fatalf("expected 1 source from include, got %d", len(cfg.Sources))
	}
}
