package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// InternalStats handles the request to retrieve internal statistics of the service.
// It calls the Service's GetInternalStats method and returns the result as a JSON response.
// On success, it responds with HTTP 200 and the statistics data.
// On failure, it responds with HTTP 500 and an error message.
func (s HandlerService) InternalStats(c *gin.Context) {

	resp, err := s.Service.GetInternalStats(c.Request.Context())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
