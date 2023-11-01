package clone

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/cidverse/reposync/pkg/config"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/rs/zerolog/log"
)

func FetchProject(source config.RepoSource, target string) {
	if _, err := os.Stat(target); err != nil {
		log.Debug().Str("directory", target).Msg("target directory is empty, cloning project")
		_, cloneErr := git.PlainClone(target, false, &git.CloneOptions{
			URL:           source.Url,
			RemoteName:    "origin",
			ReferenceName: plumbing.ReferenceName(source.Ref),
			Progress:      os.Stdout,
			Auth:          GetRepoAuth(source),
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
				Auth:       GetRepoAuth(source),
			})
			if fetchErr != nil && fetchErr.Error() != "already up-to-date" {
				log.Warn().Err(fetchErr).Str("directory", target).Str("repo", source.Url).Msg("repository fetch failed")
				return
			}

			// pull
			workTree, workTreeErr := repo.Worktree()
			if workTreeErr == nil {
				pullErr := workTree.Pull(&git.PullOptions{
					RemoteName:    "origin",
					ReferenceName: plumbing.ReferenceName(source.Ref),
					Progress:      os.Stdout,
					Auth:          GetRepoAuth(source),
					SingleBranch:  true,
				})
				if pullErr != nil {
					log.Warn().Err(pullErr).Str("directory", target).Str("repo", source.Url).Msg("repository pull failed")
					return
				}
			}

			log.Info().Str("directory", target).Str("repo", source.Url).Msg("repository updated successfully")
		}
	}
}

func GetRepoAuth(source config.RepoSource) transport.AuthMethod {
	var auth transport.AuthMethod

	// config file
	if source.Auth.Username != "" {
		auth = &http.BasicAuth{
			Username: source.Auth.Username,
			Password: source.Auth.Password,
		}
	}

	// env
	hostSlug := getHostnameSlug(source.Url)
	envUsername, _ := os.LookupEnv("REPOSYNC_" + hostSlug + "_USERNAME")
	envPassword, _ := os.LookupEnv("REPOSYNC_" + hostSlug + "_PASSWORD")
	if envUsername != "" && envPassword != "" {
		auth = &http.BasicAuth{
			Username: envUsername,
			Password: envPassword,
		}
	}

	return auth
}

func getHostnameSlug(input string) string {
	url, err := url.Parse(input)
	if err != nil {
		return ""
	}

	output := strings.TrimPrefix(url.Hostname(), "www.")
	output = strings.ReplaceAll(output, ".", "_")
	output = strings.ToUpper(output)
	return output
}
