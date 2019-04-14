package main

import (
	"fmt"
	"labtime/pkg/api"
	"labtime/pkg/config"
	"log"
)

// main is the entry point of the application.
func main() {
	err := generateDefualtConfig()
	handleErr(err)

	appConfig, err := config.Load("./config/custom/config.yml")
	handleErr(err)

	// fmt.Println("Hello World!")
	// fmt.Println(appConfig)

	gapi, err := api.Connect(appConfig.Instance.Project, appConfig.Instance.AccessToken, appConfig.Instance.URL)
	handleErr(err)

	// users, err := gapi.GetUsers()
	// handleErr(err)
	// for _, user := range users {
	// 	fmt.Println(user.Name)
	// }

	// issueID := 0
	// issues, err := gapi.GetIssues()
	// for _, issue := range issues {
	// 	fmt.Println(issue.WebURL)
	// }
	// notes, err := gapi.GetIssueNotes(133)
	// handleErr(err)
	// for _, note := range notes {
	// 	fmt.Println(note.NoteableType)
	// 	fmt.Println(note.String())
	// }

	overview, err := gapi.GetTimeIssue(
		133, appConfig.QueryConfig.TimeEntryDetectionRegex, appConfig.QueryConfig.TimeEntryExtractionRegex)
	handleErr(err)
	totalTime := float64(0)
	for _, timeSpent := range overview.TimeSpents {
		fmt.Printf("%s: %v days (raw: %v) @ %v - %v :: %v \n",
			timeSpent.Author.Name,
			timeSpent.Spent.TotalDaysRounded(),
			timeSpent.Spent.TotalDaysRaw(),
			timeSpent.CreatedAt.Format("02/01/06"),
			timeSpent.RawBody,
			timeSpent.Spent)
		totalTime += timeSpent.Spent.TotalDaysRounded()
	}
	fmt.Printf("Total time spent: %v days \n", totalTime)
}

// handleErr eases the use of errors.
func handleErr(err error) {
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func generateDefualtConfig() error {
	appConfig := config.Default()
	return appConfig.Save("./config/default/config.yml")
}
