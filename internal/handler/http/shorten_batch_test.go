package handlehttp_test

import (
	"net/http"
	"testing"

	"github.com/mp1947/ya-url-shortener/internal/dto"
	"github.com/stretchr/testify/assert"
	"resty.dev/v3"
)

func TestBatchShortenURL(t *testing.T) {
	url, shutdown := setupTestServer()
	defer shutdown()
	url += "/api/shorten/batch"

	tests := []struct {
		requestData      any
		testName         string
		expectedRespCode int
	}{
		{
			testName:         "test bad request",
			requestData:      `{"name": "whatever"}`,
			expectedRespCode: http.StatusBadRequest,
		},
		{
			testName: "test request with batch of urls",
			requestData: []dto.BatchShortenRequest{
				{
					CorrelationID: "0maasdasd",
					OriginalURL:   "https://google.com",
				},
				{
					CorrelationID: "0maasdasd",
					OriginalURL:   "https://yandex.ru",
				},
			},
			expectedRespCode: http.StatusCreated,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			client := resty.New()
			defer func() {
				_ = client.Close()
			}()

			resp, err := client.R().
				SetBody(test.requestData).
				Post(url)

			assert.NoError(t, err)
			assert.Equal(t, test.expectedRespCode, resp.StatusCode())
		})
	}

}
