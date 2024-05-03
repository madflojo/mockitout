package variable

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func createTestRequestContext() *RequestContext {
	r := RequestContext{}

	RandomMap["randomMock"] = func() string {
		return "randomValue"
	}

	mockBody := strings.NewReader(`{"test": "body"}`)
	r.Request, _ = http.NewRequest("GET", "test:8080/url/:param?testquery=queryvalue", mockBody)
	r.Request.Header.Add("testheader", "headervalue")

	r.Params = make(httprouter.Params, 0)
	r.Params = append(r.Params, httprouter.Param{Key: "testparam", Value: "paramvalue"})

	return &r
}

func TestNewRequestContext(t *testing.T) {

	testRequest := &http.Request{}
	testParams := httprouter.Params{}

	r := NewRequestContext(testRequest, nil, testParams)

	assert.NotNil(t, r)
}

func TestParseVariable(t *testing.T) {
	testMatrix := map[string]struct {
		inputVariable string
		expectError   bool
		expectValue   string
	}{
		"Valid Header": {
			inputVariable: "header.testheader",
			expectError:   false,
			expectValue:   "headervalue",
		},
		"Invalid Header": {
			inputVariable: "header.badheader",
			expectError:   true,
			expectValue:   "",
		},
		"Valid Query": {
			inputVariable: "query.testquery",
			expectError:   false,
			expectValue:   "queryvalue",
		},
		"Invalid Query": {
			inputVariable: "query.badquery",
			expectError:   true,
			expectValue:   "",
		},
		"Valid Param": {
			inputVariable: "param.testparam",
			expectError:   false,
			expectValue:   "paramvalue",
		},
		"Invalid Param": {
			inputVariable: "param.badparam",
			expectError:   true,
			expectValue:   "",
		},
		"Valid Environment": {
			inputVariable: "environment.testenv",
			expectError:   false,
			expectValue:   "envvalue",
		},
		"Invalid Environment": {
			inputVariable: "environment.badenv",
			expectError:   true,
			expectValue:   "",
		},
		"Valid Random": {
			inputVariable: "$randomMock",
			expectError:   false,
			expectValue:   "randomValue",
		},
		"Invalid Random": {
			inputVariable: "$badRandom",
			expectError:   true,
			expectValue:   "",
		},
		"Valid Text Body": {
			inputVariable: "body",
			expectError:   false,
			expectValue:   `{"test": "body"}`,
		},
		"Valid Json Body": {
			inputVariable: "body.test",
			expectError:   false,
			expectValue:   "body",
		},
		"Invalid Json Body": {
			inputVariable: "body.bad",
			expectError:   true,
			expectValue:   "",
		},
		"Invalid Prefix": {
			inputVariable: "BadVariable",
			expectError:   true,
			expectValue:   "",
		},
		"Blank": {
			inputVariable: "",
			expectError:   true,
			expectValue:   "",
		},
	}

	for name, tc := range testMatrix {
		t.Run(name, func(t *testing.T) {
			r := createTestRequestContext()
			t.Setenv("testenv", "envvalue")

			value, err := r.ParseVariable(tc.inputVariable)
			assert.Equal(t, tc.expectValue, value)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetTestBody(t *testing.T) {
	testMatrix := map[string]struct {
		inputBody   string
		expectError bool
		expectValue string
	}{
		"Valid Body": {
			inputBody:   "test body",
			expectError: false,
			expectValue: "test body",
		},
		"Valid Json Body": {
			inputBody:   `{"test": "value"}`,
			expectError: false,
			expectValue: `{"test": "value"}`,
		},
		"Blank": {
			inputBody:   "",
			expectError: false,
			expectValue: "",
		},
	}

	for name, tc := range testMatrix {
		t.Run(name, func(t *testing.T) {
			r := createTestRequestContext()
			r.Request.Body = io.NopCloser(strings.NewReader(tc.inputBody))

			value, err := r.getTextBody("body")
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectValue, value)
			}
		})
	}
}

func TestGetBodyJsonVariable(t *testing.T) {
	testMatrix := map[string]struct {
		inputBody     string
		inputVariable string
		expectError   bool
		expectValue   string
	}{
		"Valid Json": {
			inputBody:     `{"test": "value"}`,
			inputVariable: "test",
			expectError:   false,
			expectValue:   "value",
		},
		"Valid Nested Json": {
			inputBody:     `{"test": {"nested": "value2"}}`,
			inputVariable: "test.nested",
			expectError:   false,
			expectValue:   "value2",
		},
		"Valid Json Map": {
			inputBody:     `{"test": {"nested": {"key": "value"}}}`,
			inputVariable: "test.nested",
			expectError:   false,
			expectValue:   `{"key":"value"}`,
		},
		"Valid Json Array": {
			inputBody:     `{"test": {"nested": ["value1", "value2"]}}`,
			inputVariable: "test.nested",
			expectError:   false,
			expectValue:   `["value1","value2"]`,
		},
		"Valid Json Int": {
			inputBody:     `{"test": 1}`,
			inputVariable: "test",
			expectError:   false,
			expectValue:   "1",
		},
		"Invalid Json": {
			inputBody:     `{"test": "value"`,
			inputVariable: "test",
			expectError:   true,
			expectValue:   "",
		},
		"Missing Json Variable": {
			inputBody:     `{"test": "value"}`,
			inputVariable: "bad",
			expectError:   true,
			expectValue:   "",
		},
		"Blank": {
			inputBody:     "",
			inputVariable: "test",
			expectError:   true,
			expectValue:   "",
		},
	}

	for name, tc := range testMatrix {
		t.Run(name, func(t *testing.T) {
			r := createTestRequestContext()
			r.Request.Body = io.NopCloser(strings.NewReader(tc.inputBody))

			value, err := r.getBodyJsonVariable(tc.inputVariable)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectValue, value)
			}
		})
	}
}
