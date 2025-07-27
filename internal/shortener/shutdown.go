package shortener

import (
	"context"
	"errors"
	"log"
	"syscall"

	"github.com/mp1947/ya-url-shortener/internal/repository/database"
	"go.uber.org/zap"
)

// Shutdown gracefully shuts down the Shortener service, including the HTTP and gRPC servers,
// and closes any open database connections. It logs the shutdown process, ensures the logger
// is properly synced, and closes the communication channel. The method accepts a context
// for controlling the shutdown timeout and returns an error if any part of the shutdown fails.
func (s *Shortener) Shutdown(ctx context.Context) error {

	s.Logger.Info("received shutdown signal, gracefully shutting down shortener...")

	defer close(s.service.CommCh)

	defer func() {
		if err := s.Logger.Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
			log.Printf("error while syncing logger: %v", err)
		}
	}()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.Logger.Error("http server graceful shutdown error", zap.Error(err))
		return err
	}

	if s.grpcServer != nil {
		s.Logger.Info("gracefully shutting down gRPC server")
		s.grpcServer.GracefulStop()
	}

	if s.repo.GetType() == "database" {
		s.Logger.Info("closing connections to the database")
		s.service.Storage.(*database.Database).Close()
	}

	s.Logger.Info("goodbye!")

	return nil
}
