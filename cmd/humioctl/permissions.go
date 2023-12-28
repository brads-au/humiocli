package main

import (
	"github.com/spf13/cobra"
)

func newPermissionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "permissions",
		Short: "List valid role and token permissions",
	}

	cmd.AddCommand(newPermissionsListCmd())

	return cmd
}
