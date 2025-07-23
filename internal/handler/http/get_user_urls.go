package handlehttp

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUserURLs handles the HTTP request to retrieve all URLs associated with the authenticated user.
//
// @Summary      Get user's URLs
// @Description  Returns a list of URLs that belong to the authenticated user.
// @Tags         urls
// @Produce      json
// @Success      200 {array} models.UserURL "List of user's URLs"
// @NoContent    204 "No URLs found for the user"
// @Failure      401 {object} gin.H "Unauthorized"
// @Failure      500 {object} gin.H "Internal server error"
// @Router       /api/user/urls [get]
// @Security     ApiKeyAuth
//
// The handler expects the user ID to be set in the context (typically by authentication middleware).
// If the user is not authenticated, it responds with HTTP 401 Unauthorized.
// If the user has no URLs, it responds with HTTP 204 No Content.
// On success, it returns HTTP 200 OK with a JSON array of URLs.
// On internal errors, it responds with HTTP 500 Internal Server Error.
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
