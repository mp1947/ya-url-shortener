package handlegrpc

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/mp1947/ya-url-shortener/internal/proto"
)

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
