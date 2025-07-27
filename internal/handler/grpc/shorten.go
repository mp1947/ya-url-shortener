package handlegrpc

import (
	"context"
	"strings"

	pb "github.com/mp1947/ya-url-shortener/internal/proto"
	"google.golang.org/grpc/codes"
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

	userID, token, err := g.getDataFromMD(ctx)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var response pb.ShortenURLResp

	shortURL, err := g.Service.ShortenURL(
		ctx,
		string(in.Url),
		userID,
	)

	if err != nil {
		return nil, status.Errorf(codes.AlreadyExists, "error while shortening URL: %v", err)
	}

	response.ShortURL = strings.Replace(shortURL, *g.Cfg.BaseHTTPURL, *g.Cfg.BaseGRPCURL, 1)
	response.JwtToken = token

	return &response, nil
}
