package handlehttp_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONShortenURL(t *testing.T) {

	shortenPath := "/api/shorten"

	type request struct {
		httpMethod  string
		requestBody io.Reader
		path        string
	}

	tests := []struct {
		request                 request
		testName                string
		expectedRespContentType string
		expectedRespCode        int
	}{
		{
			testName: "test wrong http method",
			request: request{
				httpMethod:  http.MethodGet,
				requestBody: nil,
				path:        shortenPath,
			},
			expectedRespCode: http.StatusBadRequest,
		},
		{
			testName: "test empty body",
			request: request{
				httpMethod:  http.MethodPost,
				requestBody: nil,
				path:        shortenPath,
			},
			expectedRespCode: http.StatusBadRequest,
		},
		{
			testName: "test correct request",
			request: request{
				httpMethod:  http.MethodPost,
				requestBody: strings.NewReader(testJSONRequest),
				path:        shortenPath,
			},
			expectedRespCode: http.StatusCreated,
		},
	}

	gin.SetMode(gin.TestMode)

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(
				tc.request.httpMethod,
				tc.request.path,
				tc.request.requestBody,
			)

			t.Logf("sending %s request to %s", c.Request.Method, c.Request.RequestURI)

			hs.JSONShortenURL(c)

			result := w.Result()

			body := result.Body
			defer body.Close()

			bodyData, err := io.ReadAll(body)
			statusCode := result.StatusCode

			require.NoError(t, err)

			if statusCode == http.StatusCreated {
				t.Logf("json response body is: %v", string(bodyData))
				assert.NotEmpty(t, bodyData)
			}

			assert.Equal(t, tc.expectedRespCode, statusCode)
		})
	}
}
