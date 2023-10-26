package api

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
		"name": string(name),
	}

	err := s.client.Query(&query, variables)
	return query.SearchDomains, err
}
