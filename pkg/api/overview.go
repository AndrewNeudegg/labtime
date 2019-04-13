package api

import (
	"regexp"

	"github.com/xanzy/go-gitlab"
)

type Overview struct {
}

func (g *GitlabAPI) GetTimeIssue(issueID int, timeEntryRegex string, timeMatchRegex string) (*Overview, error) {
	notes, err := g.GetIssueNotes(issueID)
	if err != nil {
		return nil, err
	}
	timeNotes, err := extractTimeNotes(notes, timeEntryRegex)
	if err != nil {
		return nil, err
	}

}

func (g *GitlabAPI) GetTimeMR(mergeRequestID int) (*Overview, error) {
	return nil, nil
}

func extractTimeNotes(notes []*gitlab.Note, timeEntryRegex string) ([]*gitlab.Note, error) {
	timeNotes := []*gitlab.Note{}
	r, err := regexp.Compile(timeEntryRegex)
	if err != nil {
		return nil, err
	}
	for _, note := range notes {
		if r.MatchString(note.Body) {
			timeNotes = append(timeNotes, note)
		}
	}
	return timeNotes, nil
}
