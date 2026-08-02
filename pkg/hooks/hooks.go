package hooks

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

// Event names.
const (
	EventProjectAdded = "project-added"
)

// Config defines command hooks that are executed when certain events occur.
type Config struct {
	// ProjectAdded is executed when a new project has been added.
	ProjectAdded Commands `yaml:"project-added"`
}

// ForEvent returns the commands registered for the given event name.
func (c *Config) ForEvent(event string) (Commands, error) {
	switch event {
	case EventProjectAdded:
		return c.ProjectAdded, nil
	default:
		return nil, fmt.Errorf("unknown hook event %q", event)
	}
}

// Commands is a list of commands. It can be defined in the config as a single command string or as a list of command strings.
type Commands []string

func (c *Commands) UnmarshalYAML(value *yaml.Node) error {
	var commands []string

	if value.Kind == yaml.ScalarNode && value.Tag != "!!null" {
		commands = []string{value.Value}
	} else {
		if err := value.Decode(&commands); err != nil {
			return err
		}
	}

	*c = commands
	return nil
}

func expandVariables(command string, vars map[string]string) string {
	for key, value := range vars {
		command = strings.ReplaceAll(command, "{{"+key+"}}", value)
	}

	return command
}

// Expanded returns the commands with the given variables applied.
func Expanded(commands Commands, vars map[string]string) []string {
	expanded := make([]string, len(commands))
	for i, command := range commands {
		expanded[i] = expandVariables(command, vars)
	}

	return expanded
}

// Execute runs the provided commands with the given variables applied.
func Execute(commands Commands, vars map[string]string) {
	for _, command := range Expanded(commands, vars) {
		log.Info().Str("command", command).Msg("executing hook")

		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.Command("cmd", "/C", command)
		} else {
			cmd = exec.Command("sh", "-c", command)
		}
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Error().Err(err).Str("command", command).Msg("failed to execute hook")
		}
	}
}
