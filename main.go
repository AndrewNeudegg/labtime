package main

import (
	"fmt"
	"labtime/pkg/api"
	"labtime/pkg/config"
	"log"
)

// main is the entry point of the application.
func main() {
	appConfig, err := config.Load("./config/custom/config.yml")
	handleErr(err)

	fmt.Println("Hello World!")
	fmt.Println(appConfig)

	gapi, err := api.Connect(appConfig.Instance.Project, appConfig.Instance.AccessToken, appConfig.Instance.URL)
	handleErr(err)

	// users, err := gapi.GetUsers()
	// handleErr(err)
	// for _, user := range users {
	// 	fmt.Println(user.Name)
	// }

	// issueID := 0
	issues, err := gapi.GetIssues()
	for _, issue := range issues {
		fmt.Println(issue.WebURL)
	}
	notes, err := gapi.GetIssueNotes(133)
	handleErr(err)
	for _, note := range notes {
		fmt.Println(note.NoteableType)
		fmt.Println(note.String())
	}
}

// handleErr eases the use of errors.
func handleErr(err error) {
	if err != nil {
		log.Fatalf(err.Error())
	}
}
