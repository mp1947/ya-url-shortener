package service

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/mp1947/ya-url-shortener/config"
)

type Event struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type EventProcessor struct {
	File        *os.File
	Encoder     *json.Encoder
	currentUUID int
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
		currentUUID: 0,
	}, nil
}

func (ep *EventProcessor) writeEvent(e *Event) error {
	return ep.Encoder.Encode(&e)
}

func (ep *EventProcessor) incrementUUID() {
	ep.currentUUID++
}

func (ep *EventProcessor) setUUID(uuid int) {
	ep.currentUUID = uuid
}

func (ep *EventProcessor) RestoreFromFile(cfg config.Config) (int, error) {
	file, err := os.OpenFile(*cfg.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return 0, err
	}
	scanner := bufio.NewScanner(file)
	currentUUID := 0

	for scanner.Scan() {
		var event Event
		line := scanner.Text()

		if err := json.Unmarshal([]byte(line), &event); err != nil {
			return 0, err
		}

		currentUUID += 1
	}
	ep.setUUID(currentUUID)
	return currentUUID, nil
}
