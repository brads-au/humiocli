package api

import (
	"github.com/shurcooL/graphql"
)

type SearchDomains struct {
	client *Client
}

type SearchDomain struct {
	Id           string
	Name         string
	Description  string
	SavedQueries []SavedQuery
	UsersV2      struct {
		Users []struct {
			Id       string
			Username string
		}
	}
}

func (c *Client) SearchDomains() *SearchDomains { return &SearchDomains{client: c} }

func (s *SearchDomains) Get(name string) (SearchDomain, error) {
	var query struct {
		SearchDomains SearchDomain `graphql:"searchDomain(name: $name)"`
	}

	variables := map[string]interface{}{
		"name": graphql.String(name),
	}

	err := s.client.Query(&query, variables)
	return query.SearchDomains, err
}

func (s *SearchDomains) List() ([]SearchDomain, error) {
	var query struct {
		SearchDomains []SearchDomain `graphql:"searchDomains"`
	}

	err := s.client.Query(&query, nil)
	return query.SearchDomains, err
}

func (s *SearchDomains) UpdateUserPermissions(searchDomain string, identity string, role string) error {
	type GroupRoleAssignment struct {
		GroupId graphql.String `json:"groupId"`
		RoleId  graphql.String `json:"roleId"`
	}

	type UserRoleAssignment struct {
		UserId graphql.String `json:"userId"`
		RoleId graphql.String `json:"roleId"`
	}

	// Groups... TBD
	//var groups []GroupRoleAssignment
	groups := []GroupRoleAssignment{}

	var users []UserRoleAssignment
	users = append(
		users,
		UserRoleAssignment{
			UserId: graphql.String(identity),
			RoleId: graphql.String(role),
		})

	var mutation struct {
		ChangePermissions []struct {
			// We have to make a selection, so just take __typename
			Typename graphql.String `graphql:"__typename"`
		} `graphql:"changeUserAndGroupRolesForSearchDomain(searchDomainId: $searchDomainId, groups: $groups, users: $users)"`
	}

	variables := map[string]interface{}{
		"searchDomainId": graphql.String(searchDomain),
		"groups":         groups,
		"users":          users,
	}

	//fmt.Printf("\n---DEBUG: %+v\n\n", variables)

	return s.client.Mutate(&mutation, variables)
}
