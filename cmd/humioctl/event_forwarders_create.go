package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newEventForwardersCreateCmd() *cobra.Command {
	description := ""
	properties := ""
	enabled := true

	cmd := cobra.Command{
		Use:   "create <forwarderName> <topic> <bootstrap-server>",
		Short: "Create a event forwarder.",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			forwarderName := args[0]
			topic := args[1]
			bootstrapServer := args[2]
			client := NewApiClient(cmd)

			// Construct properties
			properties = "bootstrap.servers=" + bootstrapServer

			err := client.EventForwarders().Create(forwarderName, description, properties, topic, enabled)
			exitOnError(cmd, err, "Error creating event forwarder")

			fmt.Fprintf(cmd.OutOrStdout(), "Successfully created event forwarder %s\n", forwarderName)
		},
	}

	cmd.Flags().StringVar(&description, "description", description, "Sets an optional description")
	//cmd.Flags().BoolVar(&enable, "description",  "Sets a repository connection with the chosen filter.")

	return &cmd
}
