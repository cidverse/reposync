package cmd

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/cidverse/go-vcsapp/pkg/platform/api"
	"github.com/cidverse/go-vcsapp/pkg/vcsapp"
	"github.com/cidverse/reposync/pkg/clone"
	"github.com/cidverse/reposync/pkg/config"
	"github.com/cidverse/reposync/pkg/repository"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func cloneCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "clone",
		Aliases: []string{},
		Short:   `clones all repositories from configured remote servers`,
		Run: func(cmd *cobra.Command, args []string) {
			// flags
			dryRun, err := cmd.Flags().GetBool("dry-run")
			if err != nil {
				log.Fatal().Err(err).Msg("failed to parse dry-run flag")
			}
			silent, err := cmd.Flags().GetBool("silent")
			if err != nil {
				log.Fatal().Err(err).Msg("failed to parse silent flag")
			}

			// config
			c, err := config.Load()
			if err != nil {
				log.Fatal().Err(err).Str("file", cfg.ConfigFile).Msg("failed to parse config file")
			}

			// state
			stateFile := config.StateFile()
			state, err := config.LoadState(stateFile)
			if err != nil {
				log.Fatal().Err(err).Str("file", cfg.ConfigFile).Msg("failed to parse state file")
			}
			defer func(state *config.SyncState) { // ensure state is updated on exit
				saveErr := config.SaveState(stateFile, state)
				if saveErr != nil {
					log.Fatal().Err(saveErr).Msg("failed to save state")
				}
			}(state)

			// servers
			for _, s := range c.Servers {
				// skip if no local dir is specified
				if s.Mirror.LocalDir == "" {
					log.Debug().Str("server", s.Server).Str("type", s.Type).Msg("no local dir specified, skipping")
					continue
				}

				// setup platform
				log.Info().Str("server", s.Server).Str("type", s.Type).Msg("querying server")
				platformAuth, platformAuthErr := config.AuthToPlatformConfig(s.Type, s.Server, s.Auth)
				if platformAuthErr != nil {
					log.Error().Err(platformAuthErr).Msg("failed to initialize platform auth")
					continue
				}
				platform, platformErr := vcsapp.NewPlatform(platformAuth)
				if platformErr != nil {
					log.Error().Err(platformErr).Msg("failed to initialize platform")
					continue
				}

				// query repositories
				repos, repoErr := platform.Repositories(api.RepositoryListOpts{IncludeBranches: false, IncludeCommitHash: false})
				if repoErr != nil {
					log.Fatal().Err(repoErr).Msg("failed to list repositories")
				}
				log.Info().Int("count", len(repos)).Str("server", s.Server).Msg("received repository list")

				// process repositories
				for _, r := range repos {
					uniqueId := fmt.Sprintf("%s/%d", r.PlatformId, r.Id)
					log.Info().Str("namespace", r.Namespace).Str("repo", r.Name).Msg("processing repository")

					// check rules
					if config.EvaluateRules(s.Mirror.Rules, s.Mirror.DefaultAction, r) == config.RuleActionExclude {
						log.Debug().Str("namespace", r.Namespace).Str("repo", r.Name).Msg("repository excluded by rules, skipping")
						continue
					}

					// remote
					remote := r.CloneURL
					if s.Mirror.CloneMethod == config.CloneMethodSSH {
						remote = r.CloneSSH
					}

					// expected state
					expectedState := config.RepositoryState{
						ID:        uniqueId,
						Namespace: r.Namespace,
						Name:      r.Name,
						Remote:    remote,
						Directory: filepath.Join(s.Mirror.LocalDir, r.Namespace, r.Name),
						LastSync:  time.Now(),
					}

					// current state
					currentState, inCurrentState := state.Repositories[uniqueId]

					// run actions based on state
					if !inCurrentState {
						log.Debug().Str("repo", r.Namespace+"/"+r.Name).Str("current-dir", currentState.Directory).Str("expected-dir", expectedState.Directory).Msg("repository not present, cloning")
						if dryRun {
							continue
						}

						// add untracked repository in expected location, add to state (can occur if clone is interrupted and the state is not updated)
						if repository.Exists(expectedState.Directory) {
							log.Info().Str("repo", r.Namespace+"/"+r.Name).Str("dir", expectedState.Directory).Msg("adding existing repository to state")
							state.Repositories[uniqueId] = expectedState
							continue
						}

						// clone repository
						cloneErr := repository.CloneRepository(expectedState.Directory, remote, silent)
						if cloneErr != nil {
							log.Error().Err(cloneErr).Str("repo", r.Namespace+"/"+r.Name).Msg("failed to clone repository")
							continue
						}
					} else if currentState.Directory != expectedState.Directory {
						log.Debug().Str("repo", r.Namespace+"/"+r.Name).Str("current-dir", currentState.Directory).Str("expected-dir", expectedState.Directory).Msg("repository present in different location, moving")
						if dryRun {
							continue
						}

						// move repository
						moveErr := repository.MoveRepository(currentState.Directory, expectedState.Directory)
						if moveErr != nil {
							log.Error().Err(moveErr).Str("repo", r.Namespace+"/"+r.Name).Str("current-dir", currentState.Directory).Str("expected-dir", expectedState.Directory).Msg("failed to move repository")
							continue
						}
					} else if currentState.Directory == expectedState.Directory {
						log.Debug().Str("repo", r.Namespace+"/"+r.Name).Str("dir", expectedState.Directory).Msg("repository already present in expected location")
					}

					// update remote (allows switching between https/ssh for all repositories)
					updateRemoteErr := repository.UpdateRemote(expectedState.Directory, remote, silent)
					if updateRemoteErr != nil {
						log.Error().Err(updateRemoteErr).Str("repo", r.Namespace+"/"+r.Name).Msg("failed to update remote")
						continue
					}

					// add to state
					state.Repositories[uniqueId] = expectedState
				}
			}

			// clone sources
			for _, s := range c.Sources {
				log.Debug().Str("remote", s.Url).Str("remote-ref", s.Ref).Str("target", s.TargetDir).Msg("processing project")
				clone.FetchProject(s, s.TargetDir)
			}
		},
	}

	cmd.PersistentFlags().BoolP("dry-run", "d", false, "dry run")
	cmd.PersistentFlags().BoolP("silent", "s", false, "silent (omit stdout/stderr output) from cli commands")

	return cmd
}
