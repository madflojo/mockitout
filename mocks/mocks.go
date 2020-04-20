/*
Package mocks is used to provide functionality for reading and parsing Mocks files.
These files are used to define the mock end-points to be loaded and virtualized by
this service.

The below is a sample mocks definition in YAML format.

  routes:
    hello:
      path: "/hi"
      response_headers:
        "content-type": "application/json"
        "server": "MockItOut"
      # Multi-line values can be created like this
      body: |
        {
          "greeting": "Hello",
          "name": "World"
        }

*/
package mocks

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

// Mocks defines the main mocks file structure.
type Mocks struct {
	// Routes is a map of Route values, each route is a mocked URI.
	Routes map[string]Route `yaml:"routes"`

	// Paths is a map of Path to Route names. This can be used to quickly lookup a
	// path and match it to a named route configuration.
	Paths map[string]string
}

// Route is the primary config for each mocked URI.
type Route struct {
	// Path is the URI value being mocked.
	Path string `yaml:"path"`

	// ResponseHeaders is a map of custom HTTP response headers.
	ResponseHeaders map[string]string `yaml:"response_headers"`

	// ReturnCode is the HTTP return code to reply with.
	ReturnCode int `yaml:"return_code"`

	// Body is the HTTP payload returned to be returned by the server.
	Body string `yaml:"body"`
}

// FromFile will read the Mocks file from the specified file path and return a
// Mocks configuration.
func FromFile(filepath string) (Mocks, error) {
	var m Mocks

	// Read file
	c, err := ioutil.ReadFile(filepath)
	if err != nil {
		return m, fmt.Errorf("could not read Mocks file at %s - %s", filepath, err)
	}

	// Parse the YAML
	err = yaml.Unmarshal(c, &m)
	if err != nil {
		return m, fmt.Errorf("error parsing Mocks file - %s", err)
	}

	// Check validity of Mocks file
	if len(m.Routes) < 1 {
		return m, fmt.Errorf("no routes defined in Mocks file")
	}

	// Setup helper values
	m.Paths = make(map[string]string)
	for k, v := range m.Routes {
		// Create lookup map
		m.Paths[v.Path] = k
	}

	return m, nil
}

// GenExampleFile will generate a basic example Mocks file. This is used
// primarily for testing.
func GenExampleFile() (*os.File, error) {
	data := []byte(`
routes:
  hello:
    path: "/hi"
    response_headers:
      "content-type": "application/json"
      "server": "MockItOut"
    # Multi-line values can be created like this
    body: | 
      {
        "greeting": "Hello",
        "name": "World"
      }
  deny:
    path: "/no"
    response_headers:
      "content-type": "application/json"
      "server": "MockItOut"
    body: |
      {"status": false}
    return_code: 403
`)

	fh, err := ioutil.TempFile("", "mocks_example")
	if err != nil {
		return fh, fmt.Errorf("Error creating temp file - %s", err)
	}

	_, err = fh.Write(data)
	if err != nil {
		return fh, fmt.Errorf("Error writing temp file data - %s", err)
	}

	err = fh.Close()
	if err != nil {
		return fh, fmt.Errorf("Error closing temp file - %s", err)
	}

	return fh, nil
}
