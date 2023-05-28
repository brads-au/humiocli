package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newRolesRemoveCmd() *cobra.Command {

	cmd := cobra.Command{
		Use:   "remove <roleName>",
		Short: "Remove a role [Root Only]",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			roleName := args[0]
			client := NewApiClient(cmd)

			err := client.Roles().Remove(roleName)
			exitOnError(cmd, err, "Error removing the role")

			fmt.Fprintf(cmd.OutOrStdout(), "Successfully removed role: %s\n", roleName)
		},
	}

	return &cmd
}
