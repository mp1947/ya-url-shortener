package handler

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mp1947/ya-url-shortener/config"
	"github.com/mp1947/ya-url-shortener/internal/logger"
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
var fileStoragePath = "./test.out"
var cfg = config.Config{
	ListenAddr:      &listenAddr,
	BaseURL:         &baseURL,
	FileStoragePath: &fileStoragePath,
}
var storage = &inmemory.Memory{}
var storageInitializedErr = storage.Init(cfg, context.Background())
var hs = initTestHandlerService()

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

	gin.SetMode(gin.TestMode)

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(test.request.httpMethod, "/", test.request.requestBody)

			t.Logf("sending %s request to %s", c.Request.Method, c.Request.RequestURI)

			hs.ShortenURL(c)

			result := w.Result()

			body := result.Body
			defer body.Close() //nolint:errcheck

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

	storage.Save(randomID, testURL)

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

			hs.GetOriginalURLByID(c)

			result := w.Result()

			respStatusCode := result.StatusCode
			location := result.Header.Get("Location")

			defer result.Body.Close() //nolint:errcheck

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
			defer body.Close() //nolint:errcheck

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

	if storageInitializedErr != nil {
		log.Fatalf("error initializing storage: %v", storageInitializedErr)
	}

	logger, _ := logger.InitLogger()

	service := service.ShortenService{Storage: storage, Logger: logger}

	return HandlerService{Service: &service, Cfg: cfg}
}
