package service

import (
	"context"

	"github.com/mp1947/ya-url-shortener/internal/entity"
	"go.uber.org/zap"
)

func (s *ShortenService) DeleteURLsBatch(
	ctx context.Context,
	shortURLs entity.BatchDeleteShortURLs,
) error {
	s.Logger.Info(
		"putting short urls to delete into channel",
		zap.Any("data", shortURLs),
	)
	s.CommCh <- shortURLs
	return nil
}

func (s *ShortenService) ProcessDeletions() {
	s.Logger.Info("starting deletions processing goroutine")
	for data := range s.CommCh {
		ctx, cancel := context.WithCancel(context.Background())
		s.Logger.Info("received new data for deletion", zap.Any("data", data))
		rowsDeleted, err := s.Storage.DeleteBatch(ctx, data)
		if err != nil {
			s.Logger.Warn("error batch-deleting short urls", zap.Error(err))
			cancel()
		}
		s.Logger.Info(
			"data has been deleted from the database",
			zap.Any("data", data.ShortURLs),
			zap.String("user_id", data.UserID),
			zap.Int64("rows_deleted", rowsDeleted),
		)
		cancel()
	}
}
