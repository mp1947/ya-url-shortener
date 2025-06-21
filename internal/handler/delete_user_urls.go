package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mp1947/ya-url-shortener/internal/model"
)

// DeleteUserURLs handles the deletion of a batch of user-specific short URLs.
//
// @Summary      Delete user short URLs
// @Description  Deletes a batch of short URLs belonging to the authenticated user. The request body must be a JSON array of short URL identifiers.
// @Tags         urls
// @Accept       json
// @Produce      json
// @Param        userURLsToDelete  body      []string  true  "Array of short URL identifiers to delete"
// @Success      202  {object}  map[string]string  "in progress"
// @Failure      400  {object}  map[string]string  "incorrect request body or no urls provided for deletion"
// @Failure      401  {object}  nil                "unauthorized"
// @Router       /api/user/urls [delete]
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
		model.BatchDeleteShortURLs{
			ShortURLs: userURLsToDelete,
			UserID:    fmt.Sprintf("%s", userID),
		},
	)

	c.JSON(http.StatusAccepted, gin.H{
		"message": "in progress",
	})

}
