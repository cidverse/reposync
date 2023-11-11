package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func indexCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "index",
		Aliases: []string{},
		Long:    "Indexes all repositories in the provided directories, adding them to the state.\nThis allows the first run of mirror to move existing repositories instead of cloning them from scratch.",
		Run: func(cmd *cobra.Command, args []string) {
			// process directories
			for _, dir := range args {
				log.Info().Str("dir", dir).Msg("indexing directory")

				// TODO: implement
				// search for git repositories
				// setup platform and lookup repository information via API
				// add to state
			}
		},
	}
}
