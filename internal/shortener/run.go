package shortener

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

// Run starts the Shortener service by launching background processes for handling deletions,
// running the HTTP and optional gRPC servers, and waits for a termination signal (SIGINT or SIGTERM).
// Upon receiving a shutdown signal, it gracefully shuts down all running services within a 10-second timeout.
// Logs errors encountered during server execution or shutdown.
func (s *Shortener) Run() {

	go s.service.ProcessDeletions()

	go func() {
		if err := s.runHTTP(); err != nil {
			s.Logger.Error("error running HTTP server", zap.Error(err))
		}
	}()

	if s.grpcServer != nil {
		go func() {
			if err := s.runGRPC(); err != nil {
				s.Logger.Error("error running gRPC server", zap.Error(err))
			}
		}()
	}

	gracefuShutdownCh := make(chan os.Signal, 1)

	signal.Notify(gracefuShutdownCh, syscall.SIGINT, syscall.SIGTERM)

	<-gracefuShutdownCh

	shutdownCtx, shutdownCtxCancel := context.WithTimeout(context.Background(), time.Second*10)
	defer shutdownCtxCancel()

	if err := s.Shutdown(shutdownCtx); err != nil {
		s.Logger.Fatal("error during shutdown", zap.Error(err))
	}

	<-shutdownCtx.Done()

}
