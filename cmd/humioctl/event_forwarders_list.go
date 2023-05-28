package main

import (
	"github.com/humio/cli/cmd/internal/format"
	"github.com/spf13/cobra"
)

func newEventForwardersListCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "list",
		Short: "List all event forwarders.",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			client := NewApiClient(cmd)

			eventForwarders, err := client.EventForwarders().List()
			exitOnError(cmd, err, "Error fetching event forwarders")

			//cmd.Printf("%+v", eventForwarders)

			rows := make([][]format.Value, len(eventForwarders))
			for i, eventForwarder := range eventForwarders {
				rows[i] = []format.Value{
					format.String(eventForwarder.Name),
					format.String(eventForwarder.ID),
					format.String(eventForwarder.ForwarderInfo.Topic),
					yesNo(eventForwarder.ForwarderInfo.Enabled),
				}
			}

			printOverviewTable(cmd, []string{"Name", "ID", "Topic", "Enabled"}, rows)
		},
	}

	return &cmd
}
