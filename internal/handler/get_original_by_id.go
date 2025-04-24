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

	url, err := s.Service.GetOriginalURL(c.Request.Context(), id)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	if url.IsDeleted {
		c.Data(http.StatusGone, contentTypePlain, nil)
		return
	}

	if url.OriginalURL != "" {
		c.Header("Location", url.OriginalURL)
		c.Data(http.StatusTemporaryRedirect, contentTypePlain, nil)
		return
	}
	c.Data(http.StatusBadRequest, contentTypePlain, nil)
}
