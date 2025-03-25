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

		supportsGzip := shouldUseGzip(c.GetHeader("Accept-Encoding"))

		if !supportsGzip {
			c.Next()
			return
		}

		gzw, err := gzip.NewWriterLevel(c.Writer, gzip.BestSpeed)

		if err != nil {
			io.WriteString(c.Writer, err.Error())
			return
		}

		defer gzw.Close()
		c.Writer = &gzipWriter{
			ResponseWriter: c.Writer,
			writer:         gzw,
		}
		c.Writer.Header().Set("Content-Encoding", "gzip")
		gzw.Flush()

		c.Next()
	}
}

func shouldUseGzip(acceptEncoding string) bool {
	if acceptEncoding == "" {
		return false
	}

	encodings := strings.Split(acceptEncoding, ",")
	for _, enc := range encodings {
		enc = strings.ToLower(strings.TrimSpace(enc))
		if strings.Contains(enc, "gzip") {
			if strings.Contains(enc, "q=") {
				parts := strings.Split(enc, "q=")
				if len(parts) > 1 {
					qValue := strings.TrimSpace(parts[1])
					if qValue == "0.0" || qValue == "0" {
						return false
					}
				}
			}
			return true
		}

	}
	return false
}
