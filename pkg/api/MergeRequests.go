package api

import "github.com/xanzy/go-gitlab"

func (g *GitlabAPI) GetMR(MRID int) (*gitlab.MergeRequest, error) {
	MR, _, err := g.client.MergeRequests.GetMergeRequest(g.project, MRID, nil)
	return MR, err
}

func (g *GitlabAPI) GetMRs() ([]*gitlab.MergeRequest, error) {
	var err error

	opt := &gitlab.ListProjectMergeRequestsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 10,
			Page:    1,
		},
	}

	retMRs := []*gitlab.MergeRequest{}

	for {
		MRs, resp, err := g.client.MergeRequests.ListProjectMergeRequests(g.project, opt)
		if err != nil {
			break
		}
		retMRs = append(retMRs, MRs...)
		// Exit the loop when we've seen all pages.
		if resp.CurrentPage >= resp.TotalPages {
			break
		}
		// Update the page number to get the next page.
		opt.Page = resp.NextPage
	}

	return retMRs, err
}

func (g *GitlabAPI) GetMRNotes(MRID int) ([]*gitlab.Note, error) {
	var err error

	opt := &gitlab.ListMergeRequestNotesOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 10,
			Page:    1,
		},
	}

	retMRNotes := []*gitlab.Note{}

	for {
		notes, resp, err := g.client.Notes.ListMergeRequestNotes(g.project, MRID, opt)
		if err != nil {
			break
		}
		retMRNotes = append(retMRNotes, notes...)
		// Exit the loop when we've seen all pages.
		if resp.CurrentPage >= resp.TotalPages {
			break
		}
		// Update the page number to get the next page.
		opt.Page = resp.NextPage
	}

	return retMRNotes, err
}

// GetTimeMR extracts the times from a specific MR.
func (g *GitlabAPI) GetTimeMR(MRID int, timeEntryRegex string, timeMatchRegex string) (*Overview, error) {
	notes, err := g.GetMRNotes(MRID)
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
		ID:         MRID,
		TimeSpents: timeEntries,
	}, nil
}
