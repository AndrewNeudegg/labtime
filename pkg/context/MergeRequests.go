package context

import (
	"fmt"
	"strings"

	"github.com/xanzy/go-gitlab"

	log "github.com/sirupsen/logrus"
)

//
// MR access functions
//

func (c *Context) ListProjectMRs() []*gitlab.MergeRequest {
	log.Info(fmt.Sprintf("Listing project MRs for %s.", c.AppConfig.Instance.Project))
	MRs, err := c.API.GetMRs()
	log.Info(fmt.Sprintf("Found %v MRs.", len(MRs)))
	handleErr(err)
	return MRs
}

func (c *Context) GetMRByID(MRID int) *gitlab.MergeRequest {
	log.Info(fmt.Sprintf("Fetching project MR matching ID %v.", MRID))
	MR, err := c.API.GetMR(MRID)
	handleErr(err)
	log.Info(fmt.Sprintf("Found MR matching %v.", MR.IID))
	return MR
}

//
// Utility Functions
//

func (c *Context) ConcatMRLabels(MR *gitlab.MergeRequest, seperator string) string {
	log.Info(fmt.Sprintf(`Concatonating %v labels with the seperator "%s".`, len(MR.Labels), seperator))
	return strings.Join(MR.Labels, seperator)
}

func (c *Context) TotalMRTimeSpent(MR *gitlab.MergeRequest) float64 {
	log.Info(fmt.Sprintf(`Calculating total time spent on MR %v.`, MR.IID))
	overview, err := c.API.GetTimeMR(MR.IID, c.AppConfig.QueryConfig.TimeEntryDetectionRegex, c.AppConfig.QueryConfig.TimeEntryExtractionRegex)
	handleErr(err)
	log.Info(fmt.Sprintf(`Found %v time entries`, len(overview.TimeSpents)))

	totalTime := float64(0)
	for _, timeSpent := range overview.TimeSpents {
		totalTime += timeSpent.Spent.TotalDaysRounded()
	}
	return totalTime
}
