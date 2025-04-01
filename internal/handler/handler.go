package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mp1947/ya-url-shortener/config"
	"github.com/mp1947/ya-url-shortener/internal/dto"
	shrterr "github.com/mp1947/ya-url-shortener/internal/errors"
	"github.com/mp1947/ya-url-shortener/internal/repository"
	"github.com/mp1947/ya-url-shortener/internal/service"
)

const (
	contentType       = "text/plain; charset=utf-8"
	requestBindingErr = "invalid request: error parsing request params"
	requestBodyGetErr = "error getting request body"
)

type HandlerService struct {
	Service service.Service
	Cfg     config.Config
	Storage repository.Repository
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

	shortURL, err := s.Service.ShortenURL(s.Cfg, string(body))

	if errors.Is(err, shrterr.ErrOriginalURLAlreadyExists) {
		c.Data(http.StatusConflict, contentType, []byte(shortURL))
		return
	} else if err != nil {
		c.Data(
			http.StatusInternalServerError,
			contentType,
			[]byte("internal server error while shorten url"),
		)
		return
	}

	c.Data(http.StatusCreated, contentType, []byte(shortURL))

}

func (s HandlerService) GetOriginalURLByID(c *gin.Context) {
	if c.Request.Method != http.MethodGet {
		c.Data(http.StatusBadRequest, contentType, nil)
		return
	}

	id := c.Param("id")

	originalURL, err := s.Service.GetOriginalURL(id)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	if originalURL != "" {
		c.Header("Location", originalURL)
		c.Data(http.StatusTemporaryRedirect, contentType, nil)
		return
	}
	c.Data(http.StatusBadRequest, contentType, nil)
}

func (s HandlerService) JSONShortenURL(c *gin.Context) {
	var request dto.ShortenRequest
	rawRequest, err := c.GetRawData()

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			dto.ShortenResponse{Result: requestBodyGetErr},
		)
		return
	}

	if err := json.Unmarshal(rawRequest, &request); err != nil {

		c.JSON(
			http.StatusBadRequest,
			dto.ShortenResponse{Result: requestBindingErr},
		)
		return
	}

	shortURL, err := s.Service.ShortenURL(s.Cfg, string(request.URL))

	if errors.Is(err, shrterr.ErrOriginalURLAlreadyExists) {
		c.JSON(http.StatusConflict, dto.ShortenResponse{Result: shortURL})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error while shorten url",
		})
		return
	}

	c.JSON(http.StatusCreated, dto.ShortenResponse{Result: shortURL})

}

func (s HandlerService) Ping(c *gin.Context) {
	err := s.Storage.Ping()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

func (s HandlerService) BatchShortenURL(c *gin.Context) {
	var batchRequestData []dto.BatchShortenRequest
	if err := c.ShouldBindJSON(&batchRequestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "incorrect request body (error json binding)",
		})
		return
	}

	if len(batchRequestData) < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "incorrect request body (items in array less than 1)",
		})
		return
	}

	data, err := s.Service.ShortenURLBatch(s.Cfg, batchRequestData)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error while batch url shorten",
		})
		return
	}

	c.JSON(http.StatusCreated, data)
}
