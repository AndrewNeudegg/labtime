// Package api is used to access gitlab api restfully.
package api

import (
	"fmt"

	"github.com/xanzy/go-gitlab"
)

// GitlabAPI stores the current instance client.
type GitlabAPI struct {
	client  *gitlab.Client
	project string
}

// Connect establishes a connection with a Gitlab instance.
func Connect(project string, token string, domain string) (*GitlabAPI, error) {
	api := GitlabAPI{}
	api.project = project
	api.client = gitlab.NewClient(nil, token)
	err := api.client.SetBaseURL(domain)
	if err != nil {
		return nil, err
	}
	return &api, nil
}

func (g *GitlabAPI) GetUsers() ([]*gitlab.User, error) {
	if g.client == nil {
		return nil, fmt.Errorf("client has not been instantiated")
	}
	users, _, err := g.client.Users.ListUsers(nil)
	return users, err
}
