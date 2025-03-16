package middleware

import (
	"io"

	"github.com/gin-gonic/gin"
)

type gzipWriter struct {
	gin.ResponseWriter
	writer io.Writer
}

func (gzw *gzipWriter) Write(p []byte) (int, error) {
	return gzw.writer.Write(p)
}

func (gzw *gzipWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		gzw.ResponseWriter.Header().Set("Content-Encoding", "gzip")
	}
	gzw.ResponseWriter.WriteHeader(statusCode)
}
