package handler_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"resty.dev/v3"
)

func TestGetUserURLs(t *testing.T) {
	t.Run("test get user urls", func(t *testing.T) {
		url, shutdown := setupTestServer()
		defer shutdown()
		client := resty.New()
		defer func() {
			_ = client.Close()
		}()

		resp, err := client.R().Get(url + "/api/user/urls")

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode())
	})
}
