// Package main initializes and starts the URL shortener web application,
// setting up configuration, logging, storage, services, and the HTTP server.
package main

import (
	"context"
	"log"

	"github.com/mp1947/ya-url-shortener/config"
	"github.com/mp1947/ya-url-shortener/internal/logger"
	"github.com/mp1947/ya-url-shortener/internal/model"
	"github.com/mp1947/ya-url-shortener/internal/repository"
	"github.com/mp1947/ya-url-shortener/internal/repository/database"
	"github.com/mp1947/ya-url-shortener/internal/router"
	"github.com/mp1947/ya-url-shortener/internal/service"
	"go.uber.org/zap"
)

func main() {

	cfg := config.Config{}
	cfg.InitConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, err := logger.InitLogger()

	if err != nil {
		log.Printf("error while initializing logger: %v\n", err)
	}

	defer func() {
		if syncErr := logger.Sync(); err != nil {
			log.Fatalf("error while syncing logger: %v", syncErr)
		}
		logger.Info("logger has been synced")
	}()

	logger.Info(
		"initializing web application with config",
		zap.String("host", *cfg.ListenAddr),
		zap.String("base_url", *cfg.BaseURL),
	)

	storage, err := repository.CreateRepository(logger, cfg, ctx)

	if err != nil {
		logger.Fatal("error creating repository", zap.Error(err))
	}

	logger.Info(
		"storage has been initialized",
		zap.String("type", storage.GetType()),
	)

	if storage.GetType() == "database" {
		defer storage.(*database.Database).Close()
	}

	service := service.ShortenService{
		Storage: storage,
		Logger:  logger,
		Cfg:     &cfg,
		CommCh:  make(chan model.BatchDeleteShortURLs),
	}
	defer close(service.CommCh)

	go service.ProcessDeletions()

	r := router.CreateRouter(cfg, &service, storage, logger)

	logger.Info(
		"router has been created. web server is ready to start",
	)

	if err := r.Run(*cfg.ListenAddr); err != nil {
		logger.Fatal("error starting web server", zap.Error(err))
	}
}
