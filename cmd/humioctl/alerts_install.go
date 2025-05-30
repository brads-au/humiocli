// Copyright © 2020 Humio Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"

	"github.com/humio/cli/internal/api"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func newAlertsInstallCmd() *cobra.Command {
	var (
		filePath, url, name string
	)

	cmd := cobra.Command{
		Use:   "install [flags] <view>",
		Short: "Installs an alert in a view",
		Long: `Install an alert from a URL or from a local file.

The install command allows you to install alerts from a URL or from a local file, e.g.

  $ humioctl alerts install viewName --name alertName --url=https://example.com/acme/alert.yaml

  $ humioctl alerts install viewName --name alertName --file=./alert.yaml

  $ humioctl alerts install viewName --file=./alert.yaml
`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var content []byte
			var err error

			// Check that we got the right number of argument
			// if we only got <view> you must supply --file or --url.
			if l := len(args); l == 1 {
				if filePath != "" {
					content, err = getBytesFromFile(filePath)
				} else if url != "" {
					content, err = getBytesFromURL(url)
				} else {
					cmd.Printf("You must specify a path using --file or --url\n")
					os.Exit(1)
				}
			}
			exitOnError(cmd, err, "Failed to load the alert")

			client := NewApiClient(cmd)
			viewName := args[0]

			var alert api.Alert
			err = yaml.Unmarshal(content, &alert)
			exitOnError(cmd, err, "Alert format is invalid")

			if name != "" {
				alert.Name = name
			}

			_, err = client.Alerts().Add(viewName, &alert)
			exitOnError(cmd, err, "Error creating alert")

			fmt.Fprintln(cmd.OutOrStdout(), "Alert created")
		},
	}

	cmd.Flags().StringVar(&filePath, "file", "", "The local file path to the alert to install.")
	cmd.Flags().StringVar(&url, "url", "", "A URL to fetch the alert file from.")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Install the alert under a specific name, ignoring the `name` attribute in the alert file.")

	return &cmd
}
