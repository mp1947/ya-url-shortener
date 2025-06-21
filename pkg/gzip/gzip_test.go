package gzip_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mp1947/ya-url-shortener/pkg/gzip"
	"github.com/stretchr/testify/assert"
)

func TestWrite(t *testing.T) {
	var buf bytes.Buffer
	gzw := gzip.GzipWriter{
		Writer: &buf,
	}
	data := []byte("test payload")
	n, err := gzw.Write(data)
	assert.NoError(t, err)
	assert.Equal(t, len(data), n)
	assert.Equal(t, data, buf.Bytes())
}

func TestWriteHeader(t *testing.T) {
	recorder := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(recorder)

	gzw := &gzip.GzipWriter{
		ResponseWriter: c.Writer,
	}

	gzw.WriteHeader(http.StatusNoContent)

	assert.Equal(t, http.StatusNoContent, gzw.Status())

}
