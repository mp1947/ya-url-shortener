package handlegrpc

import (
	"context"
	"errors"
	"strings"

	pb "github.com/mp1947/ya-url-shortener/internal/proto"
)

func (g *GRPCService) GetOriginalURLByShort(
	ctx context.Context,
	in *pb.GetOriginalURLByShortReq,
) (*pb.GetOriginalURLByShortResp, error) {
	var response pb.GetOriginalURLByShortResp
	parsedURL := strings.Split(in.ShortURL, "/")
	var id string
	if len(parsedURL) >= 1 {
		id = parsedURL[len(parsedURL)-1]
	} else {
		return nil, errors.New("id not correct")
	}
	url, err := g.Service.GetOriginalURL(ctx, id)
	if err != nil {
		return nil, err
	}
	response.OriginalURL = url.OriginalURL
	return &response, nil
}
