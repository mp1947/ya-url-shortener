package handlegrpc

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/mp1947/ya-url-shortener/internal/proto"
)

// ShortenURL handles the gRPC request to shorten a given URL.
// It receives a ShortenURLReq containing the original URL, generates a unique identifier,
// and calls the underlying service to create a shortened URL.
// Returns a ShortenURLResp with the shortened URL or an error if the operation fails.
func (g *GRPCService) ShortenURL(
	ctx context.Context,
	in *pb.ShortenURLReq,
) (*pb.ShortenURLResp, error) {
	var response pb.ShortenURLResp

	shortURL, err := g.Service.ShortenURL(
		ctx,
		string(in.Url),
		uuid.NewString(),
	)

	if err != nil {
		return nil, err
	}

	response.ShortURL = shortURL

	return &response, nil
}
