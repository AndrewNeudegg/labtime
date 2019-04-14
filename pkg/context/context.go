package context

import (
	"fmt"
	"labtime/pkg/api"
	"labtime/pkg/config"

	log "github.com/sirupsen/logrus"
)

// Context exposes all methods to pongo2.
type Context struct {
	API       *api.GitlabAPI
	AppConfig *config.Configuration
}

func CreateContext(AppConfig *config.Configuration) *Context {
	return &Context{API: nil, AppConfig: AppConfig}
}

func (c *Context) Connect() (err error) {
	log.Info("Attempting to connect to Gitlab Instance.")
	if c.API != nil {
		return fmt.Errorf("API has already been instantiated")
	}
	log.Info("Building API Client.")
	c.API, err = api.Connect(c.AppConfig.Instance.Project, c.AppConfig.Instance.AccessToken, c.AppConfig.Instance.URL)
	log.Info("API Client built.")
	return err
}

// handleErr eases the use of errors.
func handleErr(err error) {
	if err != nil {
		log.Fatalf(err.Error())
	}
}
