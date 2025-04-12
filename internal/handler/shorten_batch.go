package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mp1947/ya-url-shortener/internal/dto"
)

func (s HandlerService) BatchShortenURL(c *gin.Context) {
	var batchRequestData []dto.BatchShortenRequest
	userID, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "unexpected internal server error (uuid not exists)",
		})
		return
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
