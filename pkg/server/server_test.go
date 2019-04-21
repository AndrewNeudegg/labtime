package server2

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/phayes/freeport"
	"gotest.tools/assert"
)

// TODO: Add interactive test so that developers can see what is actually being returned.

func Test_ServerBasic(t *testing.T) {
	port, err := freeport.GetFreePort()
	if err != nil {
		t.Errorf(err.Error())
	}

	webApp := WebApp{
		Hostname: "localhost",
		Port:     port,
		Routes: map[string]Route{
			"/":     Route{Content: "Hello Home"},
			"/help": Route{Content: "Hello Help"},
		},
	}

	go func() {
		if err := webApp.Launch(); err != nil {
			t.Fatalf(err.Error())
		}
	}()

	for path, route := range webApp.Routes {
		fmt.Printf("Testing route: %s\n", path)
		resp, err := http.Get(fmt.Sprintf("http://%s:%v%s", webApp.Hostname, webApp.Port, path))
		if err != nil {
			t.Fatalf(err.Error())
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf(err.Error())
		}
		fmt.Printf("Testing that %s == %s\n", string(body), route.Content)
		assert.Assert(t, string(body) == route.Content)
	}
}
