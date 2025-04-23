package inmemory

import (
	"context"
	"strconv"

	"github.com/mp1947/ya-url-shortener/internal/eventlog"
	"github.com/mp1947/ya-url-shortener/internal/model"

	shrterr "github.com/mp1947/ya-url-shortener/internal/errors"
)

func (s *Memory) Save(
	ctx context.Context,
	shortURLID,
	originalURL string,
	userID string,
) error {
	if s.data[shortURLID] == "" {
		s.data[shortURLID] = originalURL

		s.EP.IncrementUUID()
		event := eventlog.Event{
			UUID:        strconv.Itoa(s.EP.CurrentUUID),
			ShortURL:    shortURLID,
			OriginalURL: originalURL,
			UserID:      userID,
			IsDeleted:   false,
		}
		s.shortURLToEvent[shortURLID] = event
		if !s.isInRestoreMode {
			s.EP.WriteEvent(&event)
		}
		return nil
	}
	return shrterr.ErrOriginalURLAlreadyExists
}

func (s *Memory) SaveBatch(
	ctx context.Context,
	urls []model.URLWithCorrelation,
	userID string,
) (bool, error) {
	for _, v := range urls {
		// s.data[v.ShortURLID] = v.OriginalURL
		err := s.Save(ctx, v.ShortURLID, v.OriginalURL, userID)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}
