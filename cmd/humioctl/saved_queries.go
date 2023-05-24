package main

import (
	"github.com/humio/cli/api"
	"github.com/humio/cli/cmd/internal/format"
	"github.com/spf13/cobra"
)

func newSavedQueriesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "saved-queries [flags]",
		Short: "Manage saved queries",
	}

	cmd.AddCommand(newSavedQueriesListCmd())
	cmd.AddCommand(newSavedQueriesShowCmd())
	cmd.AddCommand(newSavedQueriesCreateCmd())
	cmd.AddCommand(newSavedQueriesDeleteCmd())

	return cmd
}

func printSavedQueriesTable(cmd *cobra.Command, savedQueries *api.SearchDomain) {
	if len(savedQueries.SavedQueries) == 0 {
		return
	}

	var rows [][]format.Value
	for _, query := range savedQueries.SavedQueries {
		rows = append(rows, []format.Value{
			format.String(savedQueries.Name),
			format.String(query.Id),
			format.String(query.Name),
		})
	}

	printOverviewTable(cmd, []string{"SearchDomain", "QueryId", "QueryName"}, rows)
}

func printSavedQueriesTableDetailed(cmd *cobra.Command, savedQueries *api.SearchDomain) {
	if len(savedQueries.SavedQueries) == 0 {
		return
	}

	var rows [][]format.Value
	for _, query := range savedQueries.SavedQueries {
		rows = append(rows, []format.Value{
			format.String(savedQueries.Name),
			format.String(query.Id),
			format.String(query.Name),
			format.String(query.Query.QueryString),
			format.String(query.Query.Start),
			format.String(query.Query.End),
			format.Bool(query.Query.IsLive),
		})
	}

	printOverviewTable(cmd, []string{"SearchDomain", "QueryId", "QueryName", "QueryString", "Start", "End", "IsLive"}, rows)
}
