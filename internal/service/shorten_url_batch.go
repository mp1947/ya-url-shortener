package service

import (
	"context"
	"fmt"

	"github.com/mp1947/ya-url-shortener/config"
	"github.com/mp1947/ya-url-shortener/internal/dto"
	"github.com/mp1947/ya-url-shortener/internal/entity"
	"github.com/mp1947/ya-url-shortener/internal/usecase"
	"go.uber.org/zap"
)

func (s *ShortenService) ShortenURLBatch(
	ctx context.Context,
	cfg config.Config,
	batchData []dto.BatchShortenRequest,
	userID string,
) ([]dto.BatchShortenResponse, error) {
	s.Logger.Info(
		"processing batch of urls",
		zap.Any("batch_data", batchData),
	)

	urls := make([]entity.URLWithCorrelation, len(batchData))
	result := make([]dto.BatchShortenResponse, len(batchData))

	for i, v := range batchData {
		shortURLID := usecase.GenerateIDFromURL(v.OriginalURL)
		urls[i] = entity.URLWithCorrelation{
			ShortURLID:    shortURLID,
			OriginalURL:   v.OriginalURL,
			CorrelationID: v.CorrelationID,
		}
		result[i] = dto.BatchShortenResponse{
			CorrelationID: v.CorrelationID,
			ShortURL:      fmt.Sprintf("%s/%s", *cfg.BaseURL, shortURLID),
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
