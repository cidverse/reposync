package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// Version will be set at build time
var Version string

// CommitHash will be set at build time
var CommitHash string

// BuildAt will be set at build time
var BuildAt string

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:  "version",
	Long: "Prints the version number of reposync",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("reposync v" + Version + "-" + CommitHash + " " + runtime.GOOS + "/" + runtime.GOARCH + " BuildDate=" + BuildAt)
	},
}
