package main

import (
	"log"

	"github.com/mp1947/ya-url-shortener/config"
	"github.com/mp1947/ya-url-shortener/internal/logger"
	"github.com/mp1947/ya-url-shortener/internal/repository/inmemory"
	"github.com/mp1947/ya-url-shortener/internal/router"
	"github.com/mp1947/ya-url-shortener/internal/service"
	"go.uber.org/zap"
)

func main() {

	logger, err := logger.InitLogger()

	if err != nil {
		log.Fatalf("error while initializing logger: %v", err)
	}

	defer logger.Sync() //nolint:errcheck

	logger.Info("initializing web server")

	logger.Info("creating and initializing storage")

	storage := &inmemory.Memory{}
	storage.Init()

	logger.Info(
		"storage has been initialized",
		zap.String("type", storage.StorageType),
	)

	logger.Info("creating shortener service")

	service := service.ShortenService{Storage: storage}

	logger.Info("parsing configuration parameters")
	cfg := config.Config{}
	cfg.ParseFlags()

	logger.Info(
		"config has been initialized",
		zap.String("host", *cfg.ListenAddr),
		zap.String("base_url", *cfg.BaseURL),
	)

	r := router.CreateRouter(cfg, service, logger)

	logger.Info(
		"router has been created. web server is ready to start",
	)

	if err := r.Run(*cfg.ListenAddr); err != nil {
		logger.Fatal("error starting web server", zap.Error(err))
	}
}
