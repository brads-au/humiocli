package main

import (
	"strings"
	"time"

	"github.com/humio/cli/cmd/internal/format"
	"github.com/spf13/cobra"
)

func newTokensListCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "list [flags]",
		Short: "List authentication tokens.",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			client := NewApiClient(cmd)

			tokens, err := client.Tokens().List()
			exitOnError(cmd, err, "Error fetching token list")

			var rows [][]format.Value
			for i := 0; i < len(tokens); i++ {
				token := tokens[i]

				var expireTime string
				if token.ExpireAt > 0 {
					et := time.UnixMilli(token.ExpireAt)
					expireTime = et.Format(time.RFC3339)
				} else {
					expireTime = "not set"
				}

				var status string
				currentTime := time.Now().UTC()
				unixTime := currentTime.UnixMilli()
				if token.ExpireAt == 0 {
					status = "ok"
				} else if token.ExpireAt < unixTime {
					status = "expired"
				} else {
					status = "ok"
				}

				var views string
				if strings.Compare(token.Type, "ViewPermissionToken") == 1 {
					views = strings.Join(token.Views, ",")
				}

				rows = append(rows, []format.Value{
					format.String(token.Id),
					format.String(token.Name),
					format.String(token.Type),
					format.String(views),
					format.String(status),
					format.String(expireTime),
				})
			}

			printOverviewTable(cmd, []string{"Id", "Name", "Type", "Views", "Status", "Expire At"}, rows)
		},
	}

	return &cmd
}
