package shortener

import (
	"net"

	handlegrpc "github.com/mp1947/ya-url-shortener/internal/handler/grpc"
	"github.com/mp1947/ya-url-shortener/internal/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/reflection"
)

// runGRPC starts the gRPC server for the Shortener service.
// It sets up the TCP listener on the configured address, registers the gRPC service and reflection,
// and begins serving incoming gRPC requests. If any error occurs during setup or serving,
// it logs the error and returns it.
func (s *Shortener) runGRPC() error {
	s.Logger.Info("preparing to start grpc server")
	l, err := net.Listen("tcp", *s.cfg.GRPCServerAddress)
	if err != nil {
		s.Logger.Fatal("error creating gcrp listener", zap.Error(err))
		return err
	}

	proto.RegisterShortenerServer(s.grpcServer, handlegrpc.NewGRPCService(&s.service, s.cfg))
	reflection.Register(s.grpcServer)

	s.Logger.Info("starting grpc server on address", zap.String("address", *s.cfg.GRPCServerAddress))

	if err := s.grpcServer.Serve(l); err != nil {
		s.Logger.Fatal("error starting grpc server", zap.Error(err))
		return err
	}
	return nil
}
