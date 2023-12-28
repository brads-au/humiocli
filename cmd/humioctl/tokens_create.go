package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newTokensCreateCmd() *cobra.Command {
	var viewName string

	cmd := &cobra.Command{
		Use:   "create [flags] <token-name> <token-type> <permissions...>",
		Short: "Create a token",
		Long: `Create an authorisation token for a View (or repository), Organization or System.

Token types: View, Organization or System
Permissions: Please see "humioctl permissions list"

Example:
$ humioctl tokens add test-token-humio View -s humio ReadAccess ChangeIngestTokens
		`,
		Args: cobra.MinimumNArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			tokenType := args[1]
			permissions := args[2:]
			client := NewApiClient(cmd)

			token, err := client.Tokens().Add(name, tokenType, permissions, viewName)
			exitOnError(cmd, err, "Error adding token")

			fmt.Fprintf(cmd.OutOrStdout(), "Created token %q with secret: %q\n", name, token)
		},
	}

	cmd.Flags().StringVarP(&viewName, "view", "s", "", "Add the token to which view/repository, accepts name or ID (required for view tokens only)")

	return cmd
}
