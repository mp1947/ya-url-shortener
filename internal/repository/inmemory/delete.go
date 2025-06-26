package inmemory

import (
	"context"

	"github.com/mp1947/ya-url-shortener/internal/eventlog"
	"github.com/mp1947/ya-url-shortener/internal/model"
)

// DeleteBatch removes a batch of short URLs for a user from memory.
// Returns the number of processed URLs and an error if any.
func (s *Memory) DeleteBatch(
	ctx context.Context,
	shortURLs model.BatchDeleteShortURLs,
) (int64, error) {
	var counter int64
	for _, v := range shortURLs.ShortURLs {
		event := s.shortURLToEvent[v]
		if event.UserID == shortURLs.UserID {
			s.data[v] = ""
			s.shortURLToEvent[v] = eventlog.Event{}
		}
		counter++
	}
	return counter, nil
}
