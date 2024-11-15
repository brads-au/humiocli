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
	"github.com/spf13/cobra"
)

func newAlertsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alerts",
		Short: "Manage alerts",
	}

	cmd.AddCommand(newAlertsListCmd())
	cmd.AddCommand(newAlertsInstallCmd())
	cmd.AddCommand(newAlertsExportCmd())
	cmd.AddCommand(newAlertsExportAllCmd())
	cmd.AddCommand(newAlertsRemoveCmd())
	cmd.AddCommand(newAlertsShowCmd())

	return cmd
}
