package cmd

import (
	"github.com/cidverse/reposync/pkg/clone"
	"github.com/cidverse/reposync/pkg/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func syncCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "sync",
		Aliases: []string{"s"},
		Run: func(cmd *cobra.Command, args []string) {
			c, err := config.Load()
			if err != nil {
				log.Fatal().Err(err).Str("file", cfg.ConfigFile).Msg("failed to parse config file")
			}

			// clone sources
			for _, s := range c.Sources {
				log.Debug().Str("remote", s.Url).Str("remote-ref", s.Ref).Str("target", s.TargetDir).Msg("processing project")
				clone.FetchProject(s, s.TargetDir)
			}
		},
	}
}
