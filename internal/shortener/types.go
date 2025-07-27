// Package shortener provides the core logic for initializing and running the URL shortener service.
// It defines the Shortener struct, which encapsulates the repository, service layer, configuration, and logger.
// The package includes functions to initialize the Shortener with all necessary dependencies and to start the service.
package shortener

import (
	"net/http"

	"github.com/mp1947/ya-url-shortener/config"
	"github.com/mp1947/ya-url-shortener/internal/repository"
	"github.com/mp1947/ya-url-shortener/internal/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Shortener encapsulates the core components required for running the URL shortener service,
// including HTTP and gRPC servers, configuration, logging, repository, and business logic service.
type Shortener struct {
	httpServer *http.Server
	grpcServer *grpc.Server
	cfg        *config.Config
	Logger     *zap.Logger
	repo       repository.Repository
	service    service.ShortenService
}
