package main

import (
	"github.com/spf13/cobra"
)

func newSavedQueriesDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <searchDomain> <savedQuery>",
		Short: "Delete a saved query in a search domain.",
		Long:  "Delete a saved query in a search domain. A search domain can be either a Repository or a View.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			savedQuery := args[1]
			searchDomain := args[0]
			client := NewApiClient(cmd)

			err := client.SavedQueries().Delete(savedQuery, searchDomain)
			exitOnError(cmd, err, "Error deleting saved query.")

			cmd.Printf("Successfully deleted saved query with id %s in %s\n", savedQuery, searchDomain)
		},
	}
}
