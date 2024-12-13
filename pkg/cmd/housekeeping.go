package cmd

import (
	"slices"

	"github.com/cidverse/reposync/pkg/config"
	"github.com/cidverse/reposync/pkg/repository"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func houseKeepingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "housekeeping",
		Aliases: []string{"hk"},
		Short:   `runs various housekeeping tasks for repositories (git gc, git prune, git fsck)`,
		Run: func(cmd *cobra.Command, args []string) {
			// flags
			jobs, err := cmd.Flags().GetStringArray("jobs")
			if err != nil {
				log.Fatal().Err(err).Msg("failed to parse jobs flag")
			}
			silent, err := cmd.Flags().GetBool("silent")
			if err != nil {
				log.Fatal().Err(err).Msg("failed to parse silent flag")
			}

			// state
			stateFile := config.StateFile()
			state, err := config.LoadState(stateFile)
			if err != nil {
				log.Fatal().Err(err).Str("file", configFile).Msg("failed to parse state file")
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

				// repack
				if slices.Contains(jobs, "repack") {
					repackErr := repository.RepackRepository(r.Directory, silent)
					if repackErr != nil {
						log.Error().Err(repackErr).Str("repository", r.Directory).Msg("failed to repack repository")
					}
				}

				// prune
				if slices.Contains(jobs, "prune") {
					pruneErr := repository.PruneRepository(r.Directory, silent)
					if pruneErr != nil {
						log.Error().Err(pruneErr).Str("repository", r.Directory).Msg("failed to fsck repository")
					}
				}

				// gc
				if slices.Contains(jobs, "gc") {
					gcErr := repository.GCRepository(r.Directory, silent)
					if gcErr != nil {
						log.Error().Err(gcErr).Str("repository", r.Directory).Msg("failed to fsck repository")
					}
				}

				// fsck
				if slices.Contains(jobs, "fsck") {
					fsckErr := repository.FsckRepository(r.Directory, silent)
					if fsckErr != nil {
						log.Error().Err(fsckErr).Str("repository", r.Directory).Msg("failed to fsck repository")
					}
				}

				// commit-graph
				if slices.Contains(jobs, "commit-graph") {
					commitGraphErr := repository.RegenerateCommitGraph(r.Directory, silent)
					if commitGraphErr != nil {
						log.Error().Err(commitGraphErr).Str("repository", r.Directory).Msg("failed to regenerate commit")
					}
				}

				log.Info().Str("repository", r.Directory).Msg("housekeeping completed")
			}
		},
	}

	cmd.Flags().StringArrayP("jobs", "j", []string{"repack", "prune", "gc", "fsck", "commit-graph"}, "jobs to run (repack, prune, gc, fsck, commit-graph)")
	cmd.Flags().BoolP("silent", "s", false, "silent (omit stdout/stderr output) from cli commands")

	return cmd
}
