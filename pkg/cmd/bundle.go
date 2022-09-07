package cmd

import (
	"github.com/cidverse/reposync/pkg/clone"
	"github.com/cidverse/reposync/pkg/config"
	"github.com/gosimple/slug"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/thoas/go-funk"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	rootCmd.AddCommand(bundleCmd)
	bundleCmd.Flags().StringP("config", "s", "", "config file, default ./reposync.yaml")
}

var bundleCmd = &cobra.Command{
	Use: "bundle",
	Run: func(cmd *cobra.Command, args []string) {
		config, configErr := config.LoadConfig(cfg.ConfigFile)
		if configErr != nil {
			log.Fatal().Err(configErr).Str("file", configFile).Msg("failed to parse config file")
		}

		for key, bundle := range config.Bundle {
			log.Info().Str("bundle", key).Msg("updating bundle")

			// clone or update
			for _, s := range bundle.Sources {
				cacheDir := getBundleCacheDir(s)
				clone.FetchProject(s, cacheDir)
			}

			// merge content
			for _, s := range bundle.Sources {
				// bundle sources
				contentDir := getBundleCacheContentDir(s)
				var files []string
				fileWalkErr := filepath.WalkDir(contentDir, func(path string, d fs.DirEntry, err error) error {
					// check for directory skip
					if d.IsDir() {
						if d.Name() == ".git" || d.Name() == ".idea" {
							return filepath.SkipDir
						}
					} else {
						extension := getFileExtension(d.Name())
						if len(s.Bundle.Extensions) > 0 {
							if funk.ContainsString(s.Bundle.Extensions, extension) {
								files = append(files, path)
							}
						} else {
							files = append(files, path)
						}
					}

					return nil
				})
				if fileWalkErr != nil {
					log.Fatal().Err(fileWalkErr).Msg("failed to query repository file list")
				}

				// targetDir
				targetDir, _ := filepath.Abs(bundle.TargetDir)
				if s.Bundle.TargetPrefix != "" {
					targetDir = filepath.Join(targetDir, s.Bundle.TargetPrefix)
				}

				for _, file := range files {
					contentDir := getBundleCacheContentDir(s)
					relativeFile, _ := filepath.Rel(contentDir, file)
					targetFile := filepath.Join(targetDir, relativeFile)

					log.Debug().Str("source-file", relativeFile).Str("target-file", targetFile).Msg("copy file")
					copyFile(file, targetFile)
				}
				log.Info().Int("count", len(files)).Str("target-dir", targetDir).Msg("updated files in target directory")
			}
		}
	},
}

func getBundleCacheDir(source config.RepoSyncSource) string {
	currentDir, _ := os.Getwd()
	contentDir := filepath.Join(currentDir, ".bundle-cache", slug.Make(source.Url+source.Ref))
	return contentDir
}

func getBundleCacheContentDir(source config.RepoSyncSource) string {
	contentDir := getBundleCacheDir(source)
	if source.Bundle.SourcePrefix != "" {
		contentDir = filepath.Join(contentDir, source.Bundle.SourcePrefix)
	}

	return contentDir
}

func getFileExtension(file string) string {
	split := strings.SplitN(file, ".", 2)
	if len(split) == 2 {
		return "." + split[1]
	}

	return ""
}

func copyFile(source string, target string) error {
	// create folders if needed
	dirName := filepath.Dir(target)
	if _, err := os.Stat(dirName); err != nil {
		mkErr := os.MkdirAll(dirName, os.ModePerm)
		if mkErr != nil {
			return mkErr
		}
	}

	// copy file
	input, err := os.ReadFile(source)
	if err != nil {
		return err
	}
	err = os.WriteFile(target, input, 0644)
	if err != nil {
		return err
	}

	return nil
}
