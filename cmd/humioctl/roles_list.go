package main

import (
	"github.com/humio/cli/cmd/internal/format"
	"github.com/spf13/cobra"
)

func newRolesListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Lists all roles. [Root Only]",
		Run: func(cmd *cobra.Command, args []string) {
			client := NewApiClient(cmd)

			roles, err := client.Roles().List()
			exitOnError(cmd, err, "Error fetching role list")

			rows := make([][]format.Value, len(roles))
			for i, user := range roles {
				rows[i] = []format.Value{
					format.String(user.DisplayName),
					format.String(user.ID),
				}
			}

			printOverviewTable(cmd, []string{"DisplayName", "ID"}, rows)
		},
	}
}
