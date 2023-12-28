package main

import (
	"github.com/humio/cli/cmd/internal/format"
	"github.com/spf13/cobra"
)

func newPermissionsListCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "list [flags]",
		Short: "List valid role and token permissions",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			client := NewApiClient(cmd)

			permissionTypes := []string{
				"Permission", // ViewPermissions, used on views and repositories
				"OrganizationPermission",
				"SystemPermission",
			}

			var rows [][]format.Value

			for _, permissionType := range permissionTypes {
				permissions, err := client.Permissions().List(permissionType)
				exitOnError(cmd, err, "Error fetching permission list")

				for i := 0; i < len(permissions); i++ {
					permission := permissions[i]

					rows = append(rows, []format.Value{
						format.String(permissionType),
						format.String(permission.Name),
						format.String(permission.Description),
						format.Bool(permission.IsDeprecated),
					})
				}
			}

			printOverviewTable(cmd, []string{"Type", "Name", "Description", "IsDepreciated"}, rows)
		},
	}

	return &cmd
}
