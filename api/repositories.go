package api

import (
	"fmt"
	"strings"
)

type Repositories struct {
	client *Client
}

type Repository struct {
	ID                     string
	Name                   string
	Description            string
	RetentionDays          float64 `graphql:"timeBasedRetention"`
	IngestRetentionSizeGB  float64 `graphql:"ingestSizeBasedRetention"`
	StorageRetentionSizeGB float64 `graphql:"storageSizeBasedRetention"`
	SpaceUsed              int64   `graphql:"compressedByteSize"`
	AutomaticSearch        bool
	DefaultQuery           SavedQuery
}

func (c *Client) Repositories() *Repositories { return &Repositories{client: c} }

func (r *Repositories) Get(name string) (Repository, error) {
	var query struct {
		Repository Repository `graphql:"repository(name: $name)"`
	}

	variables := map[string]interface{}{
		"name": name,
	}

	err := r.client.Query(&query, variables)

	if err != nil {
		// The graphql error message is vague if the repo already exists, so add a hint.
		return query.Repository, fmt.Errorf("%w. Does the repo already exist?", err)
	}

	return query.Repository, nil
}

type RepoListItem struct {
	ID        string
	Name      string
	SpaceUsed int64 `graphql:"compressedByteSize"`
}

func (r *Repositories) List() ([]RepoListItem, error) {
	var query struct {
		Repositories []RepoListItem `graphql:"repositories"`
	}

	err := r.client.Query(&query, nil)
	return query.Repositories, err
}

func (r *Repositories) Create(name string) error {
	var mutation struct {
		CreateRepository struct {
			Repository Repository
		} `graphql:"createRepository(name: $name)"`
	}

	variables := map[string]interface{}{
		"name": name,
	}

	err := r.client.Mutate(&mutation, variables)
	if err != nil {
		// The graphql error message is vague if the repo already exists, so add a hint.
		return fmt.Errorf("%w. Does the repo already exist?", err)
	}

	return nil
}

func (r *Repositories) Delete(name, reason string, allowDataDeletion bool) error {
	existingRepo, err := r.Get(name)
	if err != nil {
		return err
	}
	safeToDelete := allowDataDeletion || existingRepo.SpaceUsed == 0
	if !safeToDelete {
		return fmt.Errorf("repository contains data and data deletion not allowed")
	}

	var mutation struct {
		DeleteSearchDomain struct {
			// We have to make a selection, so just take __typename
			Typename string `graphql:"__typename"`
		} `graphql:"deleteSearchDomain(name: $name, deleteMessage: $reason)"`
	}
	variables := map[string]interface{}{
		"name":   name,
		"reason": reason,
	}

	return r.client.Mutate(&mutation, variables)
}

type DefaultGroupEnum string

const (
	DefaultGroupEnumMember     DefaultGroupEnum = "Member"
	DefaultGroupEnumAdmin      DefaultGroupEnum = "Admin"
	DefaultGroupEnumEliminator DefaultGroupEnum = "Eliminator"
)

func (e DefaultGroupEnum) String() string {
	return string(e)
}

func (e *DefaultGroupEnum) ParseString(s string) bool {
	switch strings.ToLower(s) {
	case "member":
		*e = DefaultGroupEnumMember
		return true
	case "admin":
		*e = DefaultGroupEnumAdmin
		return true
	case "eliminator":
		*e = DefaultGroupEnumEliminator
		return true
	default:
		return false
	}
}

func (r *Repositories) UpdateUserGroup(name, username string, groups ...DefaultGroupEnum) error {
	if len(groups) == 0 {
		return fmt.Errorf("at least one group must be defined")
	}

	var mutation struct {
		UpdateDefaultGroupMembershipsMutation struct {
			// We have to make a selection, so just take __typename
			Typename string `graphql:"__typename"`
		} `graphql:"updateDefaultGroupMemberships(input: {viewName: $name, userName: $username, groups: $groups})"`
	}
	variables := map[string]interface{}{
		"name":     name,
		"username": username,
		"groups":   groups,
	}

	return r.client.Mutate(&mutation, variables)
}

func (r *Repositories) UpdateTimeBasedRetention(name string, retentionInDays float64, allowDataDeletion bool) error {
	existingRepo, err := r.Get(name)
	if err != nil {
		return err
	}
	safeToDelete := allowDataDeletion || existingRepo.SpaceUsed == 0

	var mutation struct {
		UpdateRetention struct {
			// We have to make a selection, so just take __typename
			Typename string `graphql:"__typename"`
		} `graphql:"updateRetention(repositoryName: $name, timeBasedRetention: $retentionInDays)"`
	}
	variables := map[string]interface{}{
		"name":            name,
		"retentionInDays": (*float64)(nil),
	}
	if retentionInDays > 0 {
		if retentionInDays < existingRepo.RetentionDays || existingRepo.RetentionDays == 0 {
			if !safeToDelete {
				return fmt.Errorf("repository contains data and data deletion not allowed")
			}
		}
		variables["retentionInDays"] = retentionInDays
	}

	return r.client.Mutate(&mutation, variables)
}

func (r *Repositories) UpdateStorageBasedRetention(name string, storageInGB float64, allowDataDeletion bool) error {
	existingRepo, err := r.Get(name)
	if err != nil {
		return err
	}
	safeToDelete := allowDataDeletion || existingRepo.SpaceUsed == 0

	var mutation struct {
		UpdateRetention struct {
			// We have to make a selection, so just take __typename
			Typename string `graphql:"__typename"`
		} `graphql:"updateRetention(repositoryName: $name, storageSizeBasedRetention: $storageInGB)"`
	}
	variables := map[string]interface{}{
		"name":        name,
		"storageInGB": (*float64)(nil),
	}
	if storageInGB > 0 {
		if storageInGB < existingRepo.StorageRetentionSizeGB || existingRepo.StorageRetentionSizeGB == 0 {
			if !safeToDelete {
				return fmt.Errorf("repository contains data and data deletion not allowed")
			}
		}
		variables["storageInGB"] = storageInGB
	}

	return r.client.Mutate(&mutation, variables)
}

func (r *Repositories) UpdateIngestBasedRetention(name string, ingestInGB float64, allowDataDeletion bool) error {
	existingRepo, err := r.Get(name)
	if err != nil {
		return err
	}
	safeToDelete := allowDataDeletion || existingRepo.SpaceUsed == 0

	var mutation struct {
		UpdateRetention struct {
			// We have to make a selection, so just take __typename
			Typename string `graphql:"__typename"`
		} `graphql:"updateRetention(repositoryName: $name, ingestSizeBasedRetention: $ingestInGB)"`
	}
	variables := map[string]interface{}{
		"name":       name,
		"ingestInGB": (*float64)(nil),
	}
	if ingestInGB > 0 {
		if ingestInGB < existingRepo.IngestRetentionSizeGB || existingRepo.IngestRetentionSizeGB == 0 {
			if !safeToDelete {
				return fmt.Errorf("repository contains data and data deletion not allowed")
			}
		}
		variables["ingestInGB"] = ingestInGB
	}

	return r.client.Mutate(&mutation, variables)
}

func (r *Repositories) UpdateDescription(name, description string) error {
	var mutation struct {
		UpdateDescription struct {
			// We have to make a selection, so just take __typename
			Typename string `graphql:"__typename"`
		} `graphql:"updateDescriptionForSearchDomain(name: $name, newDescription: $description)"`
	}

	variables := map[string]interface{}{
		"name":        name,
		"description": description,
	}

	return r.client.Mutate(&mutation, variables)
}

func (r *Repositories) UpdateAutomaticSearch(name string, automaticSearch bool) error {
	var mutation struct {
		UpdateAutomaticSearch struct {
			// We have to make a selection, so just take __typename
			Typename string `graphql:"__typename"`
		} `graphql:"setAutomaticSearching(name: $name, automaticSearch: $automaticSearch)"`
	}

	variables := map[string]interface{}{
		"name":            string(name),
		"automaticSearch": bool(automaticSearch),
	}

	return r.client.Mutate(&mutation, variables)
}

func (r *Repositories) UpdateDefaultSavedQuery(viewName, query string) error {
	queryInfo, err := r.client.SavedQueries().Get(query, viewName)
	if err != nil {
		return fmt.Errorf("unable to get saved query: %w", err)
	}
	queryId := queryInfo.SavedQueries[0].Id

	var mutation struct {
		SetDefaultSavedQuery struct {
			// We have to make a selection, so just take __typename
			Typename string `graphql:"__typename"`
		} `graphql:"setDefaultSavedQuery(input: { savedQueryId: $savedQueryId, viewName: $viewName })"`
	}

	variables := map[string]interface{}{
		"savedQueryId": string(queryId),
		"viewName":     string(viewName),
	}

	return r.client.Mutate(&mutation, variables)
}
