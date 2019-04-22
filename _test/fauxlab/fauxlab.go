package fauxlab

import (
	"labtime/pkg/webapp"
	"path"
)

const (
	// ConfigurationRoot defines where the yaml web app configuration files should live.
	ConfigurationRoot string = "_test/fauxlab/configurations"
	// DefaultHost specifies where any calling code should look for the dummy webApp.
	DefaultHost string = "localhost"
	// DefaultPort specifies where any calling code should look for the dummy webApp.
	DefaultPort int = 8473
	// DefaultBreaker is the value that when passed to the channel for a method will terminate the server.
	DefaultBreaker int = 0
)

// Server is the fauxlab server.
type Server struct{}

// IssuesAPI will load and start an example API for the issues.
func (s *Server) IssuesAPI(breaker chan int) {
	IssuesAPIConfig := path.Join(ConfigurationRoot, "issues.yaml")
	apiHelper(IssuesAPIConfig, breaker)
}

func apiHelper(configPath string, breaker chan int) {
	webApp := webapp.WebApp{}
	webApp.Load(configPath)
	configureServer(&webApp)
	launchServer(&webApp)

	for {
		v := <-breaker
		if v == DefaultBreaker {
			break
		}
	}

	terminateServer(&webApp)
}

func launchServer(we *webapp.WebApp) {
	go func() {
		if err := we.Launch(); err != nil {
			panic(err)
		}
	}()
}

func terminateServer(we *webapp.WebApp) {

}

func configureServer(we *webapp.WebApp) {
	we.Hostname = DefaultHost
	we.Port = DefaultPort
}
