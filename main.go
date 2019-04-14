package main

import (
	"io/ioutil"
	"labtime/pkg/config"
	"labtime/pkg/context"
	"labtime/pkg/templator"
	"os"

	log "github.com/sirupsen/logrus"
)

// main is the entry point of the application.
func main() {
	log.Info("Starting application...")

	err := generateDefualtConfig()
	handleErr(err)

	appConfig, err := config.Load("./config/custom/config.yml")
	handleErr(err)

	ud := make(map[string]interface{})
	ctx := context.CreateContext(appConfig)
	handleErr(ctx.Connect())
	ud["ctx"] = ctx

	// output, err := templator.RenderTemplate("./templates/IssueOverview.csv.j2", templator.CreateContext(ud))
	output, err := templator.RenderTemplate("./templates/MROverview.csv.j2", templator.CreateContext(ud))
	handleErr(err)
	err = ioutil.WriteFile("./output/MROverview.csv", []byte(output), 0644)
	handleErr(err)
}

// handleErr eases the use of errors.
func handleErr(err error) {
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func generateDefualtConfig() error {
	log.Info("Generating the default configuration.")
	appConfig := config.Default()
	defer log.Info("Generated the default configuration to: ./config/default/config.yml")
	return appConfig.Save("./config/default/config.yml")
}

// init: configuration
func init() {
	// Log as JSON instead of the default ASCII formatter.
	// log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}
