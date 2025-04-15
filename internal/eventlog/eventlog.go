package eventlog

import (
	"encoding/json"
	"os"

	"github.com/mp1947/ya-url-shortener/config"
)

type Event struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	UserID      string `json:"user_uuid"`
	IsDeleted   bool   `json:"is_deleted"`
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
