package handler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mp1947/ya-url-shortener/config"
	"github.com/mp1947/ya-url-shortener/internal/repository"
	"github.com/mp1947/ya-url-shortener/internal/usecase"
)

const (
	randomIDStringLength = 8
	contentType          = "text/plain; charset=utf-8"
)

func ShortenURL(
	cfg config.Config,
	storage repository.Repository,
) gin.HandlerFunc {

	return func(c *gin.Context) {
		if c.Request.Method == http.MethodPost {
			body, err := io.ReadAll(c.Request.Body)

			if err != nil || string(body) == "" {
				c.Data(http.StatusBadRequest, contentType, nil)
				return
			}

			shortID := usecase.GenerateRandomID(randomIDStringLength)
			storage.Save(shortID, string(body))

			shortURL := fmt.Sprintf("%s/%s", *cfg.BaseURL, shortID)

			c.Data(http.StatusCreated, contentType, []byte(shortURL))
			return
		}

		c.Data(http.StatusBadRequest, contentType, nil)
	}
}

func GetOriginalURLByID(storage repository.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodGet {
			c.Data(http.StatusBadRequest, contentType, nil)
			return
		}

		id := c.Param("id")

		originalURL := storage.Get(id)
		if originalURL != "" {
			c.Header("Location", originalURL)
			c.Data(http.StatusTemporaryRedirect, contentType, nil)
			return
		}
		c.Data(http.StatusBadRequest, contentType, nil)
	}
}
