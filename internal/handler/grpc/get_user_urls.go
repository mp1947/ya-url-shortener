package handlegrpc

import (
	"context"
	"strings"

	pb "github.com/mp1947/ya-url-shortener/internal/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// GetUserURLS retrieves all URLs associated with the current user.
// It extracts the user ID from the incoming gRPC metadata, fetches the user's URLs
// from the service layer, and returns them in the response.
// Returns an error if metadata is missing or if there is a problem retrieving URLs.
func (g *GRPCService) GetUserURLS(
	ctx context.Context,
	in *pb.Empty,
) (*pb.GetUserURLSResp, error) {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Internal, "internal error: metadata not found")
	}

	userID := md["user_id"][0]

	urls, err := g.Service.GetUserURLs(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error retrieving user URLs: %v", err)
	}

	response := &pb.GetUserURLSResp{
		UserURLs: make([]*pb.GetUserURLSResp_UserURL, 0, len(urls)),
	}

	for _, url := range urls {
		response.UserURLs = append(response.UserURLs, &pb.GetUserURLSResp_UserURL{
			ShortURL:    strings.Replace(url.ShortURL, *g.Cfg.BaseHTTPURL, *g.Cfg.BaseGRPCURL, 1),
			OriginalURL: url.OriginalURL,
		})
	}

	return response, nil
}
