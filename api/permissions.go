package api

import (
	graphql "github.com/cli/shurcooL-graphql"
)

type Permissions struct {
	client *Client
}

type PermissionValues struct {
	Name         string
	Description  string
	IsDeprecated bool
}

type PermissionData struct {
	Name       string
	EnumValues []PermissionValues `json:"enumValues"`
}

func (c *Client) Permissions() *Permissions { return &Permissions{client: c} }

func (p *Permissions) List(permissionType string) ([]PermissionValues, error) {
	var query struct {
		Permissions PermissionData `graphql:"__type(name: $name)"`
	}

	variables := map[string]interface{}{
		"name": graphql.String(permissionType),
	}

	err := p.client.Query(&query, variables)
	if err != nil {
		return nil, err
	}

	return query.Permissions.EnumValues, nil
}

func (p *Permissions) Check(permissionType string, permissions []string) (bool, error) {
	permissionSet := make(map[string]bool)

	validPermissions, err := p.client.Permissions().List(permissionType)
	if err != nil {
		return false, err
	}

	for _, v := range validPermissions {
		permissionSet[v.Name] = true
	}

	for _, permission := range permissions {
		if !permissionSet[permission] {
			return false, nil
		}
	}

	return true, nil
}
