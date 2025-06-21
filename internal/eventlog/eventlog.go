package eventlog

import (
	"encoding/json"
	"os"

	"github.com/mp1947/ya-url-shortener/config"
)

// Event represents a URL shortening event with associated metadata.
type Event struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	UserID      string `json:"user_uuid"`
	IsDeleted   bool   `json:"is_deleted"`
}

// EventProcessor handles event logging by writing encoded events to a file and tracking the current UUID.
type EventProcessor struct {
	File        *os.File
	Encoder     *json.Encoder
	CurrentUUID int
}

// NewEventProcessor creates a new instance of EventProcessor using the provided configuration.
// It opens the file specified by cfg.FileStoragePath for writing, creating it if it does not exist,
// and appending to it if it does. The function returns a pointer to the initialized EventProcessor
// or an error if the file cannot be opened.
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
