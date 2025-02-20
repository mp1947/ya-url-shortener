package app

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testURL = "https://console.yandex.cloud/"
)

func TestHandleOriginal(t *testing.T) {

	type request struct {
		httpMethod  string
		requestBody io.Reader
	}

	tests := []struct {
		testName                string
		request                 request
		expectedRespCode        int
		expectedRespContentType string
	}{
		{
			testName: "test wrong http method",
			request: request{
				httpMethod:  http.MethodGet,
				requestBody: nil,
			},
			expectedRespCode: http.StatusBadRequest,
		},
		{
			testName: "test empty body",
			request: request{
				httpMethod:  http.MethodPost,
				requestBody: nil,
			},
			expectedRespCode: http.StatusBadRequest,
		},
		{
			testName: "test correct request",
			request: request{
				httpMethod:  http.MethodPost,
				requestBody: strings.NewReader(testURL),
			},
			expectedRespCode: http.StatusCreated,
		},
	}

	urls := &Urls{IDToURL: map[string]string{}}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			r := httptest.NewRequest(test.request.httpMethod, "/", test.request.requestBody)
			w := httptest.NewRecorder()

			urls.HandleOriginal(w, r)

			body := w.Result().Body

			defer body.Close()

			bodyData, err := io.ReadAll(body)
			statusCode := w.Result().StatusCode

			require.NoError(t, err)

			if statusCode == http.StatusCreated {
				t.Logf("response body is: %v", string(bodyData))
				assert.NotEmpty(t, bodyData)
			}

			assert.Equal(t, test.expectedRespCode, w.Result().StatusCode)
		})
	}
}

func TestHandleShort(t *testing.T) {

	randomID := generateURLID(randomIDStringLength)

	urls := &Urls{IDToURL: map[string]string{
		randomID: testURL,
	}}

	type request struct {
		httpMethod    string
		originalURLID string
	}
	tests := []struct {
		testName           string
		request            request
		expectedStatusCode int
		expectedLocation   string
	}{
		{
			testName: "test incorrect id",
			request: request{
				httpMethod:    http.MethodGet,
				originalURLID: "/doesnotexists",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedLocation:   "",
		},
		{
			testName: "test correct id",
			request: request{
				httpMethod:    http.MethodGet,
				originalURLID: randomID,
			},
			expectedStatusCode: http.StatusTemporaryRedirect,
			expectedLocation:   testURL,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			r := httptest.NewRequest(test.request.httpMethod, "/", nil)
			r.SetPathValue("id", test.request.originalURLID)

			w := httptest.NewRecorder()

			urls.HandleShort(w, r)

			respStatusCode := w.Result().StatusCode

			assert.Equal(t, test.expectedStatusCode, respStatusCode)

			if respStatusCode == http.StatusTemporaryRedirect {
				location := w.Result().Header.Get("Location")
				assert.Equal(t, test.expectedLocation, location)
			}
		})
	}
}
