package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mp1947/ya-url-shortener/internal/usecase"
	"github.com/stretchr/testify/assert"
)

func TestGetOriginalURLByID(t *testing.T) {

	randomID := usecase.GenerateIDFromURL(testURL)

	userID := uuid.New().String()

	storage.Save(context.TODO(), randomID, testURL, userID) //nolint: errcheck

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
