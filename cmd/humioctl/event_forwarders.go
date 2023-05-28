package main

import (
	"github.com/spf13/cobra"
)

func newEventForwardersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "event-forwarders [flags]",
		Short: "Manage event forwarders",
	}

	cmd.AddCommand(newEventForwardersListCmd())
	// cmd.AddCommand(newSavedQueriesShowCmd())
	cmd.AddCommand(newEventForwardersCreateCmd())
	// cmd.AddCommand(newSavedQueriesDeleteCmd())

	return cmd
}
