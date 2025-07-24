// Package main initializes and starts the URL shortener web application,
// setting up configuration, logging, storage, services, and the HTTP server.
package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mp1947/ya-url-shortener/config"
	handlegrpc "github.com/mp1947/ya-url-shortener/internal/handler/grpc"
	"github.com/mp1947/ya-url-shortener/internal/interceptor"
	"github.com/mp1947/ya-url-shortener/internal/logger"
	"github.com/mp1947/ya-url-shortener/internal/model"
	"github.com/mp1947/ya-url-shortener/internal/proto"
	"github.com/mp1947/ya-url-shortener/internal/repository"
	"github.com/mp1947/ya-url-shortener/internal/repository/database"
	"github.com/mp1947/ya-url-shortener/internal/router"
	"github.com/mp1947/ya-url-shortener/internal/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
		zap.String("host", *cfg.HTTPServerAddress),
		zap.String("base_url", *cfg.BaseHTTPURL),
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
		Addr:    *cfg.HTTPServerAddress,
		Handler: r.Handler(),
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptor.AuthUnaryInterceptor))

	go func() {
		logger.Info("preparing to start http web server")
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

	go func() {
		logger.Info("preparing to start grpc server")
		l, err := net.Listen("tcp", *cfg.GRPCServerAddress)
		if err != nil {
			logger.Fatal("error creating gcrp listener", zap.Error(err))
		}

		proto.RegisterShortenerServer(grpcServer, handlegrpc.NewGRPCService(&service, cfg))
		reflection.Register(grpcServer)

		logger.Info("starting grpc server on address", zap.String("address", *cfg.GRPCServerAddress))

		if err := grpcServer.Serve(l); err != nil {
			logger.Fatal("error starting grpc server", zap.Error(err))
		}

	}()

	gracefuShutdownCh := make(chan os.Signal, 1)

	signal.Notify(gracefuShutdownCh, syscall.SIGINT, syscall.SIGTERM)

	<-gracefuShutdownCh

	logger.Info("received shutdown signal, gracefully shutting down http server")

	shutdownCtx, shutdownCtxCancel := context.WithTimeout(context.Background(), time.Second*10)
	defer shutdownCtxCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("http server graceful shutdown error", zap.Error(err))
	}

	grpcServer.GracefulStop()
	logger.Info("grpc server has been stopped")

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
