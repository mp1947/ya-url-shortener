package inmemory

import (
	"bufio"
	"context"
	"encoding/json"
	"os"

	"github.com/mp1947/ya-url-shortener/internal/eventlog"
	"github.com/mp1947/ya-url-shortener/internal/model"
)

// Get retrieves the original URL and related information associated with the given shortURL from the in-memory storage.
// It returns a model.URL containing the original URL, the short URL ID, and a deletion status.
// If the shortURL does not exist, the OriginalURL field will be empty.
func (s *Memory) Get(ctx context.Context, shortURL string) (model.URL, error) {
	return model.URL{
		OriginalURL: s.data[shortURL],
		ShortURLID:  shortURL,
		IsDeleted:   false,
	}, nil
}

// GetType returns the type of storage used by the Memory repository as a string.
func (s *Memory) GetType() string {
	return s.StorageType
}

// GetURLsByUserID retrieves all URLs associated with the specified user ID from the file storage.
// It reads each line from the storage file, unmarshals it into an event structure, and collects
// the URLs that belong to the given user. Returns a slice of UserURL and an error if any occurs
// during file operations or JSON unmarshalling.
func (s *Memory) GetURLsByUserID(ctx context.Context, userID string) ([]model.UserURL, error) {
	file, err := os.OpenFile(*s.cfg.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var result []model.UserURL

	for scanner.Scan() {
		var event eventlog.Event
		line := scanner.Text()

		if err := json.Unmarshal([]byte(line), &event); err != nil {
			return nil, err
		}
		if event.UserID == userID {
			result = append(result, model.UserURL{
				ShortURLID:  event.ShortURL,
				OriginalURL: event.OriginalURL,
			})
		}
	}

	return result, nil
}
