package webapp

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"gopkg.in/yaml.v2"
)

// WebApp represents a dummy webapp for use in testing.
type WebApp struct {
	Hostname string           `json:"hostname,omitempty"` // Hostname specifies the host that the demo server will be served on.
	Port     int              `json:"port,omitempty"`     // Port specifies the TCP port that the server will open on.
	Routes   map[string]Route `json:"routes,omitempty"`   // Routes specifies the different URL offerings for this demo web app.
	server   *http.Server     // server is where the http server is stored while in use.
}

// Route specifies a URL Path and content that this web application will expose.
type Route struct {
	Content string `json:"content,omitempty"` // Content specifies the content that will live on this URL.
}

// routeHandler allows each route to supply the required method for serving its own content.
func (we *WebApp) routeHandler(wr http.ResponseWriter, re *http.Request) {
	fmt.Fprintf(wr, we.Routes[re.URL.Path].Content)
}

// Launch will start a web app.
func (we *WebApp) Launch() error {
	if we.server != nil {
		return fmt.Errorf("the server cannot be started multiple times")
	}
	we.server = &http.Server{Addr: buildAddress(we.Hostname, we.Port)}
	http.HandleFunc("/", we.routeHandler)
	return we.server.ListenAndServe()
}

// Terminate should halt a server that is already running.
func (we *WebApp) Terminate() error {
	if we.server == nil {
		return fmt.Errorf("cannot terminate an empty server")
	}
	if err := we.server.Shutdown(context.TODO()); err != nil {
		return err
	}
	we.server = nil
	return nil
}

// Save will write the struct to yaml and is a contained convenience method.
func (we *WebApp) Save(path string) error {
	y, err := yaml.Marshal(we)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, []byte(y), 0644)
}

// Load is a contained convenience method
func (we *WebApp) Load(path string) error {
	fileAsBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(fileAsBytes, &we)
}

func buildAddress(host string, port int) string {
	return fmt.Sprintf("%s:%v", host, port)
}
