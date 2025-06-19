package handler_test

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestDeleteUserURLS(t *testing.T) {

	gin.SetMode(gin.TestMode)
	t.Run("test delete user urls", func(t *testing.T) {
		// 	w := httptest.NewRecorder()
		// 	c, _ := gin.CreateTestContext(w)

		// 	c.Request = httptest.NewRequest(http.MethodDelete, "/", nil)
		// 	hs.DeleteUserURLs(c)

		// 	result := w.Result()
		// 	respStatusCode := result.StatusCode
		// 	defer result.Body.Close()
		// 	t.Logf("unauthenticated delete request returns %v status code", respStatusCode)
		// 	assert.Equal(t, respStatusCode, http.StatusUnauthorized)
	})
}
