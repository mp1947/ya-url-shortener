package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mp1947/ya-url-shortener/internal/dto"
)

// BatchShortenURL handles batch URL shortening requests.
//
// @Summary      Batch shorten URLs
// @Description  Accepts a batch of URLs and returns their shortened versions.
// @Tags         shortener
// @Accept       json
// @Produce      json
// @Param        body  body      []dto.BatchShortenRequest  true  "Batch shorten request"
// @Success      201   {array}   dto.BatchShortenResponse
// @Failure      400   {object}  map[string]string  "incorrect request body"
// @Failure      500   {object}  map[string]string  "error while batch url shorten"
// @Router       /api/shorten/batch [post]
// @Security     ApiKeyAuth
//
// BatchShortenURL expects a JSON array of BatchShortenRequest objects in the request body,
// validates the input, and returns a JSON array of shortened URLs. Returns HTTP 400 for invalid input
// and HTTP 500 for internal errors.
func (s HandlerService) BatchShortenURL(c *gin.Context) {
	var batchRequestData []dto.BatchShortenRequest
	userID, _ := c.Get("user_id")

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
