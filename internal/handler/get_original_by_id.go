package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s HandlerService) GetOriginalURLByID(c *gin.Context) {
	if c.Request.Method != http.MethodGet {
		c.Data(http.StatusBadRequest, contentTypePlain, nil)
		return
	}

	id := c.Param("id")

	originalURL, err := s.Service.GetOriginalURL(c.Request.Context(), id)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	if originalURL != "" {
		c.Header("Location", originalURL)
		c.Data(http.StatusTemporaryRedirect, contentTypePlain, nil)
		return
	}
	c.Data(http.StatusBadRequest, contentTypePlain, nil)
}
