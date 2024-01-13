package cmd

import (
	"github.com/cidverse/reposync/pkg/config"
	"github.com/cidverse/reposync/pkg/repository"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func updateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update",
		Aliases: []string{"u"},
		Short:   `pulls all changes for repositories cloned by reposync`,
		Run: func(cmd *cobra.Command, args []string) {
			// flags
			silent, err := cmd.Flags().GetBool("silent")
			if err != nil {
				log.Fatal().Err(err).Msg("failed to parse silent flag")
			}

			// state
			stateFile := config.StateFile()
			state, err := config.LoadState(stateFile)
			if err != nil {
				log.Fatal().Err(err).Str("file", cfg.ConfigFile).Msg("failed to parse state file")
			}
			defer func(state *config.SyncState) { // ensure state is updated
				saveErr := config.SaveState(stateFile, state)
				if saveErr != nil {
					log.Fatal().Err(saveErr).Msg("failed to save state")
				}
			}(state)

			// iterate over all repositories
			for _, r := range state.Repositories {
				// check if repository still exists
				if !repository.Exists(r.Directory) {
					log.Warn().Str("repository", r.Directory).Msg("repository not found, removing from state")
					delete(state.Repositories, r.ID)
					continue
				}

				// fetch
				fetchErr := repository.FetchRepository(r.Directory, silent)
				if fetchErr != nil {
					log.Error().Err(fetchErr).Str("repository", r.Directory).Msg("failed to fetch repository")
				}

				// pull
				pullErr := repository.PullRepository(r.Directory, silent)
				if pullErr != nil {
					log.Error().Err(pullErr).Str("repository", r.Directory).Msg("failed to pull repository")
				}

				log.Info().Str("repository", r.Directory).Msg("repository updated")
			}
		},
	}

	cmd.PersistentFlags().BoolP("silent", "s", false, "silent (omit stdout/stderr output) from cli commands")

	return cmd
}
