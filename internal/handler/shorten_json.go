package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mp1947/ya-url-shortener/internal/dto"
	shrterr "github.com/mp1947/ya-url-shortener/internal/errors"
)

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
