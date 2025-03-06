package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testURL = "https://console.yandex.cloud/"
)

func TestGenerateIDFromURL(t *testing.T) {
	tests := []struct {
		testName             string
		inputLength          int
		expectedOutputLength int
	}{
		{
			testName:             "test returned string length",
			expectedOutputLength: resultStrLength,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			urlID := GenerateIDFromURL(testURL)
			assert.Equal(t, test.expectedOutputLength, len(urlID))
		})
	}
}
