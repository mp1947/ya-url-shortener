package shortener

import (
	"net/http"

	"go.uber.org/zap"
)

// runHTTP starts the HTTP web server for the Shortener service.
// It checks the configuration to determine whether to use TLS or not.
// If TLS is enabled, it starts the server with the provided TLS certificate and key files.
// Otherwise, it starts a standard HTTP server.
// Logs are generated for server startup and any errors encountered.
// Returns an error if the server fails to start, except when the error is http.ErrServerClosed.
func (s *Shortener) runHTTP() error {
	s.Logger.Info("preparing to start http web server")
	if *s.cfg.ShouldUseTLS {
		s.Logger.Info("starting web server with tls config", zap.Any("config", *s.cfg.TLSConfig))
		if err := s.httpServer.ListenAndServeTLS(
			s.cfg.TLSConfig.CrtFilePath,
			s.cfg.TLSConfig.KeyFilePath,
		); err != nil && err != http.ErrServerClosed {

			s.Logger.Fatal("error starting http web server with tls", zap.Error(err))
			return err

		}
	} else {
		s.Logger.Info("starting http web server")
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.Logger.Fatal("error starting http web server", zap.Error(err))
			return err
		}
	}
	return nil
}
