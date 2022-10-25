package examples

import "fmt"

type Endpoints struct {
	CurrentUserUrl    string `json:"current_user_url"`
	AuthorizationsUrl string `json:"authorizations_url"`
	RepositoryUrl     string `json:"repository_url"`
}

func GetEndpoints() (*Endpoints, error) {
	resp, err := httpClient.Get("https://api.github.com")
	if err != nil {
		return nil, err
	}

	fmt.Printf("status: %s\n", resp.Status)
	fmt.Printf("status code: %v\n", resp.StatusCode)
	fmt.Printf("resp: %s\n", resp.String())

	var endpoints Endpoints
	if err := resp.UnmarshalJson(&endpoints); err != nil {
		return nil, err
	}

	fmt.Printf("repositories URL: %s\n", endpoints.RepositoryUrl)
	return &endpoints, err
}
