package mocks

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestFromFile(t *testing.T) {
	// Testing with Valid Data
	data := make(map[string][]byte)
	data["valid yaml"] = []byte(`
routes:
  hello:
    path: "/hi"
    response_headers:
      "content-type": "application/json"
    # Multi-line values can be created like this
    body: '''
    {
      "greeting": "Hello",
      "name": "World"
    }
    '''
    return_code: 200
  `)

	for k, v := range data {
		t.Run("Testing "+k, func(t *testing.T) {

			fh, err := ioutil.TempFile("", "mocks_example")
			if err != nil {
				t.Fatalf("Error creating temp file - %s", err)
			}
			defer os.Remove(fh.Name())

			_, err = fh.Write(v)
			if err != nil {
				t.Fatalf("Error writing temp file data - %s", err)
			}

			err = fh.Close()
			if err != nil {
				t.Fatalf("Error closing temp file - %s", err)
			}

			_, err = FromFile(fh.Name())
			if err != nil {
				t.Fatalf("Unexpected Failure when loading valid yaml - %s", err)
			}
		})
	}

	// Testing with Invalid Data
	data = make(map[string][]byte)
	data["invalid yaml"] = []byte(`{"string":"bro this is json"}`)
	data["empty yaml"] = []byte("")
	data["broken yaml"] = []byte(`
routes:
  hello:
    path: "hi"
    response_headers:
      - "content-type": "application/json"
  `)

	for k, v := range data {
		t.Run("Testing "+k, func(t *testing.T) {

			fh, err := ioutil.TempFile("", "mocks_example")
			if err != nil {
				t.Fatalf("Error creating temp file - %s", err)
			}
			defer os.Remove(fh.Name())

			_, err = fh.Write(v)
			if err != nil {
				t.Fatalf("Error writing temp file data - %s", err)
			}

			err = fh.Close()
			if err != nil {
				t.Fatalf("Error closing temp file - %s", err)
			}

			_, err = FromFile(fh.Name())
			if err == nil {
				t.Fatalf("Expected Failure when loading invalid yaml, got nil")
			}
		})
	}

	t.Run("Testing with no file", func(t *testing.T) {
		_, err := FromFile("thisfilewillneverexistoratleastiwillalwaysthinkso")
		if err == nil {
			t.Fatalf("Expected Failure when loading missing file, got nil")
		}
	})
}

func TestGenExampleFile(t *testing.T) {
	fh, err := GenExampleFile()
	if err != nil {
		t.Fatalf("Unexpected error generating example YAML - %s", err)
	}
	defer os.Remove(fh.Name())

	_, err = FromFile(fh.Name())
	if err != nil {
		t.Fatalf("Unexpected error reading example YAML - %s", err)
	}
}
