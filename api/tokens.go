package api

import (
	"fmt"
	"strings"

	graphql "github.com/cli/shurcooL-graphql"
)

type Tokens struct {
	client *Client
}

type Token struct {
	Id          string
	Name        string
	CreatedAt   int64
	ExpireAt    int64
	IpFilter    string
	Permissions []string
	Views       []string
	Type        string
}

type TokenData struct {
	ViewPermissionToken         ViewPermissionToken         `graphql:"...on ViewPermissionsToken"`
	SystemPermissionToken       SystemPermissionToken       `graphql:"...on SystemPermissionsToken"`
	OrganizationPermissionToken OrganizationPermissionToken `graphql:"...on OrganizationPermissionsToken"`
}

type ViewPermissionToken struct {
	Id          string
	Name        string
	CreatedAt   int64
	ExpireAt    int64
	IpFilter    string
	Permissions []string
	Views       []struct {
		Name string
	}
	Type string `graphql:"__typename"`
}

type SystemPermissionToken struct {
	Id          string
	Name        string
	CreatedAt   int64
	ExpireAt    int64
	IpFilter    string
	Permissions []string
	Type        string `graphql:"__typename"`
}

type OrganizationPermissionToken struct {
	Id          string
	Name        string
	CreatedAt   int64
	ExpireAt    int64
	IpFilter    string
	Permissions []string
	Type        string `graphql:"__typename"`
}

type Permission string
type OrganizationPermission string
type SystemPermission string

func (c *Client) Tokens() *Tokens { return &Tokens{client: c} }

func (t *Tokens) List() ([]Token, error) {
	var query struct {
		Tokens struct {
			Results []TokenData
		} `graphql:"tokens(typeFilter: [ViewPermissionToken, SystemPermissionToken, OrganizationPermissionToken], sortBy: Name)"`
	}

	variables := map[string]interface{}{}

	err := t.client.Query(&query, variables)
	if err != nil {
		return nil, err
	}

	var tokens []Token
	for _, tokenData := range query.Tokens.Results {
		switch {
		case tokenData.ViewPermissionToken.Type == "ViewPermissionsToken":
			var viewSlice []string
			for _, viewData := range tokenData.ViewPermissionToken.Views {
				viewSlice = append(viewSlice, viewData.Name)
			}
			tokens = append(tokens, Token{
				Id:          tokenData.ViewPermissionToken.Id,
				Name:        tokenData.ViewPermissionToken.Name,
				CreatedAt:   tokenData.ViewPermissionToken.CreatedAt,
				ExpireAt:    tokenData.ViewPermissionToken.ExpireAt,
				IpFilter:    tokenData.ViewPermissionToken.IpFilter,
				Permissions: tokenData.ViewPermissionToken.Permissions,
				Views:       viewSlice,
				Type:        tokenData.ViewPermissionToken.Type,
			})
		case tokenData.SystemPermissionToken.Type == "SystemPermissionsToken":
			tokens = append(tokens, Token{
				Id:          tokenData.SystemPermissionToken.Id,
				Name:        tokenData.SystemPermissionToken.Name,
				CreatedAt:   tokenData.SystemPermissionToken.CreatedAt,
				ExpireAt:    tokenData.SystemPermissionToken.ExpireAt,
				IpFilter:    tokenData.SystemPermissionToken.IpFilter,
				Permissions: tokenData.SystemPermissionToken.Permissions,
				Type:        tokenData.SystemPermissionToken.Type,
			})
		case tokenData.OrganizationPermissionToken.Type == "OrganizationPermissionsToken":
			tokens = append(tokens, Token{
				Id:          tokenData.OrganizationPermissionToken.Id,
				Name:        tokenData.OrganizationPermissionToken.Name,
				CreatedAt:   tokenData.OrganizationPermissionToken.CreatedAt,
				ExpireAt:    tokenData.OrganizationPermissionToken.ExpireAt,
				IpFilter:    tokenData.OrganizationPermissionToken.IpFilter,
				Permissions: tokenData.OrganizationPermissionToken.Permissions,
				Type:        tokenData.OrganizationPermissionToken.Type,
			})
		}
	}

	return tokens, nil
}

func (t *Tokens) Add(name string, tokenType string, permissions []string, viewName string) (string, error) {
	switch strings.ToLower(tokenType) {
	case "system":
		result, errCheck := t.client.Permissions().Check("SystemPermission", permissions)
		if errCheck != nil {
			return "", errCheck
		}
		if !result {
			return "", fmt.Errorf("invalid permission provided")
		}

		var validPerms []SystemPermission
		for _, permission := range permissions {
			validPerms = append(validPerms, SystemPermission(permission))
		}

		variables := map[string]interface{}{
			"name":        graphql.String(name),
			"permissions": validPerms,
		}

		var mutation struct {
			CreateSystemPermissionsToken string `graphql:"createSystemPermissionsToken(input: {name: $name, permissions: $permissions})"`
		}

		err := t.client.Mutate(&mutation, variables)
		if err != nil {
			return "", err
		}

		return mutation.CreateSystemPermissionsToken, nil

	case "organization":
		result, errCheck := t.client.Permissions().Check("OrganizationPermission", permissions)
		if errCheck != nil {
			return "", errCheck
		}
		if !result {
			return "", fmt.Errorf("invalid permission provided")
		}

		var validPerms []OrganizationPermission
		for _, permission := range permissions {
			validPerms = append(validPerms, OrganizationPermission(permission))
		}

		variables := map[string]interface{}{
			"name":        graphql.String(name),
			"permissions": validPerms,
		}

		var mutation struct {
			CreateOrganizationPermissionsToken string `graphql:"createOrganizationPermissionsToken(input: {name: $name, permissions: $permissions})"`
		}

		err := t.client.Mutate(&mutation, variables)
		if err != nil {
			return "", err
		}

		return mutation.CreateOrganizationPermissionsToken, nil

	case "view":
		result, errCheck := t.client.Permissions().Check("Permission", permissions)
		if errCheck != nil {
			return "", errCheck
		}
		if !result {
			return "", fmt.Errorf("invalid permission provided")
		}

		var validPerms []Permission
		for _, permission := range permissions {
			validPerms = append(validPerms, Permission(permission))
		}

		viewId, errView := t.client.Views().GetViewID(viewName)
		if errView != nil {
			return "", fmt.Errorf("unable to get view ID")
		}

		variables := map[string]interface{}{
			"name":        graphql.String(name),
			"viewId":      graphql.String(viewId),
			"permissions": validPerms,
		}

		var mutation struct {
			CreateViewPermissionsToken string `graphql:"createViewPermissionsToken(input: {name: $name, viewIds: [$viewId], permissions: $permissions})"`
		}

		err := t.client.Mutate(&mutation, variables)
		if err != nil {
			return "", err
		}

		return mutation.CreateViewPermissionsToken, nil

	default:
		err := fmt.Errorf("invalid token type: %s. Use: system, organization or view", tokenType)
		return "", err
	}
}

func (t *Tokens) Delete(name string) error {
	var mutation struct {
		DeleteToken bool `graphql:"deleteToken(input: {id: $tokenId})"`
	}

	variables := map[string]interface{}{
		"tokenId": graphql.String(name),
	}

	return t.client.Mutate(&mutation, variables)
}
