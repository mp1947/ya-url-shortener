package handlegrpc

import (
	"context"

	pb "github.com/mp1947/ya-url-shortener/internal/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// ShortenURL handles the gRPC request to shorten a given URL.
// It receives a ShortenURLReq containing the original URL, generates a unique identifier,
// and calls the underlying service to create a shortened URL.
// Returns a ShortenURLResp with the shortened URL or an error if the operation fails.
func (g *GRPCService) ShortenURL(
	ctx context.Context,
	in *pb.ShortenURLReq,
) (*pb.ShortenURLResp, error) {

	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return nil, status.Errorf(codes.Internal, "internal error: metadata not found")
	}

	var response pb.ShortenURLResp

	shortURL, err := g.Service.ShortenURL(
		ctx,
		string(in.Url),
		md["user_id"][0],
	)

	if err != nil {
		return nil, status.Errorf(codes.AlreadyExists, "error while shortening URL: %v", err)
	}

	response.ShortURL = shortURL
	response.JwtToken = md["token"][0]

	return &response, nil
}
