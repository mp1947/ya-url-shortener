package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mp1947/ya-url-shortener/config"
	"github.com/mp1947/ya-url-shortener/internal/repository/inmemory"
	"github.com/mp1947/ya-url-shortener/internal/service"
	"github.com/mp1947/ya-url-shortener/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testURL         = "https://console.yandex.cloud/"
	testJSONRequest = `{"url": "https://console.yandex.cloud"}`
)

var listenAddr = ":8080"
var baseURL = "http://localhost:8080"
var cfg = config.Config{
	ListenAddr: &listenAddr,
	BaseURL:    &baseURL,
}

func TestShortenURL(t *testing.T) {

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

	h := initTestHandlerService()
	gin.SetMode(gin.TestMode)

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(test.request.httpMethod, "/", test.request.requestBody)

			t.Logf("sending %s request to %s", c.Request.Method, c.Request.RequestURI)

			h.ShortenURL(c)

			result := w.Result()

			body := result.Body
			defer body.Close()

			bodyData, err := io.ReadAll(body)
			statusCode := result.StatusCode

			require.NoError(t, err)

			if statusCode == http.StatusCreated {
				t.Logf("response body is: %v", string(bodyData))
				assert.NotEmpty(t, bodyData)
			}

			assert.Equal(t, test.expectedRespCode, statusCode)
		})
	}
}

func TestGetOriginalURLByID(t *testing.T) {

	randomID := usecase.GenerateIDFromURL(testURL)

	storage := &inmemory.Memory{}
	storage.Init()
	storage.Save(randomID, testURL)

	service := service.ShortenService{Storage: storage}
	h := HandlerService{Service: service}

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
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(test.request.httpMethod, "/", nil)
			c.Params = []gin.Param{
				{
					Key:   "id",
					Value: test.request.originalURLID,
				},
			}

			h.GetOriginalURLByID(c)

			result := w.Result()

			respStatusCode := result.StatusCode
			location := result.Header.Get("Location")

			defer result.Body.Close()

			assert.Equal(t, test.expectedStatusCode, respStatusCode)

			if respStatusCode == http.StatusTemporaryRedirect {
				assert.Equal(t, test.expectedLocation, location)
			}
		})
	}
}

func TestJSONShortenURL(t *testing.T) {

	shortenPath := "/api/shorten"

	type request struct {
		httpMethod  string
		requestBody io.Reader
		path        string
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

	h := initTestHandlerService()

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

			h.JSONShortenURL(c)

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

func initTestHandlerService() HandlerService {

	storage := &inmemory.Memory{}
	storage.Init()

	service := service.ShortenService{Storage: storage}

	return HandlerService{Service: service, Cfg: cfg}
}
