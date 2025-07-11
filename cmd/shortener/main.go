// Package main initializes and starts the URL shortener web application,
// setting up configuration, logging, storage, services, and the HTTP server.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	srv := &http.Server{
		Addr:    *cfg.ServerAddress,
		Handler: r.Handler(),
	}

	go func() {
		logger.Info("preparing to start web server")
		if *cfg.ShouldUseTLS {
			logger.Info("starting web server with tls config", zap.Any("config", *cfg.TLSConfig))
			if err := srv.ListenAndServeTLS(cfg.TLSConfig.CrtFilePath, cfg.TLSConfig.KeyFilePath); err != nil && err != http.ErrServerClosed {
				logger.Fatal("error starting http web server with tls", zap.Error(err))
			}
		} else {
			logger.Info("starting http web server")
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.Fatal("error starting http web server", zap.Error(err))
			}
		}
	}()

	gracefuShutdownCh := make(chan os.Signal, 1)

	signal.Notify(gracefuShutdownCh, syscall.SIGINT, syscall.SIGTERM)

	<-gracefuShutdownCh

	logger.Info("received shutdown signal, gracefully shutting down web server")

	shutdownCtx, shutdownCtxCancel := context.WithTimeout(context.Background(), time.Second*10)
	defer shutdownCtxCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown error", zap.Error(err))
	}

	<-shutdownCtx.Done()

	if storage.GetType() == "database" {
		logger.Info("closing connection to the database")
		storage.(*database.Database).Close()
	}

	logger.Info("goodbye!")
}

func printStartupInfo(l *zap.Logger) {
	l.Info("Build version", zap.String("version", buildVersion))
	l.Info("Build date", zap.String("date", buildDate))
	l.Info("Build commit", zap.String("commit", buildCommit))
}
