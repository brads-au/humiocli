package main

import (
	"os"
	"regexp"

	"github.com/humio/cli/api"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func newApplyCmd() *cobra.Command {
	var filePath, url string
	var dryRun, verbose bool
	var err error
	var content []byte

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	cmd := cobra.Command{
		Use:   "apply",
		Short: "Apply configuration [Root Only]",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			// Check that we got the right number of argument
			// if we only got <repo> you must supply --file or --url.
			if filePath != "" {
				content, err = getParserFromFile(filePath)
			} else if url != "" {
				content, err = getURLParser(url)
			} else {
				cmd.PrintErrf("If you only provide repo you must specify --file or --url\n")
				os.Exit(1)
			}
			exitOnError(cmd, err, "Failed to load the parser")

			client := NewApiClient(cmd)

			// Load the YAML
			config := Config{}
			err = yaml.Unmarshal(content, &config)
			exitOnError(cmd, err, "The yaml content is invalid.")

			// Loop through the apply types
			// Users
			for _, user := range config.Users {
				if user.Username != "" {
					rUser, _ := client.Users().Get(user.Username)

					// Create user if it doesn't exist
					if rUser.ID == "" {
						log.Info().Msgf("User creating: %s", user.Username)

						_, err := client.Users().Add(user.Username, api.UserChangeSet{})
						exitOnError(cmd, err, "Error creating the user")
					} else {
						log.Info().Msgf("User exists: %s", user.Username)
					}

					// Reconcile user permissions
					if user.SearchDomains != nil {
						for _, searchDomain := range user.SearchDomains {
							// If regex, make a list of search domains
							listSearchDomains := make([]string, 0)
							if searchDomain.Regex {
								resultSearchDomains, _ := client.SearchDomains().List()
								for _, result := range resultSearchDomains {
									match, _ := regexp.MatchString(searchDomain.Name, result.Name)
									if match {
										listSearchDomains = append(listSearchDomains, result.Name)
									}
								}
							} else {
								listSearchDomains = append(listSearchDomains, searchDomain.Name)
							}

							// Got a list, now reconcile permissions for each
							for _, result := range listSearchDomains {
								role, err := client.Roles().Get(searchDomain.Role)
								if err != nil {
									log.Error().Err(err).Msgf("Error getting role id for user: %s", user.Username)
									continue
								}

								sd, err := client.SearchDomains().Get(result)
								if err != nil {
									log.Error().Err(err).Msgf("Error getting search domain id for user: %s", user.Username)
									continue
								}

								// FIXME: Getting user ID, could I pick this up from above though
								user, _ := client.Users().Get(user.Username)

								errSD := client.SearchDomains().UpdateUserPermissions(sd.Id, user.ID, role.ID)
								if errSD != nil {
									log.Error().Err(err).Msgf("Error setting permissions for user: %s", user.Username)
									continue
								}

								log.Info().Msgf("User updated: %s on %s as %s", user.Username, sd.Name, role.DisplayName)
							}
						}
					}
				}
			}

			// Default Queries
			for _, defaultQuery := range config.DefaultQueries {
				if defaultQuery.Name != "" {
					if defaultQuery.Global {
						searchDomains, _ := client.SearchDomains().List()

						for _, searchDomain := range searchDomains {
							savedQuery, _ := client.SavedQueries().Get(defaultQuery.Name, searchDomain.Name)
							if savedQuery == nil {
								log.Info().Msgf("Query missing, creating: %s in %s", defaultQuery.Name, searchDomain.Name)
								err := client.SavedQueries().Create(defaultQuery.Name, searchDomain.Name, defaultQuery.QueryString, defaultQuery.Start, "now", false, "list-view")
								if err != nil {
									log.Error().Err(err).Msgf("Error creating query: %s in %s", defaultQuery.Name, searchDomain.Name)
									continue
								}
							} else {
								log.Info().Msgf("Query exists: %s in %s", defaultQuery.Name, searchDomain.Name)
							}
						}
					} else {
						// FIXME: I currenly only deploy global queries
					}
				}
			}

			// Repos
			listRepos, err := client.Repositories().List()
			if err != nil {
				log.Fatal().Err(err).Msg("Error listing repositories")
				exitOnError(cmd, err, "Error listing repositories")
			}

			for _, repo := range config.Repos {
				if repo.Regex {
					for _, rangeRepos := range listRepos {
						match, _ := regexp.MatchString(repo.Name, rangeRepos.Name)

						if match {
							log.Info().Msgf("Repo exists: %s", rangeRepos.Name)

							if repo.AutomaticSearch != rangeRepos.AutomaticSearch {
								log.Info().Msgf("Repo updating automatic search: %s", rangeRepos.Name)
								err := client.Repositories().UpdateAutomaticSearch(rangeRepos.Name, repo.AutomaticSearch)
								if err != nil {
									log.Error().Err(err).Msgf("Error setting automatic search: %s", rangeRepos.Name)
									continue
								}
							}

							if repo.DefaultQuery != rangeRepos.DefaultQuery.Name {
								log.Info().Msgf("Repo updating default query: %s", rangeRepos.Name)
								err := client.Repositories().UpdateDefaultSavedQuery(rangeRepos.Name, repo.DefaultQuery)
								if err != nil {
									log.Error().Err(err).Msgf("Error setting default saved query: %s", rangeRepos.Name)
									continue
								}
							}
						}
					}
				} else {
					if repo.Name != "" {
						result, _ := client.Repositories().Get(repo.Name)

						if result.ID == "" {
							log.Info().Msgf("Repo missing, creating (not yet implemented): %s", repo.Name)
							continue
						} else {
							log.Info().Msgf("Repo exists: %s", result.Name)

							if repo.AutomaticSearch != result.AutomaticSearch {
								log.Info().Msgf("Repo updating automatic search: %s", result.Name)
								err := client.Repositories().UpdateAutomaticSearch(repo.Name, repo.AutomaticSearch)
								if err != nil {
									log.Error().Err(err).Msgf("Error setting automatic search: %s", result.Name)
									continue
								}
							}

							if repo.DefaultQuery != result.DefaultQuery.Name {
								log.Info().Msgf("Repo updating default query: %s", result.Name)
								err := client.Repositories().UpdateDefaultSavedQuery(repo.Name, repo.DefaultQuery)
								if err != nil {
									log.Error().Err(err).Msgf("Error setting default saved query: %s", result.Name)
									continue
								}
							}
						}
					}
				}
			}

			// Views
			listViews, err := client.Views().List()
			if err != nil {
				log.Fatal().Err(err).Msg("Error listing views")
				exitOnError(cmd, err, "Error listing views")
			}

			for _, view := range config.Views {
				if view.Regex {
					for _, rangeViews := range listViews {
						match, _ := regexp.MatchString(view.Name, rangeViews.Name)

						if match {
							log.Info().Msgf("View exists: %s", rangeViews.Name)

							if view.AutomaticSearch != rangeViews.AutomaticSearch {
								log.Info().Msgf("View updating automatic search: %s", rangeViews.Name)
								err := client.Views().UpdateAutomaticSearch(rangeViews.Name, view.AutomaticSearch)
								if err != nil {
									log.Error().Err(err).Msgf("Error setting automatic search: %s", rangeViews.Name)
									continue
								}
							}

							if view.DefaultQuery != rangeViews.DefaultQuery.Name {
								log.Info().Msgf("View updating default query: %s", rangeViews.Name)
								err := client.Views().UpdateDefaultSavedQuery(rangeViews.Name, view.DefaultQuery)
								if err != nil {
									log.Error().Err(err).Msgf("Error setting default saved query: %s", rangeViews.Name)
									continue
								}
							}
						}
					}
				} else {
					if view.Name != "" {
						result, err1 := client.Views().Get(view.Name)
						if err1 != nil {
							log.Info().Msgf("View missing, creating (not yet implemented): %s", view.Name)
							continue
						}

						if result.Name == "" {
							log.Info().Msgf("View missing, creating (not yet implemented): %s", view.Name)
							continue
						} else {
							log.Info().Msgf("View exists: %s", result.Name)

							if view.AutomaticSearch != result.AutomaticSearch {
								log.Info().Msgf("View updating automatic search: %s", result.Name)
								err := client.Views().UpdateAutomaticSearch(view.Name, view.AutomaticSearch)
								if err != nil {
									log.Error().Err(err).Msgf("Error setting automatic search: %s", result.Name)
									continue
								}
							}

							if view.DefaultQuery != result.DefaultQuery.Name {
								log.Info().Msgf("View updating default query: %s", result.Name)
								err := client.Views().UpdateDefaultSavedQuery(view.Name, view.DefaultQuery)
								if err != nil {
									log.Error().Err(err).Msgf("Error setting default saved query: %s", result.Name)
									continue
								}
							}
						}
					}
				}
			}
		},
	}

	cmd.Flags().StringVar(&filePath, "file", "", "The local file path of the yaml file to apply.")
	cmd.Flags().StringVar(&url, "url", "", "A URL to fetch the yaml file from to apply.")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "List changes without applying them.")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "Show changes.")

	return &cmd
}

// YAML layout
type Config struct {
	Users          []Users              `yaml:"users"`
	Repos          []YamlRepos          `yaml:"repos"`
	Views          []YamlViews          `yaml:"views"`
	DefaultQueries []YamlDefaultQueries `yaml:"defaultQueries"`
}

type Users struct {
	Username      string          `yaml:"username"`
	Email         string          `yaml:"email"`
	Company       string          `yaml:"company"`
	SearchDomains []SearchDomains `yaml:"searchDomains"`
}

type SearchDomains struct {
	Name  string `yaml:"name"`
	Role  string `yaml:"role"`
	Regex bool   `yaml:"regex"`
}

type YamlRepos struct {
	Name            string `yaml:"name"`
	AutomaticSearch bool   `yaml:"automaticSearch"`
	DefaultQuery    string `yaml:"defaultQuery"`
	Regex           bool   `yaml:"regex"`
}

type YamlViews struct {
	Name            string `yaml:"name"`
	AutomaticSearch bool   `yaml:"automaticSearch"`
	DefaultQuery    string `yaml:"defaultQuery"`
	Regex           bool   `yaml:"regex"`
}

type YamlDefaultQueries struct {
	Name        string `yaml:"name"`
	Global      bool   `yaml:"global"`
	QueryString string `yaml:"queryString"`
	Start       string `yaml:"start"`
	Options     string `yaml:"options"`
}

// func reconcileUserDomains(client *api.Client, username, role, searchDomain string) error {
// 	// Check if user already has permissions
// 	resultSearchDomain, _ := client.SearchDomains().Get(searchDomain)

// 	// if searchDomain.Name != "" && searchDomain.Role != "" {
// 	// 	// Get Role
// 	// 	role, errRole := client.Roles().Get(searchDomain.Role)
// 	// 	if errRole != nil {
// 	// 		cmd.Printf("ERROR: %s", errRole)
// 	// 	}

// 	// 	resultSearchDomain, errSD := client.SearchDomains().Get(searchDomain.Name)
// 	// 	if errSD != nil {
// 	// 		cmd.Printf("ERROR: %s", errSD)
// 	// 	}

// 	// 	err := client.SearchDomains().UpdateUserPermissions(resultSearchDomain.Id, result.ID, role.ID)
// 	// 	//exitOnError(cmd, err, "Error adding user/group to search domain")
// 	// 	if err != nil {
// 	// 		cmd.Printf("ERROR adding user/group to search domain: %s\n", err)
// 	// 	} else {
// 	// 		cmd.Printf("[Users] Added %s to %s\n", user.Username, searchDomain.Name)
// 	// 	}

// 	// } else {
// 	// 	cmd.Printf("Name, filter and role fields are requires to assign a user to a repo/view: %s, %s", user.Username, searchDomain.Name)
// 	// }
// }
