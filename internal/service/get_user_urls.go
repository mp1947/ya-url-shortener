package service

import (
	"context"

	"github.com/mp1947/ya-url-shortener/config"
	"github.com/mp1947/ya-url-shortener/internal/dto"
	"go.uber.org/zap"
)

func (s *ShortenService) GetUserURLs(
	ctx context.Context,
	cfg config.Config,
	userID string,
) ([]dto.ShortenURLsByUserID, error) {
	s.Logger.Info(
		"processing shorten urls for user",
		zap.String("user_id", userID),
	)

	userURLs, err := s.Storage.GetURLsByUserID(ctx, userID)

	if err != nil {
		s.Logger.Warn("error getting urls by user id", zap.Error(err))
		return nil, err
	}

	userURLsResponse := make([]dto.ShortenURLsByUserID, len(userURLs))

	for i, v := range userURLs {
		userURLsResponse[i] = dto.ShortenURLsByUserID{
			ShortURL:    generateShortURL(*cfg.BaseURL, v.ShortURLID),
			OriginalURL: v.OriginalURL,
		}
	}
	return userURLsResponse, nil
}
