package logger_test

import (
	"testing"

	"github.com/mp1947/ya-url-shortener/internal/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestInitLogger(t *testing.T) {
	t.Run("init logger", func(t *testing.T) {
		l, err := logger.InitLogger()
		assert.NoError(t, err)
		assert.IsType(t, &zap.Logger{}, l)
	})
}
