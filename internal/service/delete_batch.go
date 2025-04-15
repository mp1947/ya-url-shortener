package service

import "context"

func (s *ShortenService) DeleteURLsBatch(
	ctx context.Context,
	shortURLs []string,
	userID string,
) error {
	return s.Storage.DeleteBatch(ctx, shortURLs, userID)
}
