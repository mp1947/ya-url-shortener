package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
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
	err := s.Service.DeleteURLsBatch(
		c.Request.Context(),
		userURLsToDelete,
		fmt.Sprintf("%s", userID),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error processing delete request",
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "in progress",
	})

}
