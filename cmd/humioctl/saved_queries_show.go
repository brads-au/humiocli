package main

import (
	"github.com/spf13/cobra"
)

func newSavedQueriesShowCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "show <searchDomain> <savedQuery>",
		Short: "Show details about a saved query in a search domain.",
		Long:  "Show details about a saved query in a search domain. A search domain can be either a Repository or a View.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			savedQuery := args[1]
			searchDomainName := args[0]
			client := NewApiClient(cmd)

			savedQueries, err := client.SavedQueries().Get(savedQuery, searchDomainName)
			exitOnError(cmd, err, "Error fetching saved queries")

			printSavedQueriesTableDetailed(cmd, savedQueries)
		},
	}

	return &cmd
}
