package context

// import (
// 	"strings"

// 	"github.com/xanzy/go-gitlab"
// )

// //
// // Issue access functions
// //

// func (c *Context) ListProjectMRs() []*gitlab.MergeRequest {
// 	issues, err := c.API.GetIssues()
// 	handleErr(err)
// 	return issues
// }

// func (c *Context) GetMRByID(issueID int) *gitlab.MergeRequest {
// 	issue, err := c.API.GetIssue(issueID)
// 	handleErr(err)
// 	return issue
// }

// //
// // Utility Functions
// //

// func (c *Context) ConcatMRLabels(issue *gitlab.MergeRequest, seperator string) string {
// 	return strings.Join(issue.Labels, seperator)
// }

// func (c *Context) TotalTimeSpent(issue *gitlab.MergeRequest) float64 {
// 	overview, err := c.API.GetTimeIssue(issue.IID, c.AppConfig.QueryConfig.TimeEntryDetectionRegex, c.AppConfig.QueryConfig.TimeEntryExtractionRegex)
// 	handleErr(err)
// 	totalTime := float64(0)
// 	for _, timeSpent := range overview.TimeSpents {
// 		totalTime += timeSpent.Spent.TotalDaysRounded()
// 	}
// 	return totalTime
// }
