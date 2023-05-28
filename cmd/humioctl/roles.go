package main

import (
	"github.com/humio/cli/api"
	"github.com/humio/cli/cmd/internal/format"
	"github.com/spf13/cobra"
)

func newRolesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "roles",
		Short: "Manage roles [Root Only]",
	}

	// TODO add roles

	// cmd.AddCommand(newRolesAddCmd())
	cmd.AddCommand(newRolesRemoveCmd())
	cmd.AddCommand(newRolesListCmd())
	cmd.AddCommand(newRolesShowCmd())

	return cmd
}

func printRolesDetailsTable(cmd *cobra.Command, role *api.Role) {
	details := [][]format.Value{
		{format.String("Name"), format.String(role.DisplayName)},
		//{format.String("Description"), format.String(role.Description)},
		//{format.String("Color"), format.String(role.Color)},
		//{format.String("View Permissions"), format.String(user.ViewPermissions)},
		{format.String("ID"), format.String(role.ID)},
	}

	printDetailsTable(cmd, details)
}
