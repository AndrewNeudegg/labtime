package api

import "github.com/xanzy/go-gitlab"

func (g *GitlabAPI) GetIssue(issueID int) (*gitlab.Issue, error) {
	issue, _, err := g.client.Issues.GetIssue(g.project, issueID)
	return issue, err
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

func (g *GitlabAPI) GetIssueNotes(issueID int) ([]*gitlab.Note, error) {
	var err error

	opt := &gitlab.ListIssueNotesOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 10,
			Page:    1,
		},
	}

	retIssueNotes := []*gitlab.Note{}

	for {
		notes, resp, err := g.client.Notes.ListIssueNotes(g.project, issueID, opt)
		if err != nil {
			break
		}
		retIssueNotes = append(retIssueNotes, notes...)
		// Exit the loop when we've seen all pages.
		if resp.CurrentPage >= resp.TotalPages {
			break
		}
		// Update the page number to get the next page.
		opt.Page = resp.NextPage
	}

	return retIssueNotes, err
}

// GetTimeIssue extracts the times from a specific issue.
func (g *GitlabAPI) GetTimeIssue(issueID int, timeEntryRegex string, timeMatchRegex string) (*Overview, error) {
	notes, err := g.GetIssueNotes(issueID)
	if err != nil {
		return nil, err
	}
	timeNotes, err := extractTimeNotes(notes, timeEntryRegex)
	if err != nil {
		return nil, err
	}

	timeEntries := []TimeSpentEntry{}
	for _, timeNote := range timeNotes {

		time, err := extractTime(timeNote.Body, timeMatchRegex)
		if err != nil {
			return nil, err
		}

		timeEntries = append(timeEntries, TimeSpentEntry{
			ID: timeNote.ID,
			Author: &gitlab.Author{
				ID:       timeNote.Author.ID,
				Username: timeNote.Author.Username,
				Email:    timeNote.Author.Email,
				Name:     timeNote.Author.Name,
			},
			CreatedAt: timeNote.CreatedAt,
			RawBody:   timeNote.Body,
			Spent:     time,
		})
	}
	return &Overview{
		ID:         issueID,
		TimeSpents: timeEntries,
	}, nil
}
