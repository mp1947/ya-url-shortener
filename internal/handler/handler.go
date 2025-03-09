package handler

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mp1947/ya-url-shortener/config"
	"github.com/mp1947/ya-url-shortener/internal/service"
)

const (
	contentType = "text/plain; charset=utf-8"
)

type HandlerService struct {
	Service service.Service
	Cfg     config.Config
}

func (s HandlerService) ShortenURL(c *gin.Context) {

	if c.Request.Method != http.MethodPost {
		c.Data(http.StatusBadRequest, contentType, nil)
		return
	}

	body, err := io.ReadAll(c.Request.Body)

	if err != nil || string(body) == "" {
		c.Data(http.StatusBadRequest, contentType, nil)
		return
	}

	shortURL := s.Service.ShortenURL(s.Cfg, string(body))

	c.Data(http.StatusCreated, contentType, []byte(shortURL))

}

func (s HandlerService) GetOriginalURLByID(c *gin.Context) {
	if c.Request.Method != http.MethodGet {
		c.Data(http.StatusBadRequest, contentType, nil)
		return
	}

	id := c.Param("id")

	originalURL := s.Service.GetOriginalURL(id)
	if originalURL != "" {
		c.Header("Location", originalURL)
		c.Data(http.StatusTemporaryRedirect, contentType, nil)
		return
	}
	c.Data(http.StatusBadRequest, contentType, nil)
}
