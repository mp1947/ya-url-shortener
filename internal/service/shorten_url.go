package service

import (
	"context"
	"errors"

	"github.com/mp1947/ya-url-shortener/config"
	shrterr "github.com/mp1947/ya-url-shortener/internal/errors"
	"github.com/mp1947/ya-url-shortener/internal/usecase"
	"go.uber.org/zap"
)

func (s *ShortenService) ShortenURL(
	ctx context.Context,
	cfg config.Config,
	url string,
	userID string,
) (string, error) {
	s.Logger.Info("shortening incoming url", zap.String("original_url", url))

	shortURLID := usecase.GenerateIDFromURL(url)

	s.Logger.Info(
		"short_url id generated for url",
		zap.String("short_url_id", shortURLID),
		zap.String("original_url", url),
		zap.String("user_id", userID),
	)

	err := s.Storage.Save(ctx, shortURLID, url, userID)
	if errors.Is(err, shrterr.ErrOriginalURLAlreadyExists) {
		s.Logger.Info(
			"original_url already exists, returning error with short url",
			zap.Error(err),
			zap.String("original_url", url),
		)
		return generateShortURL(*cfg.BaseURL, shortURLID), err
	} else if err != nil {
		s.Logger.Warn("unexpected error", zap.Error(err))
		return "", err
	}

	return generateShortURL(*cfg.BaseURL, shortURLID), nil
}
