package api

import (
	"fmt"
)

type SavedQueries struct {
	client *Client
}

type SavedQuery struct {
	Id    string
	Name  string
	Query struct {
		QueryString string
		Start       string
		End         string
		IsLive      bool
	}
}

type QueryOptions struct {
	Columns string
}

func (c *Client) SavedQueries() *SavedQueries { return &SavedQueries{client: c} }

func (s *SavedQueries) List(name string) (*SearchDomain, error) {
	query, err := s.client.SearchDomains().Get(name)
	if err != nil {
		return nil, err
	}

	savedQuery := SearchDomain{
		Id:           query.Id,
		Name:         query.Name,
		SavedQueries: query.SavedQueries,
	}

	return &savedQuery, nil
}

func (s *SavedQueries) Get(query, viewName string) (*SearchDomain, error) {
	savedQueries, err := s.List(viewName)
	if err != nil {
		return nil, fmt.Errorf("unable to get saved queries: %w", err)
	}

	var matched []SavedQuery
	for _, data := range savedQueries.SavedQueries {
		if (query == data.Id) || (query == data.Name) {
			matched = append(matched, data)
		}
	}

	if matched == nil {
		return nil, fmt.Errorf("no saved query found by that name/id")
	}

	savedQuery := SearchDomain{
		Id:           savedQueries.Id,
		Name:         savedQueries.Name,
		SavedQueries: matched,
	}

	return &savedQuery, nil
}

func (s *SavedQueries) Create(name, viewName, queryString, start, end string, isLive bool, widgetType string) error {
	var mutation struct {
		CreateSavedQuery struct {
			// We have to make a selection, so just take __typename
			Typename string `graphql:"__typename"`
		} `graphql:"createSavedQuery(input:{ name: $name, viewName: $viewName, queryString: $queryString, start: $start, end: $end, isLive: $isLive, widgetType: $widgetType})"`
	}

	variables := map[string]interface{}{
		"name":        string(name),
		"viewName":    string(viewName),
		"queryString": string(queryString),
		"start":       string(start),
		"end":         string(end),
		"isLive":      bool(isLive),
		"widgetType":  string(widgetType),
	}

	err := s.client.Mutate(&mutation, variables)
	if err != nil {
		// The graphql error message is vague if the saved query already exists, so add a hint.
		return fmt.Errorf("%w. Does the saved query already exist?", err)
	}

	return nil
}

func (s *SavedQueries) Delete(query, viewName string) error {
	savedQueries, err := s.List(viewName)
	if err != nil {
		return fmt.Errorf("unable to get saved queries: %w", err)
	}

	var matched []SavedQuery
	for _, data := range savedQueries.SavedQueries {
		if (query == data.Id) || (query == data.Name) {
			matched = append(matched, data)
		}
	}

	var mutation struct {
		DeleteSavedQuery struct {
			// We have to make a selection, so just take __typename
			Typename string `graphql:"__typename"`
		} `graphql:"deleteSavedQuery(input: { id: $id, viewName: $viewName })"`
	}

	variables := map[string]interface{}{
		"id":       string(matched[0].Id),
		"viewName": string(viewName),
	}

	return s.client.Mutate(&mutation, variables)
}
