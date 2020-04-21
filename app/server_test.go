package app

import (
	"encoding/json"
	"github.com/madflojo/mockitout/config"
	"github.com/madflojo/mockitout/mocks"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"
)

// hiResp is the response JSON structure for the /hi end-point.
type hiResp struct {
	Greeting string `json:"greeting"`
	Name     string `json:"name"`
}

// noResp is the response JSON structure for the /no end-point.
type noResp struct {
	Status bool `json:"status"`
}

func TestMockServer(t *testing.T) {
	fh, err := mocks.GenExampleFile()
	if err != nil {
		t.Fatalf("Could not generate example Mocks file - %s", err)
	}
	defer os.Remove(fh.Name())

	go func() {
		err := Run(config.Config{
			Debug:          true,
			EnableTLS:      false,
			ListenAddr:     "localhost:9000",
			DisableLogging: true,
			MocksFile:      fh.Name(),
		})
		if err != nil && err != ErrShutdown {
			t.Errorf("Run unexpectedly stopped - %s", err)
		}
	}()
	// Clean up
	defer Stop()

	// Wait for app to start
	time.Sleep(10 * time.Second)

	t.Run("Check Hello Mock URL", func(t *testing.T) {
		r, err := http.Get("http://localhost:9000/hi")
		if err != nil {
			t.Errorf("Unexpected error when requesting mock URL - %s", err)
		}

		// Verify expected return code
		if r.StatusCode != 200 {
			t.Errorf("Unexpected http status code - %d", r.StatusCode)
		}

		// Verify expected headers
		if r.Header.Get("Server") != "MockItOut" {
			t.Errorf("Could not find expected Header Server - %+v", r.Header)
		}

		// Verify Body
		hi := hiResp{}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Unable to read HTTP response body - %s", err)
		}
		err = json.Unmarshal(body, &hi)
		if err != nil {
			t.Errorf("Unable to parse response JSON - %s", err)
		}
		if hi.Greeting != "Hello" {
			t.Errorf("Unexpected value from JSON - %+v", hi)
		}
	})

	for _, v := range []string{"Unk", "Andre", "Jim"} {
		t.Run("Check WildCard Mock URL with "+v, func(t *testing.T) {
			r, err := http.Get("http://localhost:9000/names/" + v)
			if err != nil {
				t.Errorf("Unexpected error when requesting mock URL - %s", err)
			}

			// Verify expected return code
			if r.StatusCode != 200 {
				t.Errorf("Unexpected http status code - %d", r.StatusCode)
			}

			// Verify expected headers
			if r.Header.Get("Server") != "WalkItOut" {
				t.Errorf("Could not find expected Header Server - %+v", r.Header)
			}
		})
	}

	t.Run("Check Deny Mock URL", func(t *testing.T) {
		r, err := http.Get("http://localhost:9000/no")
		if err != nil {
			t.Errorf("Unexpected error when requesting mock URL - %s", err)
		}

		// Verify expected return code
		if r.StatusCode != 403 {
			t.Errorf("Unexpected http status code - %d", r.StatusCode)
		}

		// Verify expected headers
		if r.Header.Get("Server") != "MockItOut" {
			t.Errorf("Could not find expected Header Server - %+v", r.Header)
		}

		// Verify Body
		no := noResp{}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Unable to read HTTP response body - %s", err)
		}
		err = json.Unmarshal(body, &no)
		if err != nil {
			t.Errorf("Unable to parse response JSON - %s", err)
		}
		if no.Status {
			t.Errorf("Unexpected value from JSON - %+v", no)
		}
	})
}
