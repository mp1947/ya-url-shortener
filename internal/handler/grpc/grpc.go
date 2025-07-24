// Package handlegrpc provides GRPC Handler functions for routing and processing API requests,
// delegating business logic to the service layer.
package handlegrpc

import (
	"github.com/mp1947/ya-url-shortener/config"
	pb "github.com/mp1947/ya-url-shortener/internal/proto"
	"github.com/mp1947/ya-url-shortener/internal/service"
)

// GRPCService implements the gRPC server for the URL shortener service.
// It embeds the generated UnimplementedShortenerServer to satisfy the gRPC interface
// and holds a reference to the core service logic via the Service field.
type GRPCService struct {
	pb.UnimplementedShortenerServer
	Service service.Service
	Cfg     *config.Config
}

// NewGRPCService creates and returns a new GRPCService instance using the provided service.
// It initializes the GRPCService with the given service implementation.
func NewGRPCService(s service.Service, cfg *config.Config) *GRPCService {
	return &GRPCService{Service: s, Cfg: cfg}
}
