package service

import (
	"context"

	"github.com/mp1947/ya-url-shortener/internal/model"
	"go.uber.org/zap"
)

// DeleteURLsBatch enqueues a batch of short URLs for deletion by sending them
// to the service's communication channel. It logs the operation and does not
// perform the deletion synchronously. The actual deletion is handled
// asynchronously by another component listening on the channel.
//
// Parameters:
//   - ctx: context for cancellation and deadlines.
//   - shortURLs: a batch of short URLs to be deleted.
func (s *ShortenService) DeleteURLsBatch(
	ctx context.Context,
	shortURLs model.BatchDeleteShortURLs,
) {
	s.Logger.Info(
		"putting short urls to delete into channel",
		zap.Any("data", shortURLs),
	)
	s.CommCh <- shortURLs

}

// ProcessDeletions starts a goroutine that listens for deletion requests on the CommCh channel.
// For each batch of data received, it attempts to delete the corresponding short URLs from the storage.
// The method logs the start of processing, each received deletion request, any errors encountered during deletion,
// and the result of each deletion operation, including the number of rows deleted and the user ID associated with the request.
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
