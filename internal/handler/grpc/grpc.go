// Package handlegrpc provides GRPC Handler functions for routing and processing API requests,
// delegating business logic to the service layer.
package handlegrpc

import (
	"context"
	"errors"

	"github.com/mp1947/ya-url-shortener/config"
	pb "github.com/mp1947/ya-url-shortener/internal/proto"
	"github.com/mp1947/ya-url-shortener/internal/service"
	"google.golang.org/grpc/metadata"
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

// getDataFromMD extracts the "user_id" and "token" values from the gRPC metadata in the provided context.
// It returns the user ID, token, and an error if the metadata is not found or the required keys are missing.
func (g *GRPCService) getDataFromMD(ctx context.Context) (userID string, token string, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		err = errors.New("metadata not found")
		return "", "", err
	}
	userID = md["user_id"][0]
	token = md["token"][0]
	return userID, token, nil
}
