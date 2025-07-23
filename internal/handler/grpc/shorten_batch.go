package handlegrpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/mp1947/ya-url-shortener/internal/dto"
	pb "github.com/mp1947/ya-url-shortener/internal/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (g *GRPCService) BatchShortenURL(
	ctx context.Context,
	in *pb.BatchShortenReq,
) (*pb.BatchShortenResp, error) {

	if len(in.BatchShortenData) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "no URLs provided to shorten")
	}

	batchData := make([]dto.BatchShortenRequest, 0, len(in.BatchShortenData))
	for _, inData := range in.BatchShortenData {
		batchData = append(batchData, dto.BatchShortenRequest{
			CorrelationID: inData.CorrelationID,
			OriginalURL:   inData.OriginalURL,
		})
	}

	dataShortened, err := g.Service.ShortenURLBatch(ctx, batchData, uuid.NewString())
	if err != nil {
		return nil, err
	}

	response := make([]*pb.BatchShortenResp_BatchShorten, 0, len(dataShortened))
	for _, d := range dataShortened {
		response = append(response, &pb.BatchShortenResp_BatchShorten{
			CorrelationID: d.CorrelationID,
			ShortURL:      d.ShortURL,
		})
	}

	return &pb.BatchShortenResp{
		BatchShortenData: response,
	}, nil
}
