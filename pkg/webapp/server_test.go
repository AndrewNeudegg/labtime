package webapp

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/phayes/freeport"
	"gotest.tools/assert"
)

// Test_ServerBasic will test the servers start up and immediate response to http requests.
func Test_ServerBasic(t *testing.T) {
	setUp(t)
	webApp := WebApp{
		Hostname: "localhost",
		Port:     getFreePort(t),
		Routes: map[string]Route{
			"/":     Route{Content: *helloWorldDynamicContent()},
			"/help": Route{Content: *helloHelpDynamicContent()},
		},
	}

	go func() {
		if err := webApp.Launch(); err != nil {
			t.Fatalf(err.Error())
		}
	}()
	time.Sleep(0)

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
		jsonBytes, err := route.Content.toJSON()
		if err != nil {
			t.Fatalf(err.Error())
		}
		assert.Assert(t, string(body) == string(jsonBytes))
	}
	tearDown(t)
}

// Test_SaveNLoadServer will; create webApp, save it to file, load it to a new instance and then attempt to connect.
// upon connecting the values will be checked against the original (pre-save) values.
func Test_SaveNLoadServer(t *testing.T) {
	setUp(t)
	// Get a temporary file.
	file := getTempFile(t)
	defer os.Remove(file.Name())
	// Create an instance of the dummy web app.
	webApp1 := WebApp{
		Hostname: "localhost",
		Port:     getFreePort(t),
		Routes: map[string]Route{
			// Following JSON from: https://developer.mozilla.org/en-US/docs/Learn/JavaScript/Objects/JSON.
			"/":     Route{Content: *helloWorldDynamicContent()},
			"/help": Route{Content: *helloHelpDynamicContent()},
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
	time.Sleep(0)

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
		jsonBytes, err := route.Content.toJSON()
		if err != nil {
			t.Fatalf(err.Error())
		}
		assert.Assert(t, string(body) == string(jsonBytes))
	}
	tearDown(t)
}

func Test_StartStop(t *testing.T) {
	setUp(t)
	webApp := WebApp{
		Hostname: "localhost",
		Port:     getFreePort(t),
		Routes: map[string]Route{
			"/":     Route{Content: *helloWorldDynamicContent()},
			"/help": Route{Content: *helloHelpDynamicContent()},
		},
	}

	go func() {
		if err := webApp.Launch(); err != http.ErrServerClosed {
			t.Fatalf(err.Error())
		}
	}()
	time.Sleep(0)

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
		jsonBytes, err := route.Content.toJSON()
		if err != nil {
			t.Fatalf(err.Error())
		}
		assert.Assert(t, string(body) == string(jsonBytes))
	}

	if err := webApp.Terminate(); err != nil {
		t.Fatalf(err.Error())
	}

	_, err := http.Get(fmt.Sprintf("http://%s:%v%s", webApp.Hostname, webApp.Port, "/"))
	if err == nil {
		t.Fatalf("should have raised a dial error as the server has been terminated")
	}
	tearDown(t)
}

func Test_DynamicContentBasic(t *testing.T) {
	setUp(t)
	// Create some content that we will serialise.
	content := make(DynamicContent)
	content["hello"] = "world"
	// Marshal the content to JSON.
	contentAsJSONBytes, err := content.toJSON()
	if err != nil {
		t.Fatalf(err.Error())
	}
	// Check that the serialised content is as expected.
	assert.Assert(t, string(contentAsJSONBytes) == `{"hello":"world"}`)

	dupeContent := make(DynamicContent)
	if err := dupeContent.fromJSON(contentAsJSONBytes); err != nil {
		t.Fatalf(err.Error())
	}

	assert.DeepEqual(t, content, dupeContent)
	tearDown(t)
}

func Test_DynamicContentComplex(t *testing.T) {
	setUp(t)
	content := make(DynamicContent)
	err := content.fromJSON([]byte(`
	{
			"glossary": {
					"title": "example glossary",
			"GlossDiv": {
							"title": "S",
				"GlossList": {
									"GlossEntry": {
											"ID": "SGML",
						"SortAs": "SGML",
						"GlossTerm": "Standard Generalized Markup Language",
						"Acronym": "SGML",
						"Abbrev": "ISO 8879:1986",
						"GlossDef": {
													"para": "A meta-markup language, used to create markup languages such as DocBook.",
							"GlossSeeAlso": ["GML", "XML"]
											},
						"GlossSee": "markup"
									}
							}
					}
			}
	}
	`))
	if err != nil {
		t.Fatalf(err.Error())
	}

	contentAsJSONBytes, err := content.toJSON()
	if err != nil {
		t.Fatalf(err.Error())
	}

	dupeContent := make(DynamicContent)
	err = dupeContent.fromJSON(contentAsJSONBytes)
	if err != nil {
		t.Fatalf(err.Error())
	}

	assert.DeepEqual(t, content, dupeContent)
	tearDown(t)
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

// generateDynamicContent will generate simple to use example content.
func generateDynamicContent(input string) *DynamicContent {
	dc := make(DynamicContent)
	dc.fromJSON([]byte(input))
	return &dc
}

func helloWorldDynamicContent() *DynamicContent {
	return generateDynamicContent(`{ "hello":"world" }`)
}

func helloHelpDynamicContent() *DynamicContent {
	return generateDynamicContent(`{ "hello":"help" }`)
}

func complexDynamicContent() *DynamicContent {
	return generateDynamicContent(`
	{
		"glossary": {
				"title": "example glossary",
		"GlossDiv": {
						"title": "S",
			"GlossList": {
								"GlossEntry": {
										"ID": "SGML",
					"SortAs": "SGML",
					"GlossTerm": "Standard Generalized Markup Language",
					"Acronym": "SGML",
					"Abbrev": "ISO 8879:1986",
					"GlossDef": {
												"para": "A meta-markup language, used to create markup languages such as DocBook.",
						"GlossSeeAlso": ["GML", "XML"]
										},
					"GlossSee": "markup"
								}
						}
				}
		}
}
	`)
}

// setUp is run before tests
func setUp(t *testing.T) {
	fmt.Printf(">>> Starting: %v\n", t.Name())
}

// tearDown is run before tests
func tearDown(t *testing.T) {
	fmt.Printf(">>> Finished: %v\n", t.Name())
}
