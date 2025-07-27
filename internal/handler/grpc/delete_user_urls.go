package handlegrpc

import (
	"context"

	"github.com/mp1947/ya-url-shortener/internal/model"
	pb "github.com/mp1947/ya-url-shortener/internal/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DeleteUserURLS handles a gRPC request to delete a batch of user URLs.
// It extracts the user ID from the incoming context metadata, then initiates
// a batch deletion of the provided short URLs associated with the user.
// The operation is performed asynchronously, and the response status is set to "pending".
// Returns an error if the user metadata is not found in the context.
func (g *GRPCService) DeleteUserURLS(
	ctx context.Context,
	in *pb.DeleteURLSReq,
) (*pb.DeleteURLSResp, error) {

	userID, _, err := g.getDataFromMD(ctx)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if len(in.ShortURLs) == 0 {
		return nil, status.Error(codes.InvalidArgument, "no short URLs provided")
	}

	var result pb.DeleteURLSResp

	g.Service.DeleteURLsBatch(
		ctx,
		model.BatchDeleteShortURLs{
			ShortURLs: in.ShortURLs,
			UserID:    userID,
		},
	)

	result.Status = "pending"

	return &result, nil
}
