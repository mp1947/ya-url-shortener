package inmemory

import (
	"context"
	"strconv"

	"github.com/mp1947/ya-url-shortener/internal/eventlog"
	"github.com/mp1947/ya-url-shortener/internal/model"

	shrterr "github.com/mp1947/ya-url-shortener/internal/errors"
)

// Save stores the mapping between a short URL ID and its original URL for a given user.
// If the short URL ID does not already exist in memory, it saves the mapping, logs the event,
// and optionally writes the event to persistent storage unless in restore mode.
// Returns an error if the short URL ID already exists.
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
			err := s.EP.WriteEvent(&event)
			if err != nil {
				return err
			}
		}
		return nil
	}
	return shrterr.ErrOriginalURLAlreadyExists
}

// SaveBatch saves a batch of URL mappings for a specific user into memory.
// It iterates over the provided slice of URLWithCorrelation, saving each URL using the Save method.
// If any save operation fails, it returns false and the encountered error immediately.
// On success, it returns true and a nil error.
//
// Parameters:
//   - ctx: context for cancellation and deadlines.
//   - urls: slice of URLWithCorrelation containing short and original URLs to save.
//   - userID: identifier of the user to associate with the URLs.
//
// Returns:
//   - bool: true if all URLs were saved successfully, false otherwise.
//   - error: error encountered during saving, or nil if successful.
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
