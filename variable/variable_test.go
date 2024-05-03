package variable

import (
	"net/http"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func createTestRequestContext() *RequestContext {
	r := RequestContext{}

	r.Request, _ = http.NewRequest("GET", "test:8080/url/:param?testquery=queryvalue", nil)
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
	r := createTestRequestContext()
	t.Setenv("testenv", "envvalue")

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
			inputVariable: "$randomFirstName",
			expectError:   false,
			expectValue:   "",
		},
		"Invalid Random": {
			inputVariable: "$badRandom",
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
			value, err := r.ParseVariable(tc.inputVariable)
			if len(tc.expectValue) > 0 {
				assert.Equal(t, tc.expectValue, value)
			}
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
