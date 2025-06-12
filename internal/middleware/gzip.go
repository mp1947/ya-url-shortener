package middleware

import (
	"io"

	"github.com/gin-gonic/gin"
)

// gzipWriter is a custom wrapper around gin.ResponseWriter that enables gzip compression
// for HTTP responses. It embeds the original ResponseWriter, adds an io.Writer for
// compressed output, and tracks the HTTP status code.
type gzipWriter struct {
	gin.ResponseWriter
	writer     io.Writer
	statusCode int
}

// Write writes the provided byte slice to the underlying gzip writer.
// It returns the number of bytes written and any error encountered during the write operation.
func (gzw *gzipWriter) Write(p []byte) (int, error) {
	return gzw.writer.Write(p)
}

// WriteHeader sets the HTTP status code for the response and writes it to the underlying ResponseWriter.
// It also stores the status code in the gzipWriter for later reference.
func (gzw *gzipWriter) WriteHeader(statusCode int) {
	gzw.statusCode = statusCode
	gzw.ResponseWriter.WriteHeader(gzw.statusCode)
}
