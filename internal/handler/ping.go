package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mp1947/ya-url-shortener/internal/repository/database"
)

// Ping is a Gin handler that checks the health of the database connection.
//
// It attempts to ping the provided database instance using the request context.
// If the database is reachable, it responds with HTTP 200 OK.
// If the database is not reachable, it responds with HTTP 500 Internal Server Error.
//
// Swagger specification:
// @Summary      Health check
// @Description  Checks the health of the database connection.
// @Tags         health
// @Produce      plain
// @Success      200  {string}  string  "OK"
// @Failure      500  {string}  string  "Internal Server Error"
// @Router       /ping [get]
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
