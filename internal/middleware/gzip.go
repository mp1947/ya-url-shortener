package middleware

import (
	"io"

	"github.com/gin-gonic/gin"
)

type gzipWriter struct {
	gin.ResponseWriter
	writer     io.Writer
	statusCode int
}

func (gzw *gzipWriter) Write(p []byte) (int, error) {
	return gzw.writer.Write(p)
}

func (gzw *gzipWriter) WriteHeader(statusCode int) {
	gzw.statusCode = statusCode
	gzw.ResponseWriter.WriteHeader(gzw.statusCode)
}
