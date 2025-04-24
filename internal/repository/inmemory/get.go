package inmemory

import (
	"bufio"
	"context"
	"encoding/json"
	"os"

	"github.com/mp1947/ya-url-shortener/internal/eventlog"
	"github.com/mp1947/ya-url-shortener/internal/model"
)

func (s *Memory) Get(ctx context.Context, shortURL string) (model.URL, error) {
	return model.URL{
		OriginalURL: s.data[shortURL],
		ShortURLID:  shortURL,
		IsDeleted:   false,
	}, nil
}

func (s *Memory) GetType() string {
	return s.StorageType
}

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
