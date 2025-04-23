package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mp1947/ya-url-shortener/internal/entity"
)

func (s HandlerService) DeleteUserURLs(c *gin.Context) {
	var userURLsToDelete []string

	userID, exists := c.Get("user_id")

	if !exists {
		c.Status(http.StatusUnauthorized)
		return
	}

	if err := c.BindJSON(&userURLsToDelete); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "incorrect request body",
		})
		return
	}

	if len(userURLsToDelete) < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "no urls provided for deletion",
		})
		return
	}

	s.Service.DeleteURLsBatch(
		c.Request.Context(),
		entity.BatchDeleteShortURLs{
			ShortURLs: userURLsToDelete,
			UserID:    fmt.Sprintf("%s", userID),
		},
	)

	c.JSON(http.StatusAccepted, gin.H{
		"message": "in progress",
	})

}
