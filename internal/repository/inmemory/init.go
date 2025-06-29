package inmemory

import (
	"context"

	"github.com/mp1947/ya-url-shortener/config"
	"github.com/mp1947/ya-url-shortener/internal/eventlog"
	"go.uber.org/zap"
)

// Init initializes the in-memory storage for the Memory repository.
// It sets up the configuration, logger, and internal data structures,
// and creates a new event processor. Returns an error if the event
// processor cannot be created.
func (s *Memory) Init(
	ctx context.Context,
	cfg config.Config,
	l *zap.Logger,
) error {
	var err error

	s.cfg = cfg
	s.isInRestoreMode = false
	s.data = make(map[string]string)
	s.shortURLToEvent = make(map[string]eventlog.Event)
	s.StorageType = "inmemory"
	s.EP, err = eventlog.NewEventProcessor(s.cfg)

	if err != nil {
		return err
	}
	return nil
}
