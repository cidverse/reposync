package cmd

import (
	"fmt"
	"os"

	"github.com/cidverse/cidverseutils/core/clioutputwriter"
	"github.com/cidverse/reposync/pkg/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func listCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{},
		Short:   `list all repositories managed by reposync`,
		Run: func(cmd *cobra.Command, args []string) {
			format, _ := cmd.Flags().GetString("format")
			columns, _ := cmd.Flags().GetStringSlice("columns")

			// state
			stateFile := config.StateFile()
			state, err := config.LoadState(stateFile)
			if err != nil {
				log.Fatal().Err(err).Str("file", configFile).Msg("failed to parse state file")
			}

			// data
			data := clioutputwriter.TabularData{
				Headers: []string{"ID", "NAME", "DIR", "REMOTE"},
				Rows:    [][]interface{}{},
			}
			for _, repo := range state.Repositories {
				data.Rows = append(data.Rows, []interface{}{
					repo.ID,
					repo.Name,
					repo.Directory,
					repo.Remote,
				})
			}

			// filter columns
			if len(columns) > 0 {
				data = clioutputwriter.FilterColumns(data, columns)
			}

			// print
			err = clioutputwriter.PrintData(os.Stdout, data, clioutputwriter.Format(format))
			if err != nil {
				log.Fatal().Err(err).Msg("failed to print data")
				os.Exit(1)
			}
		},
	}

	cmd.Flags().StringP("format", "f", string(clioutputwriter.DefaultOutputFormat()), fmt.Sprintf("output format %s", clioutputwriter.SupportedOutputFormats()))
	cmd.Flags().StringSliceP("columns", "c", []string{}, "columns to display")

	return cmd
}
