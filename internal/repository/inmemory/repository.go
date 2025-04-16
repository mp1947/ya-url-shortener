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
	"go.uber.org/zap"
)

type Memory struct {
	data            map[string]string
	cfg             config.Config
	EP              *eventlog.EventProcessor
	isInRestoreMode bool
	StorageType     string
}

func (s *Memory) Init(
	ctx context.Context,
	cfg config.Config,
	l *zap.Logger,
) error {
	var err error

	s.cfg = cfg
	s.isInRestoreMode = false
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
		if !s.isInRestoreMode {
			s.EP.WriteEvent(&event)
		}
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

func (s *Memory) RestoreFromFile(l *zap.Logger) (int, error) {
	file, err := os.OpenFile(*s.cfg.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	currentUUID := 0
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.isInRestoreMode = true

	for scanner.Scan() {
		var event eventlog.Event
		line := scanner.Text()

		if err := json.Unmarshal([]byte(line), &event); err != nil {
			return 0, err
		}
		if err := s.Save(ctx, event.ShortURL, event.OriginalURL); err != nil {
			l.Warn("error saving record to file during restore phase", zap.Error(err))
		}

		currentUUID += 1
	}
	s.EP.CurrentUUID = currentUUID
	s.isInRestoreMode = false

	return currentUUID, nil
}
