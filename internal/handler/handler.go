package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mp1947/ya-url-shortener/config"
	"github.com/mp1947/ya-url-shortener/internal/dto"
	shrterr "github.com/mp1947/ya-url-shortener/internal/errors"
	"github.com/mp1947/ya-url-shortener/internal/repository/database"
	"github.com/mp1947/ya-url-shortener/internal/service"
)

const (
	contentTypePlain  = "text/plain; charset=utf-8"
	contentTypeJSON   = "application/json; charset=utf-8"
	requestBindingErr = "invalid request: error parsing request params"
	requestBodyGetErr = "error getting request body"
)

type HandlerService struct {
	Service service.Service
	Cfg     config.Config
}

func (s HandlerService) ShortenURL(c *gin.Context) {

	userID, exists := c.Get("user_id")

	if !exists {
		userID = uuid.New().String()
	}

	if c.Request.Method != http.MethodPost {
		c.Data(http.StatusBadRequest, contentTypePlain, nil)
		return
	}

	body, err := io.ReadAll(c.Request.Body)

	if err != nil || string(body) == "" {
		c.Data(http.StatusBadRequest, contentTypePlain, nil)
		return
	}

	shortURL, err := s.Service.ShortenURL(
		c.Request.Context(),
		s.Cfg,
		string(body),
		fmt.Sprintf("%s", userID),
	)

	if errors.Is(err, shrterr.ErrOriginalURLAlreadyExists) {
		c.Data(http.StatusConflict, contentTypePlain, []byte(shortURL))
		return
	} else if err != nil {
		c.Data(
			http.StatusInternalServerError,
			contentTypePlain,
			[]byte("internal server error while shorten url"),
		)
		return
	}

	c.Data(http.StatusCreated, contentTypePlain, []byte(shortURL))

}

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

func (s HandlerService) JSONShortenURL(c *gin.Context) {
	var request dto.ShortenRequest
	rawRequest, err := c.GetRawData()

	userID, exists := c.Get("user_id")

	if !exists {
		userID = uuid.New().String()
	}

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

	shortURL, err := s.Service.ShortenURL(
		c.Request.Context(),
		s.Cfg,
		string(request.URL),
		fmt.Sprintf("%s", userID),
	)

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

func (s HandlerService) Ping(db *database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := db.Ping(c.Request.Context())
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Status(http.StatusOK)
	}
}

func (s HandlerService) BatchShortenURL(c *gin.Context) {
	var batchRequestData []dto.BatchShortenRequest
	userID, exists := c.Get("user_id")

	if !exists {
		userID = uuid.New().String()
	}
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

	data, err := s.Service.ShortenURLBatch(
		c.Request.Context(),
		s.Cfg,
		batchRequestData,
		fmt.Sprintf("%s", userID),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error while batch url shorten",
		})
		return
	}

	c.JSON(http.StatusCreated, data)
}

func (s HandlerService) GetUserURLS(c *gin.Context) {

	userID, exists := c.Get("user_id")

	if !exists {
		c.Status(http.StatusUnauthorized)
		return
	}

	resp, err := s.Service.GetUserURLs(c.Request.Context(), s.Cfg, fmt.Sprintf("%s", userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error while processing urls",
		})
		return
	}

	if len(resp) < 1 {
		c.JSON(http.StatusNoContent, nil)
		return
	}
	c.JSON(http.StatusOK, resp)
}
