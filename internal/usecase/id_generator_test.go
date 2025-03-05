package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRandomID(t *testing.T) {
	tests := []struct {
		testName             string
		inputLength          int
		expectedOutputLength int
	}{
		{
			testName:             "test returned string length",
			inputLength:          10,
			expectedOutputLength: 10,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			urlID := GenerateRandomID(test.inputLength)
			assert.Equal(t, test.expectedOutputLength, len(urlID))
		})
	}
}
