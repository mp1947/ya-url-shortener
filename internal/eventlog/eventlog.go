package eventlog

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/mp1947/ya-url-shortener/config"
	"github.com/mp1947/ya-url-shortener/internal/repository"
	"go.uber.org/zap"
)

type Event struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type EventProcessor struct {
	File        *os.File
	Encoder     *json.Encoder
	CurrentUUID int
}

func NewEventProcessor(cfg config.Config) (*EventProcessor, error) {
	file, err := os.OpenFile(
		*cfg.FileStoragePath,
		os.O_WRONLY|os.O_CREATE|os.O_APPEND,
		0666,
	)

	if err != nil {
		return nil, err
	}
	return &EventProcessor{
		File:        file,
		Encoder:     json.NewEncoder(file),
		CurrentUUID: 0,
	}, nil
}

func (ep *EventProcessor) WriteEvent(e *Event) error {
	return ep.Encoder.Encode(&e)
}

func (ep *EventProcessor) IncrementUUID() {
	ep.CurrentUUID++
}

func (ep *EventProcessor) setUUID(uuid int) {
	ep.CurrentUUID = uuid
}

func (ep *EventProcessor) RestoreFromFile(
	cfg config.Config,
	r repository.Repository,
	logger *zap.Logger,
) (int, error) {
	file, err := os.OpenFile(*cfg.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	currentUUID := 0

	for scanner.Scan() {
		var event Event
		line := scanner.Text()

		if err := json.Unmarshal([]byte(line), &event); err != nil {
			logger.Warn("unmarshal error", zap.Int("uuid", currentUUID), zap.Error(err))
		}
		if r.GetType() == "inmemory" {
			r.Save(event.ShortURL, event.OriginalURL)
		}

		currentUUID += 1
	}
	ep.setUUID(currentUUID)
	return currentUUID, nil
}
