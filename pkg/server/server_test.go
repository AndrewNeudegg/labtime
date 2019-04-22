package server2

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/phayes/freeport"
	"gotest.tools/assert"
)

// TODO: Add interactive test so that developers can see what is actually being returned.
// TODO: Add test that will examine the saving and loading of the web app routes.

func Test_ServerBasic(t *testing.T) {
	webApp := WebApp{
		Hostname: "localhost",
		Port:     getFreePort(t),
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

// Test_SaveNLoadServer will; create webApp, save it to file, load it to a new instance and then attempt to connect.
// upon connecting the values will be checked against the original (pre-save) values.
func Test_SaveNLoadServer(t *testing.T) {
	// Get a temporary file.
	file := getTempFile(t)
	defer os.Remove(file.Name())
	// Create an instance of the dummy web app.
	webApp1 := WebApp{
		Hostname: "localhost",
		Port:     getFreePort(t),
		Routes: map[string]Route{
			// Following JSON from: https://developer.mozilla.org/en-US/docs/Learn/JavaScript/Objects/JSON.
			"/": Route{Content: `{
				"squadName": "Super hero squad",
				"homeTown": "Metro City",
				"formed": 2016,
				"secretBase": "Super tower",
				"active": true,
				"members": [
				  {
					"name": "Molecule Man",
					"age": 29,
					"secretIdentity": "Dan Jukes",
					"powers": [
					  "Radiation resistance",
					  "Turning tiny",
					  "Radiation blast"
					]
				  },
				  {
					"name": "Madame Uppercut",
					"age": 39,
					"secretIdentity": "Jane Wilson",
					"powers": [
					  "Million tonne punch",
					  "Damage resistance",
					  "Superhuman reflexes"
					]
				  },
				  {
					"name": "Eternal Flame",
					"age": 1000000,
					"secretIdentity": "Unknown",
					"powers": [
					  "Immortality",
					  "Heat Immunity",
					  "Inferno",
					  "Teleportation",
					  "Interdimensional travel"
					]
				  }
				]
			  }`},
			// Following JSON from: https://json.org/example.html.
			"/help": Route{Content: `{"menu": {
				"id": "file",
				"value": "File",
				"popup": {
				  "menuitem": [
					{"value": "New", "onclick": "CreateNewDoc()"},
					{"value": "Open", "onclick": "OpenDoc()"},
					{"value": "Close", "onclick": "CloseDoc()"}
				  ]
				}
			  }}`},
		},
	}
	// Write out the webapp1 to a file.
	webApp1.Save(file.Name())
	// Attempt to load from the written file.
	webApp2 := WebApp{}
	webApp2.Load(file.Name())
	// Now attempt to launch the newly loaded webApp2 and check responses.
	go func() {
		if err := webApp2.Launch(); err != nil {
			t.Fatalf(err.Error())
		}
	}()
	// Check the new webApp2 against the values stored in webApp1.
	for path, route := range webApp1.Routes {
		fmt.Printf("Testing route: %s\n", path)
		resp, err := http.Get(fmt.Sprintf("http://%s:%v%s", webApp1.Hostname, webApp1.Port, path))
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

// getFreePort is a helper method for writing tests that may block on a specific port.
func getFreePort(t *testing.T) int {
	port, err := freeport.GetFreePort()
	if err != nil {
		t.Errorf(err.Error())
	}
	return port
}

// getTempFile is a helper method for writing tests that require saving and loading to file.
func getTempFile(t *testing.T) *os.File {
	file, err := ioutil.TempFile("", "server.*.yaml")
	if err != nil {
		t.Errorf(err.Error())
	}
	return file
}
