package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mp1947/ya-url-shortener/internal/repository/database"
)

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
