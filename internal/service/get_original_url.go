package service

import (
	"context"

	"github.com/mp1947/ya-url-shortener/internal/model"
	"go.uber.org/zap"
)

// GetOriginalURL retrieves the original URL associated with the given short URL ID.
// It logs the process of fetching the URL, including any errors encountered and the result.
// Returns the corresponding model.URL and an error if the retrieval fails.
func (s *ShortenService) GetOriginalURL(
	ctx context.Context,
	shortURLID string,
) (model.URL, error) {

	s.Logger.Info("processing short url with id", zap.String("short_url_id", shortURLID))
	data, err := s.Storage.Get(ctx, shortURLID)
	if err != nil {
		s.Logger.Warn(
			"error getting original_url by short_url_id",
			zap.String("short_url", shortURLID),
			zap.Error(err),
		)
		return model.URL{}, err
	}
	s.Logger.Info(
		"retrieved original_url by short_url_id",
		zap.String("short_url_id", shortURLID),
		zap.String("original_url", data.OriginalURL),
		zap.Bool("is_deleted", data.IsDeleted),
	)
	return data, nil
}
