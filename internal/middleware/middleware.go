package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LoggerMiddleware(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestURI := c.Request.URL.RequestURI()
		requestMethod := c.Request.Method
		t := time.Now()

		c.Next()

		duration := time.Since(t)
		status := c.Writer.Status()
		bodySize := c.Writer.Size()

		log.Info(
			"request processed",
			zap.String("request_uri", requestURI),
			zap.String("request_method", requestMethod),
			zap.Any("request_duration", duration),
			zap.Int("response_status_code", status),
			zap.Int("response_body_size", bodySize),
		)
	}
}
