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

// LoggerMiddleware returns a Gin middleware handler that logs details about each HTTP request.
// It logs the request URI, HTTP method, processing duration, response status code, and response body size
// using the provided zap.Logger instance. The middleware should be attached to a Gin router to enable
// structured logging of incoming requests and their corresponding responses.
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

// GzipMiddleware is a Gin middleware that transparently handles gzip compression and decompression for HTTP requests and responses.
//
// For incoming requests, if the "Content-Encoding" header contains "gzip", the middleware decompresses the request body before passing it to the next handler.
// For outgoing responses, if the "Accept-Encoding" header indicates support for gzip, the middleware compresses the response body using gzip and sets the "Content-Encoding: gzip" header.
// If the client does not support gzip, the response is sent uncompressed.
//
// This middleware ensures efficient bandwidth usage for clients that support gzip, while maintaining compatibility with clients that do not.
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

		defer func() {
			gzw.Flush()
			gzw.Close()
		}()

		c.Writer = &gzipWriter{
			ResponseWriter: c.Writer,
			writer:         gzw,
		}
		c.Writer.Header().Set("Content-Encoding", "gzip")

		c.Next()
	}
}

// shouldUseGzip determines whether gzip compression should be used based on the
// provided Accept-Encoding header value. It returns true if "gzip" is present
// in the header and its quality value (q) is not set to 0 or 0.0, indicating
// that the client accepts gzip encoding. Returns false if "gzip" is absent or
// explicitly declined by the client.
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
