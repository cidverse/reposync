package hooks

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestCommandsUnmarshalYAML(t *testing.T) {
	tests := []struct {
		name     string
		yaml     string
		expected Commands
	}{
		{"single command", "project-added: zoxide add {{projectDir}}", Commands{"zoxide add {{projectDir}}"}},
		{"command list", "project-added:\n  - zoxide add {{projectDir}}\n  - echo {{projectName}}", Commands{"zoxide add {{projectDir}}", "echo {{projectName}}"}},
		{"empty", "project-added:", Commands{}},
		{"missing", "other: value", Commands{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var cfg Config
			err := yaml.Unmarshal([]byte(test.yaml), &cfg)
			if err != nil {
				t.Fatalf("unmarshal failed: %v", err)
			}

			if len(cfg.ProjectAdded) != len(test.expected) {
				t.Fatalf("expected %d commands, got %d", len(test.expected), len(cfg.ProjectAdded))
			}
			for i := range test.expected {
				if cfg.ProjectAdded[i] != test.expected[i] {
					t.Errorf("expected command %q, got %q", test.expected[i], cfg.ProjectAdded[i])
				}
			}
		})
	}
}

func TestExpandVariables(t *testing.T) {
	command := expandVariables("zoxide add {{projectDir}} && echo {{projectName}}", map[string]string{
		"projectDir":  "/tmp/github/my-org/my-project",
		"projectName": "my-project",
	})

	expected := "zoxide add /tmp/github/my-org/my-project && echo my-project"
	if command != expected {
		t.Errorf("expected %q, got %q", expected, command)
	}
}
