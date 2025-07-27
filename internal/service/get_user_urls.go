package service

import (
	"context"

	"github.com/mp1947/ya-url-shortener/internal/dto"
	"go.uber.org/zap"
)

// GetUserURLs retrieves all shortened URLs associated with a specific user ID.
// It logs the processing event, fetches the URLs from the storage layer, and constructs
// a response containing both the short and original URLs for each entry.
// Returns a slice of ShortenURLsByUserID DTOs or an error if retrieval fails.
//
// Parameters:
//   - ctx: context.Context for request-scoped values and cancellation.
//   - userID: string representing the unique identifier of the user.
//
// Returns:
//   - []dto.ShortenURLsByUserID: slice containing the user's shortened URLs.
//   - error: error encountered during retrieval, or nil if successful.
func (s *ShortenService) GetUserURLs(
	ctx context.Context,
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
			ShortURL:    generateShortURL(*s.Cfg.BaseHTTPURL, v.ShortURLID),
			OriginalURL: v.OriginalURL,
		}
	}
	return userURLsResponse, nil
}
