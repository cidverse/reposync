package cmd

import (
	"github.com/cidverse/reposync/pkg/clone"
	"github.com/cidverse/reposync/pkg/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.Flags().StringP("config", "s", "", "config file, default ./reposync.yaml")
}

var syncCmd = &cobra.Command{
	Use:     "sync",
	Aliases: []string{"s"},
	Run: func(cmd *cobra.Command, args []string) {
		config, configErr := config.LoadConfig(cfg.ConfigFile)
		if configErr != nil {
			log.Fatal().Err(configErr).Str("file", configFile).Msg("failed to parse config file")
		}

		// clone sources
		for _, s := range config.Sources {
			log.Debug().Str("remote", s.Url).Str("remote-ref", s.Ref).Str("target", s.TargetDir).Msg("processing project")
			clone.FetchProject(s, s.TargetDir)
		}
	},
}
