package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s HandlerService) GetUserURLs(c *gin.Context) {
	userID, exists := c.Get("user_id")

	if !exists {
		c.Status(http.StatusUnauthorized)
		return
	}

	userIDStr := userID.(string)

	resp, err := s.Service.GetUserURLs(
		c.Request.Context(),
		userIDStr,
	)

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
