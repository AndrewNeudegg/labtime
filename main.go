package main

import (
	"fmt"
	"labtime/pkg/config"
)

// main is the entry point of the application.
func main() {
	appConfig, err := config.Load("./config/custom/config.yml")
	handleErr(err)

	fmt.Println("Hello World!")
	fmt.Println(appConfig)
}

// handleErr eases the use of errors.
func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
