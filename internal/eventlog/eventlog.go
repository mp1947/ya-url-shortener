// Package eventlog provides functionality for logging URL shortening events to a file.
package eventlog

import (
	"encoding/json"
	"os"

	"github.com/mp1947/ya-url-shortener/config"
)

// Event holds URL shortening event data.
type Event struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	UserID      string `json:"user_uuid"`
	IsDeleted   bool   `json:"is_deleted"`
}

// EventProcessor logs events to a file.
type EventProcessor struct {
	File        *os.File
	Encoder     *json.Encoder
	CurrentUUID int
}

// NewEventProcessor creates an EventProcessor with the given config.
func NewEventProcessor(cfg config.Config) (*EventProcessor, error) {
	file, err := os.OpenFile(*cfg.FileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return &EventProcessor{
		File:        file,
		Encoder:     json.NewEncoder(file),
		CurrentUUID: 0,
	}, nil
}

// WriteEvent encodes the provided Event and writes it using the underlying Encoder.
// It returns an error if the encoding or writing process fails.
func (ep *EventProcessor) WriteEvent(e *Event) error {
	return ep.Encoder.Encode(&e)
}

// IncrementUUID increments the CurrentUUID field of the EventProcessor by one.
// This method is typically used to generate a new unique identifier for events
// processed by the EventProcessor.
func (ep *EventProcessor) IncrementUUID() {
	ep.CurrentUUID++
}
