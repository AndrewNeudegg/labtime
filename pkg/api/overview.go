package api

import (
	"math"
	"regexp"
	"strconv"
	"time"

	"github.com/xanzy/go-gitlab"
)

type Overview struct {
	ID         int              `json:"id"`
	TimeSpents []TimeSpentEntry `json:"time-spents"`
}

type TimeSpentEntry struct {
	ID        int            `json:"id"`
	Author    *gitlab.Author `json:"author"`
	CreatedAt *time.Time     `json:"created-at"`
	RawBody   string         `json:"raw-body"`
	Spent     *GitlabTime    `json:"time-spent"`
}

type GitlabTime struct {
	Months  float64 `json:"months"`
	Weeks   float64 `json:"weeks"`
	Days    float64 `json:"days"`
	Hours   float64 `json:"hours"`
	Minutes float64 `json:"minutes"`
	Seconds float64 `json:"seconds"`
}

func (glt *GitlabTime) TotalDaysRaw() float64 {
	months := float64(glt.Months * gitlabWeeksInMonth * gitlabDaysInWeek)
	weeks := float64(glt.Weeks * gitlabDaysInWeek)
	days := float64(glt.Days)
	hours := float64(glt.Hours / gitlabHoursInDay)
	minutes := float64(glt.Minutes / (gitlabMinutesInHour * gitlabHoursInDay))
	seconds := float64(glt.Seconds / (gitlabMinutesInHour * gitlabHoursInDay * gitlabSecondsInMinute))
	return months + weeks + days + hours + minutes + seconds
}

func (glt *GitlabTime) TotalDaysRounded() float64 {
	months := float64(glt.Months * gitlabWeeksInMonth * gitlabDaysInWeek)
	weeks := float64(glt.Weeks * gitlabDaysInWeek)
	days := float64(glt.Days)
	hours := float64(glt.Hours / gitlabHoursInDay)
	minutes := float64(glt.Minutes / (gitlabMinutesInHour * gitlabHoursInDay))
	seconds := float64(glt.Seconds / (gitlabMinutesInHour * gitlabHoursInDay * gitlabSecondsInMinute))
	return round(months+weeks+days+hours+minutes+seconds, threeDecimalPlaces)
}

func round(input float64, unit float64) float64 {
	return math.Ceil(input*unit) / unit
}

const (
	threeDecimalPlaces = 1000

	gitlabWeeksInMonth    = 4
	gitlabDaysInWeek      = 5
	gitlabHoursInDay      = 8
	gitlabMinutesInHour   = 60
	gitlabSecondsInMinute = 60

	gitlab_Second = 1
	gitlab_Minute = gitlab_Second * gitlabSecondsInMinute
	gitlab_Hour   = gitlab_Minute * gitlabMinutesInHour
	gitlab_Day    = gitlab_Hour * gitlabHoursInDay
	gitlab_Week   = gitlab_Day * gitlabDaysInWeek
	gitlab_Month  = gitlab_Week * gitlabWeeksInMonth
)

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

func extractTime(bodyContent string, timeGroupMatcher string) (*GitlabTime, error) {
	returnTime := GitlabTime{}
	exp, err := regexp.Compile(timeGroupMatcher)
	if err != nil {
		return nil, err
	}
	match := exp.FindStringSubmatch(bodyContent)
	for i, name := range exp.SubexpNames() {
		if i != 0 && name != "" {
			var value int

			if len(match) < i {
				break
			}

			if match[i] == "" {
				value = 0
			} else {
				value, err = strconv.Atoi(match[i])
				if err != nil {
					return nil, err
				}
			}
			switch name {
			case "month":
				returnTime.Months = float64(value)
			case "weeks":
				returnTime.Weeks = float64(value)
			case "day":
				returnTime.Days = float64(value)
			case "hour":
				returnTime.Hours = float64(value)
			case "minute":
				returnTime.Minutes = float64(value)
			case "second":
				returnTime.Seconds = float64(value)
			}
		}
	}
	return &returnTime, nil
}

func (g *GitlabAPI) GetTimeMR(mergeRequestID int) (*Overview, error) {
	return nil, nil
}

func extractTimeNotes(notes []*gitlab.Note, timeEntryRegex string) ([]*gitlab.Note, error) {
	var re = regexp.MustCompile(timeEntryRegex)
	timeNotes := []*gitlab.Note{}

	for _, note := range notes {
		if re.MatchString(string(note.Body)) {
			timeNotes = append(timeNotes, note)
		}
	}
	return timeNotes, nil
}
