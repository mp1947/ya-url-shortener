package inmemory

import (
	"github.com/mp1947/ya-url-shortener/config"

	"github.com/mp1947/ya-url-shortener/internal/eventlog"
)

type Memory struct {
	data            map[string]string
	shortURLToEvent map[string]eventlog.Event
	cfg             config.Config
	EP              *eventlog.EventProcessor
	isInRestoreMode bool
	StorageType     string
}
