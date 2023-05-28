package main

import (
	"github.com/spf13/cobra"
)

func newSavedQueriesListCmd() *cobra.Command {
	detailed := false

	cmd := cobra.Command{
		Use:   "list <searchDomain>",
		Short: "List all saved queries in a search domain.",
		Long:  "Lists all saved queries in a search domain. A search domain can be either a Repository or a View.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			searchDomainName := args[0]
			client := NewApiClient(cmd)

			savedQueries, err := client.SavedQueries().List(searchDomainName)
			exitOnError(cmd, err, "Error fetching saved queries")

			if detailed != true {
				printSavedQueriesTable(cmd, savedQueries)
			} else {
				printSavedQueriesTableDetailed(cmd, savedQueries)
			}
		},
	}

	cmd.Flags().BoolVar(&detailed, "detail", detailed, "Show a detailed view, including QueryString and query timeframes.")

	return &cmd
}
