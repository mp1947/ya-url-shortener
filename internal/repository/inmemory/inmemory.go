// Package inmemory provides in-memory storage for short URLs, mainly for testing and development.
package inmemory

import (
	"github.com/mp1947/ya-url-shortener/config"

	"github.com/mp1947/ya-url-shortener/internal/eventlog"
)

// Memory represents an in-memory storage for URL shortening service data.
// It maintains mappings between original URLs and their shortened versions,
// as well as event logs associated with shortened URLs. The struct also
// holds configuration settings, an event processor for handling events,
// a flag indicating if the storage is in restore mode, and the type of storage used.

type Memory struct {
	EP              *eventlog.EventProcessor
	data            map[string]string
	shortURLToEvent map[string]eventlog.Event
	cfg             config.Config
	StorageType     string
	isInRestoreMode bool
}
