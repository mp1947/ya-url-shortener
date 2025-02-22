package app

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mp1947/ya-url-shortener/config"
)

// HandleOriginal converts provided url to the shorten by generating random Id.
// Returns 400 status code if user sent incorrect request's body, method or content-type.
func (urls *Urls) HandleOriginal(config config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodPost {
			body, err := io.ReadAll(c.Request.Body)

			if err != nil || string(body) == "" {
				c.Data(http.StatusBadRequest, contentType, nil)
				return
			}

			shortID := generateURLID(randomIDStringLength)

			urls.IDToURL[shortID] = string(body)
			shortURL := fmt.Sprintf(
				"http://%s%s%s",
				*config.ListenAddr,
				*config.BasePath,
				shortID,
			)

			c.Data(http.StatusCreated, contentType, []byte(shortURL))

			return
		}

		c.Data(http.StatusBadRequest, contentType, nil)
	}
}

// HandleShort retrieves id from the GET request, looks for
// corresponding url in urls.IDToURL map and redirects to this url via 307 HTTP to the new Location.
// If url wasn't found in urls.IDToURL map or request is incorrect - returns 400 HTTP status code.
func (urls *Urls) HandleShort(c *gin.Context) {
	if c.Request.Method != http.MethodGet {
		c.Data(http.StatusBadRequest, contentType, nil)
		return
	}

	id := c.Param("id")

	originalURL := urls.IDToURL[id]
	if originalURL != "" {
		c.Header("Location", originalURL)
		c.Data(http.StatusTemporaryRedirect, contentType, nil)
		// w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	c.Data(http.StatusBadRequest, contentType, nil)
}
