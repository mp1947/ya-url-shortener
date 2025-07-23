// Package handlegrpc provides GRPC Handler functions for routing and processing API requests,
// delegating business logic to the service layer.
package handlegrpc

import (
	pb "github.com/mp1947/ya-url-shortener/internal/proto"
	"github.com/mp1947/ya-url-shortener/internal/service"
)

type GRPCService struct {
	pb.UnimplementedShortenerServer
	Service service.Service
}

func NewGRPCService(s service.Service) *GRPCService {
	return &GRPCService{Service: s}
}
