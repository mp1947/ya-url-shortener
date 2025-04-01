package inmemory

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"strconv"

	"github.com/mp1947/ya-url-shortener/config"
	"github.com/mp1947/ya-url-shortener/internal/entity"
	shrterr "github.com/mp1947/ya-url-shortener/internal/errors"
	"github.com/mp1947/ya-url-shortener/internal/eventlog"
)

type Memory struct {
	data        map[string]string
	cfg         config.Config
	EP          *eventlog.EventProcessor
	StorageType string
}

func (s *Memory) Init(cfg config.Config, ctx context.Context) error {
	var err error

	s.cfg = cfg

	s.data = make(map[string]string)
	s.StorageType = "inmemory"
	s.EP, err = eventlog.NewEventProcessor(s.cfg)

	if err != nil {
		return err
	}
	return nil
}

func (s *Memory) Save(ctx context.Context, shortURLID, originalURL string) error {
	if s.data[shortURLID] == "" {
		s.data[shortURLID] = originalURL

		s.EP.IncrementUUID()
		event := eventlog.Event{
			UUID:        strconv.Itoa(s.EP.CurrentUUID),
			ShortURL:    shortURLID,
			OriginalURL: originalURL,
		}
		s.EP.WriteEvent(&event)
		return nil
	}
	return shrterr.ErrOriginalURLAlreadyExists
}

func (s *Memory) SaveBatch(ctx context.Context, urls []entity.URL) (bool, error) {
	for _, v := range urls {
		s.data[v.ShortURLID] = v.OriginalURL
	}
	return true, nil
}

func (s *Memory) Get(ctx context.Context, shortURL string) (string, error) {
	return s.data[shortURL], nil
}

func (s *Memory) GetType() string {
	return s.StorageType
}

func (s *Memory) Ping(ctx context.Context) error {
	return nil
}

func (s *Memory) RestoreFromFile() (int, error) {
	file, err := os.OpenFile(*s.cfg.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	currentUUID := 0
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for scanner.Scan() {
		var event eventlog.Event
		line := scanner.Text()

		if err := json.Unmarshal([]byte(line), &event); err != nil {
			return 0, err
		}
		s.Save(ctx, event.ShortURL, event.OriginalURL)

		currentUUID += 1
	}
	s.EP.CurrentUUID = currentUUID
	return currentUUID, nil
}
