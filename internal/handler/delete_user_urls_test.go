package handler_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"resty.dev/v3"
)

func TestDeleteUserURLS(t *testing.T) {
	t.Run("test bad request", func(t *testing.T) {
		baseURL, shutdown := setupTestServer()
		defer shutdown()
		client := resty.New()
		defer client.Close() //nolint: errcheck
		ids := []string{"12343213", "12345322"}

		resp, err := client.R().
			SetBody(ids).
			Delete(baseURL + "/api/user/urls")

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode())
	})
}
