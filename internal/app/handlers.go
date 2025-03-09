package app

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	contentType = "text/plain; charset=utf-8"
)

// HandleOriginal converts provided url to the shorten by generating random Id.
// Returns 400 status code if user sent incorrect request's body, method or content-type.
func (urls *Urls) HandleOriginal(c *gin.Context) {
	if c.Request.Method == http.MethodPost {
		body, err := io.ReadAll(c.Request.Body)

		if err != nil || string(body) == "" {
			// w.WriteHeader(http.StatusBadRequest)
			c.Data(http.StatusBadRequest, contentType, nil)
			return
		}

		shortID := generateURLID(randomIDStringLength)

		urls.IDToURL[shortID] = string(body)
		shortURL := fmt.Sprintf("http://localhost:8080/%s", shortID)

		c.Data(http.StatusCreated, contentType, []byte(shortURL))

		return
	}

	c.Data(http.StatusBadRequest, contentType, nil)
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
	fmt.Println(id)

	originalURL := urls.IDToURL[id]
	if originalURL != "" {
		c.Header("Location", originalURL)
		c.Data(http.StatusTemporaryRedirect, contentType, nil)
		// w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	c.Data(http.StatusBadRequest, contentType, nil)
}