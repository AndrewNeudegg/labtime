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

func (g *GitlabAPI) GetIssues() ([]*gitlab.Issue, error) {
	var err error

	opt := &gitlab.ListProjectIssuesOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 10,
			Page:    1,
		},
	}

	retIssues := []*gitlab.Issue{}

	for {
		issues, resp, err := g.client.Issues.ListProjectIssues(g.project, opt)
		if err != nil {
			break
		}
		retIssues = append(retIssues, issues...)
		// Exit the loop when we've seen all pages.
		if resp.CurrentPage >= resp.TotalPages {
			break
		}
		// Update the page number to get the next page.
		opt.Page = resp.NextPage
	}

	return retIssues, err
}
