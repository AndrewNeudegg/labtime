package context

import (
	"fmt"
	"strings"

	"github.com/xanzy/go-gitlab"

	log "github.com/sirupsen/logrus"
)

//
// Issue access functions
//

func (c *Context) ListProjectIssues() []*gitlab.Issue {
	log.Info(fmt.Sprintf("Listing project issues for %s.", c.AppConfig.Instance.Project))
	issues, err := c.API.GetIssues()
	log.Info(fmt.Sprintf("Found %v issues.", len(issues)))
	handleErr(err)
	return issues
}

func (c *Context) GetIssueByID(issueID int) *gitlab.Issue {
	log.Info(fmt.Sprintf("Fetching project issue matching ID %v.", issueID))
	issue, err := c.API.GetIssue(issueID)
	handleErr(err)
	log.Info(fmt.Sprintf("Found issue matching %v.", issue.IID))
	return issue
}

//
// Utility Functions
//

func (c *Context) ConcatIssueLabels(issue *gitlab.Issue, seperator string) string {
	log.Info(fmt.Sprintf(`Concatonating %v labels with the seperator "%s".`, len(issue.Labels), seperator))
	return strings.Join(issue.Labels, seperator)
}

func (c *Context) TotalIssueTimeSpent(issue *gitlab.Issue) float64 {
	log.Info(fmt.Sprintf(`Calculating total time spent on issue %v.`, issue.IID))
	overview, err := c.API.GetTimeIssue(issue.IID, c.AppConfig.QueryConfig.TimeEntryDetectionRegex, c.AppConfig.QueryConfig.TimeEntryExtractionRegex)
	handleErr(err)
	log.Info(fmt.Sprintf(`Found %v time entries`, len(overview.TimeSpents)))

	totalTime := float64(0)
	for _, timeSpent := range overview.TimeSpents {
		totalTime += timeSpent.Spent.TotalDaysRounded()
	}
	return totalTime
}
