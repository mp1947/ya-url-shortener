package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mp1947/ya-url-shortener/internal/auth"
	"github.com/mp1947/ya-url-shortener/internal/middleware"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestAuthMiddleware(t *testing.T) {
	userID := uuid.New()
	validToken, _ := auth.CreateToken(userID)
	t.Logf("created test token: %s", validToken)

	tests := []struct {
		testName    string
		token       string
		userID      string
		handlerPath string
	}{
		{
			testName:    "valid token",
			token:       validToken,
			userID:      userID.String(),
			handlerPath: "/valid",
		},
		{
			testName:    "not valid token",
			token:       "blahblahblah",
			userID:      uuid.NewString(),
			handlerPath: "/nonvalid",
		},
	}

	r := gin.New()
	l := zap.New(zapcore.NewNopCore())
	r.Use(middleware.AuthMiddleware(l))

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			r.GET(test.handlerPath, func(ctx *gin.Context) {
				uid, exists := ctx.Get("user_id")
				assert.NotNil(t, exists)
				assert.NotNil(t, uid)
			})

			req := httptest.NewRequest(http.MethodGet, test.handlerPath, nil)
			req.AddCookie(&http.Cookie{
				Name:  "token",
				Value: test.token,
			})

			resp := httptest.NewRecorder()

			r.ServeHTTP(resp, req)

			assert.Equal(t, http.StatusOK, resp.Code)

		})
	}
}
