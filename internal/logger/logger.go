package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitLogger initializes and returns a new zap.Logger instance configured for production use.
// The logger uses RFC3339 time format and disables caller information in log entries.
// It logs a message upon successful initialization.
// Returns the configured *zap.Logger and an error if initialization fails.
func InitLogger() (*zap.Logger, error) {

	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	config.DisableCaller = true
	config.EncoderConfig.TimeKey = "time"

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	logger.Info("logger has been successfully initialized")

	return logger, nil
}
