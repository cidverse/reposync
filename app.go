package main

import (
	"github.com/cidverse/reposync/pkg/cmd"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

// Version will be set at build time
var Version string

// CommitHash will be set at build time
var CommitHash string

// RepositoryStatus will be set at build time
var RepositoryStatus string

// BuildAt will be set at build time
var BuildAt string

// Init Hook
func init() {
	// Output to Stderr to not pollute stdout redirects with logs
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Only log the warning severity or above.
	zerolog.SetGlobalLevel(zerolog.WarnLevel)

	// version information
	cmd.Version = Version
	cmd.CommitHash = CommitHash
	cmd.RepositoryStatus = RepositoryStatus
	cmd.BuildAt = BuildAt
}

// CLI Main Entrypoint
func main() {
	cmdErr := cmd.Execute()
	if cmdErr != nil {
		log.Fatal().Err(cmdErr).Msg("internal cli library error")
	}
}
