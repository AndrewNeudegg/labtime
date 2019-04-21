package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// IntegrationTest defines all integration test methods.
// type IntegrationTest struct{}

// DemoServer defines methods avaliable for the demo server.
type DemoServer struct {
	config *DemoServerConfig // config is the current configuration.
	server *http.Server      // server is the current server.
	err    error             // err is any error that occurs within the server.
}

// DemoServerConfig offers some configuration parameters for the example server.
type DemoServerConfig struct {
	Hostname string `json:"hostname,omitempty"` // Hostname specifies the host that the demo server will be served on.
	Port     int    `json:"port,omitempty"`     // Port specifies the TCP port that the server will open on.
	Path     string `json:"path,omitempty"`     // Path defines the API path that this response should exist on.
	Content  string `json:"content,omitempty"`  // Content is the response that will return.
}

// ServeHTTP will enable the DemoServerConfig struct to be used as the request handler.
func (config *DemoServerConfig) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// for sanity we will 404 on unknown path.
	// TODO: Extend serve to handle multiple paths.
	if r.URL.Path != config.Path {
		http.Error(w, "404 Not Found", http.StatusNotFound)
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, config.Content)
}

// CreateDemoServer will create an instance of the demo server.
func CreateDemoServer(config *DemoServerConfig) (*DemoServer, error) {
	return &DemoServer{config: config}, nil
}

// Start the server.
func (ds *DemoServer) Start() error {
	ds.server = &http.Server{
		Addr:    ds.formatAddr(ds.config.Hostname, ds.config.Port, ds.config.Path),
		Handler: ds.config,
	}

	// create an error so that we can politely said error message to the user.
	var handleErr error

	go func() {
		if err := ds.server.ListenAndServe(); err != nil {
			panic(err)
			// return
		}
	}()
	// Trigger the goroutine.
	time.Sleep(0)
	return handleErr
}

// Stop will stop a currently running server.
func (ds *DemoServer) Stop() error {
	if ds.server == nil {
		return fmt.Errorf("cannot stop a nil server")
	}

	// Setting up signal capturing
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Waiting for SIGINT (pkill -2)
	<-stop
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := ds.server.Shutdown(ctx)
	return err
}

// formatAddr is a helper method for formatting the URL properly.
func (ds DemoServer) formatAddr(host string, port int, path string) string {
	return fmt.Sprintf("%s:%v/%s", host, port, path)
}

// // LaunchDemoServer will open a simplified example of the Gitlab rest API.
// func (it *IntegrationTest) LaunchDemoServer() {
// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Fprintf(w, "Demo site should only be accessed by ")
// 	})

// 	fs := http.FileServer(http.Dir("static/"))
// 	http.Handle("/static/", http.StripPrefix("/static/", fs))

// 	http.ListenAndServe(":80", nil)
// }
