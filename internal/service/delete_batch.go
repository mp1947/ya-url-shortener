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
	s.CommCh <- shortURLs
	return nil
}

func (s *ShortenService) ProcessDeletions() {
	s.Logger.Info("starting deletions processing goroutine")
	for data := range s.CommCh {
		s.Logger.Info("received new data for deletion", zap.Any("data", data))
		if err := s.Storage.DeleteBatch(context.TODO(), data); err != nil {
			s.Logger.Warn("error batch-deleting short urls", zap.Error(err))
		}
	}
}
