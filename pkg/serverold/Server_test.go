package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"gotest.tools/assert"
)

func TestCreateDemoServer(t *testing.T) {
	// For now we will just test the sunshine path.
	server, err := CreateDemoServer(&DemoServerConfig{
		Hostname: "0.0.0.0",
		Port:     8080,
		Path:     "example/url",
		Content:  `{"someKey":"someValue"}`,
	})

	if err != nil {
		t.Fatalf(err.Error())
	}

	if server == nil {
		t.Fatalf("CreateDemoServer returned nil")
	}
}

func TestStart(t *testing.T) {
	server, err := CreateDemoServer(&DemoServerConfig{
		Hostname: "0.0.0.0",
		Port:     8080,
		Path:     "example/url",
		Content:  `{"someKey":"someValue"}`,
	})

	if err := server.Start(); err != nil {
		t.Fatalf(err.Error())
	}
	// time.Sleep(time.Duration(time.Second * 0.5))
	resp, err := http.Get(fmt.Sprintf("http://%s:%v/%s",
		server.config.Hostname, server.config.Port, server.config.Path))
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf(err.Error())
	}
	assert.Assert(t, string(body) == server.config.Content)
}
