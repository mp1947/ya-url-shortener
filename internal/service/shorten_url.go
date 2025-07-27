package service

import (
	"context"
	"errors"

	shrterr "github.com/mp1947/ya-url-shortener/internal/errors"
	"github.com/mp1947/ya-url-shortener/internal/usecase"
	"go.uber.org/zap"
)

// ShortenURL generates a shortened URL for the given original URL and associates it with the specified user ID.
// It logs the process of shortening, generates a unique short URL ID, and attempts to save the mapping in storage.
// If the original URL already exists, it returns the existing short URL and a specific error.
// On other errors, it returns an empty string and the error.
// On success, it returns the generated short URL and nil error.
//
// Parameters:
//   - ctx: context for request-scoped values, cancellation, and deadlines.
//   - url: the original URL to be shortened.
//   - userID: the identifier of the user requesting the shortening.
//
// Returns:
//   - string: the shortened URL.
//   - error: error if the operation failed, or a specific error if the URL already exists.
func (s *ShortenService) ShortenURL(
	ctx context.Context,
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
		return generateShortURL(*s.Cfg.BaseHTTPURL, shortURLID), err
	} else if err != nil {
		s.Logger.Warn("unexpected error", zap.Error(err))
		return "", err
	}

	return generateShortURL(*s.Cfg.BaseHTTPURL, shortURLID), nil
}
