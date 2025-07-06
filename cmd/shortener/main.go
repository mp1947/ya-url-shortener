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

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {

	cfg := config.InitConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, err := logger.InitLogger()

	if err != nil {
		log.Printf("error while initializing logger: %v\n", err)
	}

	printStartupInfo(logger)

	defer func() {
		if syncErr := logger.Sync(); err != nil {
			log.Fatalf("error while syncing logger: %v", syncErr)
		}
		logger.Info("logger has been synced")
	}()

	logger.Info(
		"initializing web application with config",
		zap.String("host", *cfg.ServerAddress),
		zap.String("base_url", *cfg.BaseURL),
	)

	storage, err := repository.CreateRepository(logger, *cfg, ctx)

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
		Cfg:     cfg,
		CommCh:  make(chan model.BatchDeleteShortURLs),
	}
	defer close(service.CommCh)

	go service.ProcessDeletions()

	r := router.CreateRouter(*cfg, &service, storage, logger)

	logger.Info(
		"router has been created. web server is ready to start",
	)

	if *cfg.ShouldUseTLS {
		logger.Info("starting web server with tls config", zap.Any("config", *cfg.TLSConfig))
		if err := r.RunTLS(
			*cfg.ServerAddress,
			cfg.TLSConfig.CrtFilePath,
			cfg.TLSConfig.KeyFilePath,
		); err != nil {
			logger.Fatal("error starting tls web server", zap.Error(err))
		}
	} else {
		if err := r.Run(*cfg.ServerAddress); err != nil {
			logger.Fatal("error starting simple web server", zap.Error(err))
		}
	}

}

func printStartupInfo(l *zap.Logger) {
	l.Info("Build version", zap.String("version", buildVersion))
	l.Info("Build date", zap.String("date", buildDate))
	l.Info("Build commit", zap.String("commit", buildCommit))
}
