package service

import (
	"context"
	"fmt"

	"github.com/mp1947/ya-url-shortener/internal/dto"
	"github.com/mp1947/ya-url-shortener/internal/model"
	"github.com/mp1947/ya-url-shortener/internal/usecase"
	"go.uber.org/zap"
)

// ShortenURLBatch processes a batch of URL shortening requests for a specific user.
// It generates short URLs for each original URL in the batch, associates them with the provided user ID,
// and saves the batch to the storage. The function returns a slice of BatchShortenResponse containing
// the correlation IDs and the corresponding shortened URLs. If an error occurs during storage, it returns the error.
//
// Parameters:
//   - ctx: context.Context for request-scoped values, cancellation, and deadlines.
//   - batchData: slice of BatchShortenRequest containing the original URLs and correlation IDs.
//   - userID: string representing the user for whom the URLs are being shortened.
//
// Returns:
//   - []dto.BatchShortenResponse: slice containing the correlation IDs and shortened URLs.
//   - error: error encountered during the batch save operation, if any.
func (s *ShortenService) ShortenURLBatch(
	ctx context.Context,
	batchData []dto.BatchShortenRequest,
	userID string,
) ([]dto.BatchShortenResponse, error) {
	s.Logger.Info(
		"processing batch of urls",
		zap.Any("batch_data", batchData),
	)

	urls := make([]model.URLWithCorrelation, len(batchData))
	result := make([]dto.BatchShortenResponse, len(batchData))

	for i, v := range batchData {
		shortURLID := usecase.GenerateIDFromURL(v.OriginalURL)
		urls[i] = model.URLWithCorrelation{
			ShortURLID:    shortURLID,
			OriginalURL:   v.OriginalURL,
			CorrelationID: v.CorrelationID,
		}
		result[i] = dto.BatchShortenResponse{
			CorrelationID: v.CorrelationID,
			ShortURL:      fmt.Sprintf("%s/%s", *s.Cfg.BaseURL, shortURLID),
		}
	}

	_, err := s.Storage.SaveBatch(ctx, urls, userID)
	if err != nil {
		s.Logger.Warn("error while saving batch of urls", zap.Error(err))
		return nil, err
	}

	s.Logger.Info("batch of urls were successfully processed")

	return result, nil

}
