package handlegrpc

import (
	"context"
	"errors"
	"strings"

	pb "github.com/mp1947/ya-url-shortener/internal/proto"
)

// GetOriginalURLByShort retrieves the original URL corresponding to a given shortened URL.
// It parses the short URL to extract the unique identifier, then queries the service layer
// to obtain the original URL. Returns an error if the identifier is invalid or if the service
// fails to find the original URL.
//
// Parameters:
//   - ctx: The context for the request, used for cancellation and deadlines.
//   - in: The request containing the shortened URL.
//
// Returns:
//   - *pb.GetOriginalURLByShortResp: The response containing the original URL.
//   - error: An error if the operation fails.
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
