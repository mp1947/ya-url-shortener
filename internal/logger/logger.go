package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

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
