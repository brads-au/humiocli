package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newSavedQueriesCreateCmd() *cobra.Command {
	start := "1h"
	end := "now"
	isLive := false
	widgetType := "list-view"

	cmd := cobra.Command{
		Use:   "create <searchDomain> <queryName> <queryString>",
		Short: "Create a saved query in a search domain.",
		Long: `Creates a saved query within a search domain with the provided arguements. A search domain can be either a Repository or a View.
Here's an example creating a saved query called "all-warnings" in the repo called "default-repo" with a simple filter of "loglevel = WARN" including events
within the last 15 minutes.
  $ humioctl saved-queries create default-repo all-warnings "loglevel = WARN" --start "15m"
`,
		Args: cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[1]
			viewName := args[0]
			queryString := args[2]
			client := NewApiClient(cmd)

			err := client.SavedQueries().Create(name, viewName, queryString, start, end, isLive, widgetType)
			exitOnError(cmd, err, "Error creating saved query")

			fmt.Fprintf(cmd.OutOrStdout(), "Successfully created saved query: %q on %q\n", name, viewName)

		},
	}

	cmd.Flags().StringVar(&start, "start", start, "Sets query start time.")
	cmd.Flags().StringVar(&end, "end", end, "Sets query end time.")
	cmd.Flags().BoolVar(&isLive, "isLive", isLive, "Sets query live mode.")
	cmd.Flags().StringVar(&widgetType, "widgetType", widgetType, "Sets the displayed widget type.")

	return &cmd
}
