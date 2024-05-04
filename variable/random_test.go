package variable

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRandomVariable(t *testing.T) {
	testMatrix := map[string]struct {
		inputVariable string
		expectError   bool
	}{
		"Valid randomFirstName": {
			inputVariable: "randomFirstName",
			expectError:   false,
		},
		"Invalid": {
			inputVariable: "BadVariable",
			expectError:   true,
		},
	}

	for name, tc := range testMatrix {
		t.Run(name, func(t *testing.T) {
			_, err := getRandomVariable(tc.inputVariable)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRandomMap(t *testing.T) {
	for name, tc := range RandomMap {
		t.Run(name, func(t *testing.T) {
			val := tc()
			assert.NotEmpty(t, val)
		})
	}
}
