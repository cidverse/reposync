package cmd

import (
	"github.com/cidverse/reposync/pkg/config"
	"github.com/cidverse/reposync/pkg/hooks"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func hooksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "hooks",
		Aliases: []string{},
		Short:   `runs hooks for tracked projects`,
	}

	cmd.AddCommand(hooksRunCmd())

	return cmd
}

func hooksRunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run <event>",
		Short: `runs the hooks for the given event for all tracked projects`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			event := args[0]
			dryRun, err := cmd.Flags().GetBool("dry-run")
			if err != nil {
				log.Fatal().Err(err).Msg("failed to parse dry-run flag")
			}

			// config
			c, err := config.Load()
			if err != nil {
				log.Fatal().Err(err).Str("file", configFile).Msg("failed to parse config file")
			}

			// hooks for event
			commands, err := c.Hooks.ForEvent(event)
			if err != nil {
				log.Fatal().Err(err).Msg("invalid hook event")
			}

			// state
			stateFile := config.StateFile()
			state, err := config.LoadState(stateFile)
			if err != nil {
				log.Fatal().Err(err).Str("file", configFile).Msg("failed to parse state file")
			}

			// run hooks for all tracked projects
			for _, repo := range state.Repositories {
				vars := map[string]string{
					"projectDir":  repo.Directory,
					"projectName": repo.Name,
					"namespace":   repo.Namespace,
					"projectId":   repo.ID,
					"projectUrl":  repo.Remote,
				}

				if dryRun {
					for _, command := range hooks.Expanded(commands, vars) {
						log.Info().Str("command", command).Msg("would execute hook")
					}
					continue
				}

				hooks.Execute(commands, vars)
			}
		},
	}

	cmd.PersistentFlags().BoolP("dry-run", "d", false, "dry run")

	return cmd
}
