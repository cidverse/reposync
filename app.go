package main

import (
	"os"

	"github.com/cidverse/reposync/pkg/cmd"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	status  = "clean"
)

// Init Hook
func init() {
	// Output to Stderr to not pollute stdout redirects with logs
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Only log the warning severity or above.
	zerolog.SetGlobalLevel(zerolog.WarnLevel)

	// version information
	cmd.Version = version
	cmd.CommitHash = commit
	cmd.RepositoryStatus = status
	cmd.BuildAt = date
}

// CLI Main Entrypoint
func main() {
	cmdErr := cmd.Execute()
	if cmdErr != nil {
		log.Fatal().Err(cmdErr).Msg("internal cli library error")
	}
}
