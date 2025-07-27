package service

import (
	"context"

	"github.com/mp1947/ya-url-shortener/internal/dto"
	"go.uber.org/zap"
)

// GetInternalStats retrieves internal statistics from the storage layer.
// It returns a dto.InternalStatsResp containing the statistics, or an error if the operation fails.
// The method logs any errors encountered during the retrieval process.
//
// Parameters:
//   - ctx: context.Context for request-scoped values and cancellation.
//
// Returns:
//   - *dto.InternalStatsResp: Pointer to the response containing internal statistics.
//   - error: Error encountered during retrieval, if any.
func (s *ShortenService) GetInternalStats(
	ctx context.Context,
) (*dto.InternalStatsResp, error) {

	resp, err := s.Storage.GetInternalStats(ctx)

	if err != nil {
		s.Logger.Error("received error from storage", zap.Error(err))
		return nil, err
	}

	return resp, nil
}
