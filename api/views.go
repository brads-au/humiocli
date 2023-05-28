package api

import (
	"fmt"
	"sort"
	"strings"

	"github.com/shurcooL/graphql"
)

type Views struct {
	client *Client
}

type ViewConnection struct {
	RepoName string
	Filter   string
}

type ViewQueryData struct {
	Name        string
	Description string
	ViewInfo    struct {
		Connections []struct {
			Repository struct{ Name string }
			Filter     string
		}
	} `graphql:"... on View"`
	AutomaticSearch bool
	DefaultQuery    SavedQuery
}

type View struct {
	Name            string
	Description     string
	Connections     []ViewConnection
	AutomaticSearch bool
	DefaultQuery    SavedQuery
}

func (c *Client) Views() *Views { return &Views{client: c} }

func (c *Views) Get(name string) (*View, error) {
	var query struct {
		Result ViewQueryData `graphql:"searchDomain(name: $name)"`
	}

	variables := map[string]interface{}{
		"name": graphql.String(name),
	}

	err := c.client.Query(&query, variables)
	if err != nil {
		return nil, err
	}

	connections := make([]ViewConnection, len(query.Result.ViewInfo.Connections))
	for i, data := range query.Result.ViewInfo.Connections {
		connections[i] = ViewConnection{
			RepoName: data.Repository.Name,
			Filter:   data.Filter,
		}
	}

	view := View{
		Name:            query.Result.Name,
		Description:     query.Result.Description,
		Connections:     connections,
		AutomaticSearch: query.Result.AutomaticSearch,
		DefaultQuery:    query.Result.DefaultQuery,
	}

	return &view, nil
}

type ViewListItem struct {
	Name            string
	AutomaticSearch bool
	DefaultQuery    SavedQuery
}

func (c *Views) List() ([]ViewListItem, error) {
	var query struct {
		View []ViewListItem `graphql:"searchDomains"`
	}

	err := c.client.Query(&query, nil)

	sort.Slice(query.View, func(i, j int) bool {
		return strings.ToLower(query.View[i].Name) < strings.ToLower(query.View[j].Name)
	})

	return query.View, err
}

type ViewConnectionInput struct {
	RepositoryName graphql.String `json:"repositoryName"`
	Filter         graphql.String `json:"filter"`
}

func (c *Views) Create(name, description string, connections map[string]string) error {
	var mutation struct {
		CreateView struct {
			Name        string
			Description string
		} `graphql:"createView(name: $name, description: $description, connections: $connections)"`
	}

	var viewConnections []ViewConnectionInput
	for k, v := range connections {
		viewConnections = append(
			viewConnections,
			ViewConnectionInput{
				RepositoryName: graphql.String(k),
				Filter:         graphql.String(v),
			})
	}

	variables := map[string]interface{}{
		"name":        graphql.String(name),
		"description": graphql.String(description),
		"connections": viewConnections,
	}

	return c.client.Mutate(&mutation, variables)
}

func (c *Views) Delete(name, reason string) error {
	var mutation struct {
		DeleteSearchDomain struct {
			// We have to make a selection, so just take __typename
			Typename graphql.String `graphql:"__typename"`
		} `graphql:"deleteSearchDomain(name: $name, deleteMessage: $reason)"`
	}
	variables := map[string]interface{}{
		"name":   graphql.String(name),
		"reason": graphql.String(reason),
	}

	return c.client.Mutate(&mutation, variables)
}

func (c *Views) UpdateConnections(name string, connections map[string]string) error {
	var mutation struct {
		View struct {
			Name string
		} `graphql:"updateView(viewName: $viewName, connections: $connections)"`
	}

	var viewConnections []ViewConnectionInput
	for k, v := range connections {
		viewConnections = append(
			viewConnections,
			ViewConnectionInput{
				RepositoryName: graphql.String(k),
				Filter:         graphql.String(v),
			})
	}

	variables := map[string]interface{}{
		"viewName":    graphql.String(name),
		"connections": viewConnections,
	}

	return c.client.Mutate(&mutation, variables)
}

func (c *Views) UpdateDescription(name string, description string) error {
	var mutation struct {
		UpdateDescriptionMutation struct {
			// We have to make a selection, so just take __typename
			Typename graphql.String `graphql:"__typename"`
		} `graphql:"updateDescriptionForSearchDomain(name: $name, newDescription: $description)"`
	}

	variables := map[string]interface{}{
		"name":        graphql.String(name),
		"description": graphql.String(description),
	}

	return c.client.Mutate(&mutation, variables)
}

func (c *Views) UpdateAutomaticSearch(name string, automaticSearch bool) error {
	var mutation struct {
		UpdateAutomaticSearch struct {
			// We have to make a selection, so just take __typename
			Typename graphql.String `graphql:"__typename"`
		} `graphql:"setAutomaticSearching(name: $name, automaticSearch: $automaticSearch)"`
	}

	variables := map[string]interface{}{
		"name":            graphql.String(name),
		"automaticSearch": graphql.Boolean(automaticSearch),
	}

	return c.client.Mutate(&mutation, variables)
}

func (c *Views) UpdateDefaultSavedQuery(viewName, query string) error {
	queryInfo, err := c.client.SavedQueries().Get(query, viewName)
	if err != nil {
		return fmt.Errorf("Unable to get saved query: %w", err)
	}
	queryId := queryInfo.SavedQueries[0].Id

	var mutation struct {
		SetDefaultSavedQuery struct {
			// We have to make a selection, so just take __typename
			Typename graphql.String `graphql:"__typename"`
		} `graphql:"setDefaultSavedQuery(input: { savedQueryId: $savedQueryId, viewName: $viewName })"`
	}

	variables := map[string]interface{}{
		"savedQueryId": graphql.String(queryId),
		"viewName":     graphql.String(viewName),
	}

	return c.client.Mutate(&mutation, variables)
}
