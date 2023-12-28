package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newTokensDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [flags] <token-id>",
		Short: "Deletes a token.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			client := NewApiClient(cmd)

			err := client.Tokens().Delete(name)
			exitOnError(cmd, err, "Error removing token")

			fmt.Fprintf(cmd.OutOrStdout(), "Successfully removed token %q\n", name)
		},
	}

	return cmd
}
