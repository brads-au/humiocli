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
	"strconv"

	"github.com/spf13/cobra"
)

func newReposUpdateCmd() *cobra.Command {
	var allowDataDeletionFlag bool
	var descriptionFlag, automaticSearchFlag, defaultQueryFlag stringPtrFlag
	var retentionTimeFlag, ingestSizeBasedRetentionFlag, storageSizeBasedRetentionFlag float64PtrFlag

	cmd := cobra.Command{
		Use:   "update [flags] <repo>",
		Short: "Updates the settings of a repository",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			repoName := args[0]
			client := NewApiClient(cmd)

			if descriptionFlag.value == nil && retentionTimeFlag.value == nil && ingestSizeBasedRetentionFlag.value == nil && storageSizeBasedRetentionFlag.value == nil && automaticSearchFlag.value == nil && defaultQueryFlag.value == nil {
				exitOnError(cmd, fmt.Errorf("you must specify at least one flag to update"), "Nothing specified to update")
			}

			if descriptionFlag.value != nil {
				err := client.Repositories().UpdateDescription(repoName, *descriptionFlag.value)
				exitOnError(cmd, err, "Error updating repository description")
			}
			if retentionTimeFlag.value != nil {
				err := client.Repositories().UpdateTimeBasedRetention(repoName, *retentionTimeFlag.value, allowDataDeletionFlag)
				exitOnError(cmd, err, "Error updating repository retention time in days")
			}
			if ingestSizeBasedRetentionFlag.value != nil {
				err := client.Repositories().UpdateIngestBasedRetention(repoName, *ingestSizeBasedRetentionFlag.value, allowDataDeletionFlag)
				exitOnError(cmd, err, "Error updating repository ingest size based retention")
			}
			if storageSizeBasedRetentionFlag.value != nil {
				err := client.Repositories().UpdateStorageBasedRetention(repoName, *storageSizeBasedRetentionFlag.value, allowDataDeletionFlag)
				exitOnError(cmd, err, "Error updating repository storage size based retention")
			}
			if automaticSearchFlag.value != nil {
				// Convert string to bool
				automaticSearchBool, errBool := strconv.ParseBool(*automaticSearchFlag.value)
				if errBool != nil {
					exitOnError(cmd, errBool, "Error, unable to convert automatic search value to bool.")
				}

				err := client.Repositories().UpdateAutomaticSearch(repoName, automaticSearchBool)
				exitOnError(cmd, err, "Error setting automatic search")
			}
			if defaultQueryFlag.value != nil {
				query, sqErr := client.SavedQueries().Get(*defaultQueryFlag.value, repoName)
				if sqErr != nil {
					exitOnError(cmd, sqErr, "Error setting default saved query")
				}
				err := client.Repositories().UpdateDefaultSavedQuery(repoName, query.SavedQueries[0].Id)
				exitOnError(cmd, err, "Error setting default saved query")
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Successfully updated repository %q\n", repoName)
		},
	}

	cmd.Flags().BoolVar(&allowDataDeletionFlag, "allow-data-deletion", false, "Allow changing retention settings for a non-empty repository")
	cmd.Flags().Var(&descriptionFlag, "description", "The description of the repository.")
	cmd.Flags().Var(&retentionTimeFlag, "retention-time", "The retention time in days for the repository.")
	cmd.Flags().Var(&ingestSizeBasedRetentionFlag, "ingest-size-based-retention", "The ingest size based retention for the repository.")
	cmd.Flags().Var(&storageSizeBasedRetentionFlag, "storage-size-based-retention", "The storage size based retention for the repository.")
	cmd.Flags().Var(&automaticSearchFlag, "automatic-search", "Set automatic search on loading the search page. true or false.")
	cmd.Flags().Var(&defaultQueryFlag, "default-query", "Set the default saved query to be used on the search page.")

	return &cmd
}
