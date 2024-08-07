package humiographql

import (
	graphql "github.com/cli/shurcooL-graphql"
)

type FilterAlert struct {
	ID                  graphql.String   `graphql:"id"`
	Name                graphql.String   `graphql:"name"`
	Description         graphql.String   `graphql:"description"`
	QueryString         graphql.String   `graphql:"queryString"`
	Actions             []Action         `graphql:"actions"`
	Labels              []graphql.String `graphql:"labels"`
	Enabled             graphql.Boolean  `graphql:"enabled"`
	ThrottleTimeSeconds Long             `graphql:"throttleTimeSeconds"`
	ThrottleField       graphql.String   `graphql:"throttleField"`
	QueryOwnership      QueryOwnership   `graphql:"queryOwnership"`
}

type CreateFilterAlert struct {
	ViewName            RepoOrViewName     `json:"viewName"`
	Name                graphql.String     `json:"name"`
	Description         graphql.String     `json:"description,omitempty"`
	QueryString         graphql.String     `json:"queryString"`
	ActionIdsOrNames    []graphql.String   `json:"actionIdsOrNames"`
	Labels              []graphql.String   `json:"labels"`
	Enabled             graphql.Boolean    `json:"enabled"`
	ThrottleTimeSeconds Long               `json:"throttleTimeSeconds,omitempty"`
	ThrottleField       graphql.String     `json:"throttleField,omitempty"`
	RunAsUserID         graphql.String     `json:"runAsUserId,omitempty"`
	QueryOwnershipType  QueryOwnershipType `json:"queryOwnershipType"`
}

type UpdateFilterAlert struct {
	ViewName            RepoOrViewName     `json:"viewName"`
	ID                  graphql.String     `json:"id"`
	Name                graphql.String     `json:"name"`
	Description         graphql.String     `json:"description,omitempty"`
	QueryString         graphql.String     `json:"queryString"`
	ActionIdsOrNames    []graphql.String   `json:"actionIdsOrNames"`
	Labels              []graphql.String   `json:"labels"`
	Enabled             graphql.Boolean    `json:"enabled"`
	ThrottleTimeSeconds Long               `json:"throttleTimeSeconds,omitempty"`
	ThrottleField       graphql.String     `json:"throttleField,omitempty"`
	RunAsUserID         graphql.String     `json:"runAsUserId,omitempty"`
	QueryOwnershipType  QueryOwnershipType `json:"queryOwnershipType"`
}
