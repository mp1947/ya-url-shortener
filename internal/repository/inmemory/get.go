package inmemory

import (
	"bufio"
	"context"
	"encoding/json"
	"os"

	"github.com/mp1947/ya-url-shortener/internal/entity"
	"github.com/mp1947/ya-url-shortener/internal/eventlog"
)

func (s *Memory) Get(ctx context.Context, shortURL string) (entity.URL, error) {
	return entity.URL{
		OriginalURL: s.data[shortURL],
		ShortURLID:  shortURL,
		IsDeleted:   false,
	}, nil
}

func (s *Memory) GetType() string {
	return s.StorageType
}

func (s *Memory) GetURLsByUserID(ctx context.Context, userID string) ([]entity.UserURL, error) {
	file, err := os.OpenFile(*s.cfg.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var result []entity.UserURL

	for scanner.Scan() {
		var event eventlog.Event
		line := scanner.Text()

		if err := json.Unmarshal([]byte(line), &event); err != nil {
			return nil, err
		}
		if event.UserID == userID {
			result = append(result, entity.UserURL{
				ShortURLID:  event.ShortURL,
				OriginalURL: event.OriginalURL,
			})
		}
	}

	return result, nil
}
