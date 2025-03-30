package main

import (
	"log"

	"github.com/mp1947/ya-url-shortener/config"
	"github.com/mp1947/ya-url-shortener/internal/eventlog"
	"github.com/mp1947/ya-url-shortener/internal/logger"
	"github.com/mp1947/ya-url-shortener/internal/repository"
	"github.com/mp1947/ya-url-shortener/internal/repository/database"
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

	logger.Info("parsing configuration parameters")
	cfg := config.Config{}
	cfg.InitConfig()

	logger.Info(
		"config has been initialized",
		zap.String("host", *cfg.ListenAddr),
		zap.String("base_url", *cfg.BaseURL),
	)

	logger.Info("creating and initializing storage")

	var storage repository.Repository

	if *cfg.DatabaseDSN != "" {
		storage = &database.Database{}
	} else {
		storage = &inmemory.Memory{}
	}

	if err := storage.Init(cfg); err != nil {
		log.Fatal("error initializing storage", zap.Error(err))
	}

	logger.Info(
		"storage has been initialized",
		zap.String("type", storage.GetType()),
	)

	ep, err := eventlog.NewEventProcessor(cfg)

	if err != nil {
		logger.Fatal("failed creating new event processor", zap.Error(err))
	}

	logger.Info(
		"restoring records from file storage",
		zap.String("file_storage_path", *cfg.FileStoragePath),
	)
	numRecordsRestored, err := ep.RestoreFromFile(cfg, storage, logger)

	if err != nil {
		logger.Fatal("error loading data from file", zap.Error(err))
	}

	logger.Info("records restored", zap.Int("count", numRecordsRestored))

	service := service.ShortenService{
		Storage: storage,
		EP:      *ep,
	}

	r := router.CreateRouter(cfg, &service, storage, logger)

	logger.Info(
		"router has been created. web server is ready to start",
	)

	if err := r.Run(*cfg.ListenAddr); err != nil {
		logger.Fatal("error starting web server", zap.Error(err))
	}
}
