package variable

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplaceVariables(t *testing.T) {
	testMatrix := map[string]struct {
		inputData   string
		expectError bool
		expectValue string
	}{
		"Valid Header": {
			inputData:   "Test Header Data: {{header.testheader}}",
			expectError: false,
			expectValue: "Test Header Data: headervalue",
		},
		"Invalid Header": {
			inputData:   "Test Header Data: {{header.badheader}}",
			expectError: false,
			expectValue: "Test Header Data: {{header.badheader}}",
		},
		"Valid Query": {
			inputData:   "Test Query Data: {{query.testquery}}",
			expectError: false,
			expectValue: "Test Query Data: queryvalue",
		},
		"Invalid Query": {
			inputData:   "Test Query Data: {{query.badquery}}",
			expectError: false,
			expectValue: "Test Query Data: {{query.badquery}}",
		},
		"Valid Param": {
			inputData:   "Test Param Data: {{param.testparam}}",
			expectError: false,
			expectValue: "Test Param Data: paramvalue",
		},
		"Invalid Param": {
			inputData:   "Test Param Data: {{param.badparam}}",
			expectError: false,
			expectValue: "Test Param Data: {{param.badparam}}",
		},
		"Valid Environment": {
			inputData:   "Test Environment Data: {{environment.testenv}}",
			expectError: false,
			expectValue: "Test Environment Data: envvalue",
		},
		"Invalid Environment": {
			inputData:   "Test Environment Data: {{environment.badenv}}",
			expectError: false,
			expectValue: "Test Environment Data: {{environment.badenv}}",
		},
		"Valid Double Variable": {
			inputData:   "Test Double Variable Data: {{header.testheader}} {{query.testquery}}",
			expectError: false,
			expectValue: "Test Double Variable Data: headervalue queryvalue",
		},
		"Invalid Random": {
			inputData:   "Test Random Data: {{$badRandom}}",
			expectError: false,
			expectValue: "Test Random Data: {{$badRandom}}",
		},
		"Invalid Braces: Missing Close": {
			inputData:   "Test Random Data: {{$badRandom",
			expectError: false,
			expectValue: "Test Random Data: {{$badRandom",
		},
		"Invalid Braces: Missing Open": {
			inputData:   "Test Random Data: $}}",
			expectError: false,
			expectValue: "Test Random Data: $}}",
		},
		"Invalid Braces: Missing Both": {
			inputData:   "Test Random Data: {$badRandom}",
			expectError: false,
			expectValue: "Test Random Data: {$badRandom}",
		},
		"Invalid No Variable": {
			inputData:   "Test Random Data: {{}}",
			expectError: false,
			expectValue: "Test Random Data: {{}}",
		},
		"Blank": {
			inputData:   "",
			expectError: false,
			expectValue: "",
		},
	}

	for name, tc := range testMatrix {
		t.Run(name, func(t *testing.T) {
			r := createTestRequestContext()
			t.Setenv("testenv", "envvalue")

			res, err := r.ReplaceVariables(tc.inputData)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectValue, res)
			}
		})
	}
}
