package inmemory

import (
	"context"

	"github.com/mp1947/ya-url-shortener/internal/eventlog"
	"github.com/mp1947/ya-url-shortener/internal/model"
)

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
