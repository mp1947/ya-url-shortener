package inmemory

import (
	"context"

	"github.com/mp1947/ya-url-shortener/internal/eventlog"
	"github.com/mp1947/ya-url-shortener/internal/model"
)

// DeleteBatch deletes a batch of short URLs associated with a specific user from the in-memory storage.
// It iterates over the provided list of short URLs, and for each URL, if it belongs to the given user,
// it removes the URL from the storage and resets its associated event. The function returns the total
// number of processed URLs and an error if any occurred.
//
// Parameters:
//
//	ctx        - The context for cancellation and deadlines.
//	shortURLs  - A BatchDeleteShortURLs struct containing the user ID and the list of short URLs to delete.
//
// Returns:
//
//	int64 - The number of processed URLs.
//	error - An error if the operation fails, otherwise nil.
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
