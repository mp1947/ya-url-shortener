package service

import (
	"context"

	"github.com/mp1947/ya-url-shortener/internal/entity"
	"go.uber.org/zap"
)

func (s *ShortenService) GetOriginalURL(
	ctx context.Context,
	shortURLID string,
) (entity.URL, error) {

	s.Logger.Info("processing short url with id", zap.String("short_url_id", shortURLID))
	data, err := s.Storage.Get(ctx, shortURLID)
	if err != nil {
		s.Logger.Warn("error getting original_url by short_url_id", zap.String("short_url", shortURLID))
		return entity.URL{}, err
	}
	s.Logger.Info(
		"retrieved original_url by short_url_id",
		zap.String("short_url_id", shortURLID),
		zap.String("original_url", data.OriginalURL),
		zap.Bool("is_deleted", data.IsDeleted),
	)
	return data, nil
}
