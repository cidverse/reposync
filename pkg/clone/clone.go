package clone

import (
	"github.com/cidverse/reposync/pkg/config"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
)

func FetchProject(source config.RepoSyncSource, target string) {
	if _, err := os.Stat(target); err != nil {
		log.Debug().Str("directory", target).Msg("target directory is empty, cloning project")

		_, cloneErr := git.PlainClone(target, false, &git.CloneOptions{
			URL:           source.Url,
			RemoteName:    "origin",
			ReferenceName: plumbing.ReferenceName(source.Ref),
			Progress:      os.Stdout,
		})
		if cloneErr != nil {
			log.Error().Err(cloneErr).Str("directory", target).Str("repo", source.Url).Msg("repository clone failed")
			return
		}

		log.Info().Str("directory", target).Str("repo", source.Url).Msg("repository cloned successfully")
	} else {
		if _, err := os.Stat(filepath.Join(target, ".git")); err != nil {
			log.Error().Str("directory", target).Str("repo", source.Url).Msg("repository can't be updated, not a repository!")
		} else {
			repo, repoOpenErr := git.PlainOpen(filepath.Join(target, ".git"))
			if repoOpenErr != nil {
				log.Warn().Err(repoOpenErr).Str("directory", target).Str("repo", source.Url).Msg("repository open failed")
				return
			}

			// fetch
			fetchErr := repo.Fetch(&git.FetchOptions{
				RemoteName: "origin",
				Progress:   os.Stdout,
			})
			if fetchErr != nil && fetchErr.Error() != "already up-to-date" {
				log.Warn().Err(fetchErr).Str("directory", target).Str("repo", source.Url).Msg("repository fetch failed")
				return
			}

			log.Info().Str("directory", target).Str("repo", source.Url).Msg("repository updated successfully")
		}
	}
}
