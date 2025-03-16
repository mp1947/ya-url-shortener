package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
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

func GzipMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		contentEncoding := c.GetHeader("Content-Encoding")
		isRequestEncoded := strings.Contains(contentEncoding, "gzip")

		if isRequestEncoded {
			reader, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error": "invalid gzip data",
				})
				return
			}
			defer reader.Close()
			c.Request.Body = io.NopCloser(reader)
		}

		acceptEncoding := c.GetHeader("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if !supportsGzip {
			c.Next()
			return
		}

		gzw, _ := gzip.NewWriterLevel(c.Writer, gzip.BestSpeed)
		defer gzw.Close()
		c.Writer = &gzipWriter{
			ResponseWriter: c.Writer,
			writer:         gzw,
		}

		c.Next()
	}
}
